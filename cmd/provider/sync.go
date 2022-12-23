package provider

import (
	"context"
	"fmt"
	"github.com/selefra/selefra/cmd/fetch"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"path/filepath"
	"time"
)

func Sync(lockCtx context.Context) (errLogs []string, err error) {
	ui.PrintSuccessLn("Initializing provider plugins...\n")
	ctx := context.Background()
	var cof = &config.SelefraConfig{}
	err = cof.GetConfig()
	if err != nil {
		return nil, err
	}
	namespace, _, err := utils.Home()
	if err != nil {
		return nil, err
	}
	provider := registry.NewProviderRegistry(namespace)
	ui.PrintSuccessF("Selefra has been successfully installed providers!\n")
	ui.PrintSuccessF("Checking Selefra provider updates......\n")

	var hasError bool
	var ProviderRequires []*config.ProviderRequired
	for _, p := range cof.Selefra.Providers {
		configVersion := p.Version
		prov := registry.Provider{
			Name:    p.Name,
			Version: p.Version,
			Source:  "",
			Path:    p.Path,
		}
		pp, err := provider.Download(ctx, prov, true)
		if err != nil {
			hasError = true
			ui.PrintErrorF("%s@%s failed updated：%s", p.Name, p.Version, err.Error())
			errLogs = append(errLogs, fmt.Sprintf("%s@%s failed updated：%s", p.Name, p.Version, err.Error()))
			continue
		} else {
			p.Path = pp.Filepath
			p.Version = pp.Version
			err = tools.SetSelefraProvider(pp, nil, configVersion)
			if err != nil {
				hasError = true
				ui.PrintErrorF("%s@%s failed updated：%s", p.Name, p.Version, err.Error())
				errLogs = append(errLogs, fmt.Sprintf("%s@%s failed updated：%s", p.Name, p.Version, err.Error()))
				continue
			}
			ProviderRequires = append(ProviderRequires, p)
			ui.PrintSuccessF("	%s@%s all ready updated!\n", p.Name, p.Version)
		}
	}

	ui.PrintSuccessF("Selefra has been finished update providers!\n")
	global.STAG = "pull"
	for _, p := range ProviderRequires {
		for {
			locked, _ := tools.Locked(*p, *cof)
			if !locked {
				break
			}
			time.Sleep(5 * time.Second)
			//ui.PrintWarningF(fmt.Sprintf("	%s@%s is locked, waiting for unlock...\n", p.Name, p.Version))
		}
		go tools.Lock(lockCtx, *p, *cof)
		need, _ := tools.NeedFetch(*p, *cof)
		if !need {
			continue
		}
		err = fetch.Fetch(ctx, cof, p)
		if err != nil {
			ui.PrintErrorF("%s %s Synchronization failed：%s", p.Name, p.Version, err.Error())
			hasError = true
			continue
		}
		requireKey := config.GetCacheKey()
		err = tools.SetStoreValue(*cof, p, requireKey, time.Now().Format(time.RFC3339))
		if err != nil {
			ui.PrintWarningF("%s %s set cache time failed：%s", p.Name, p.Version, err.Error())
			hasError = true
			continue
		}
	}
	if hasError {
		ui.PrintErrorF(`
This may be exception, view detailed exception in %s .
`, filepath.Join(*global.WORKSPACE, "logs"))
	}

	return errLogs, nil
}
