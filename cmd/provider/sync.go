package provider

import (
	"context"
	"github.com/selefra/selefra/cmd/fetch"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"path/filepath"
)

func Sync() error {
	ui.PrintSuccessLn("Initializing provider plugins...")
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
	ui.PrintSuccessF("Selefra has been successfully installed providers!")
	ui.PrintSuccessF("Checking Selefra provider updates......")
	for _, p := range cof.Selefra.Providers {
		prov := registry.Provider{
			Name:    p.Name,
			Version: p.Version,
			Source:  "",
		}
		pp, err := provider.Download(ctx, prov, true)
		if err != nil {
			ui.PrintErrorF("%s@%s failed updated：%s", p.Name, p.Version, err.Error())
			continue
		} else {
			ui.PrintSuccessF("%s@%s all ready updated!", p.Name, p.Version)
		}

		p.Path = pp.Filepath
		p.Version = pp.Version
		err = fetch.Fetch(ctx, cof, p)
		if err != nil {
			ui.PrintErrorF("%s %s Synchronization failed：%s", p.Name, p.Version, err.Error())
			continue
		}
	}

	ui.PrintSuccessF(`
This may be exception, view detailed exception in %s .

Need help? Know on Slack or open a Github Issue: https://github.com/selefra/selefra#community
`, filepath.Join(*global.WORKSPACE, "logs"))
	return nil
}
