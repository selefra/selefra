package provider

import (
	"context"
	"github.com/selefra/selefra/cmd/fetch"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"path/filepath"
)

func Sync() error {
	ui.PrintSuccessLn("Initializing provider plugins...\n")
	ctx := context.Background()
	var cof = &config.SelefraConfig{}
	err := cof.GetConfig()
	if err != nil {
		return err
	}
	namespace, _, err := utils.Home()
	if err != nil {
		return err
	}
	provider := registry.NewProviderRegistry(namespace)
	ui.PrintSuccessF("Selefra has been successfully installed providers!\n")
	ui.PrintSuccessF("Checking Selefra provider updates......\n")

	var hasError bool
	var ProviderRequires []*config.ProviderRequired
	for _, p := range cof.Selefra.Providers {
		prov := registry.Provider{
			Name:    p.Name,
			Version: p.Version,
			Source:  "",
		}
		pp, err := provider.Download(ctx, prov, true)
		if err != nil {
			hasError = true
			ui.PrintErrorF("%s@%s failed updated：%s", p.Name, p.Version, err.Error())
			continue
		} else {
			p.Path = pp.Filepath
			p.Version = pp.Version
			tools.SetSelefraProvider(pp, nil)
			ProviderRequires = append(ProviderRequires, p)
			ui.PrintSuccessF("	%s@%s all ready updated!\n", p.Name, p.Version)
		}
	}

	ui.PrintSuccessF("Selefra has been finished update providers!\n")
	global.STAG = "pull"
	for _, p := range ProviderRequires {
		err = fetch.Fetch(ctx, cof, p)
		if err != nil {
			ui.PrintErrorF("%s %s Synchronization failed：%s", p.Name, p.Version, err.Error())
			hasError = true
			continue
		}
	}

	if hasError {
		ui.PrintErrorF(`
This may be exception, view detailed exception in %s .
`, filepath.Join(*global.WORKSPACE, "logs"))
	}

	return nil
}
