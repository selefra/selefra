package provider

import (
	"context"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/id_util"
	"github.com/selefra/selefra/cmd/fetch"
	"github.com/selefra/selefra/cmd/test"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/grpcClient"
	"github.com/selefra/selefra/pkg/httpClient"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"time"
)

type lockStruct struct {
	SchemaKey string
	Uuid      string
	Storage   *postgresql_storage.PostgresqlStorage
}

func Sync() (errLogs []string, lockSlice []lockStruct, err error) {
	ui.PrintSuccessLn("Initializing provider plugins...\n")
	ctx := context.Background()
	var cof = &config.SelefraConfig{}
	err = cof.GetConfig()
	if err != nil {
		return nil, nil, err
	}
	namespace, _, err := utils.Home()
	if err != nil {
		return nil, nil, err
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

	err = test.CheckSelefraConfig(ctx, *cof)
	if err != nil {
		ui.PrintErrorLn(err.Error())
		if global.LOGINTOKEN != "" && cof.Selefra.Cloud != nil && err == nil {
			_ = httpClient.SetupStag(global.LOGINTOKEN, cof.Selefra.Cloud.Project, httpClient.Failed)
		}
		return nil, nil, err
	}

	_, err = grpcClient.Cli.UploadLogStatus()
	if err != nil {
		ui.PrintErrorLn(err.Error())
	}
	global.STAG = "pull"
	for _, p := range ProviderRequires {
		confs, err := tools.GetProviders(cof, p.Name)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			continue
		}
		for _, conf := range confs {
			store, err := tools.GetStore(*cof, p, conf)
			if err != nil {
				hasError = true
				ui.PrintErrorF("%s@%s failed updated：%s", p.Name, p.Version, err.Error())
				errLogs = append(errLogs, fmt.Sprintf("%s@%s failed updated：%s", p.Name, p.Version, err.Error()))
				continue
			}
			ctx := context.Background()
			uuid := id_util.RandomId()
			var cp config.CliProviders
			err = yaml.Unmarshal([]byte(conf), &cp)
			if err != nil {
				ui.PrintErrorLn(err.Error())
				continue
			}
			schemaKey := config.GetSchemaKey(p, cp)
			for {
				err = store.Lock(ctx, schemaKey, uuid)
				if err == nil {
					break
				}
				time.Sleep(5 * time.Second)
			}
			lockSlice = append(lockSlice, lockStruct{
				SchemaKey: schemaKey,
				Uuid:      uuid,
				Storage:   store,
			})
			need, _ := tools.NeedFetch(*p, *cof, conf)
			if !need {
				continue
			}
			err = fetch.Fetch(ctx, cof, p, conf)
			if err != nil {
				ui.PrintErrorF("%s %s Synchronization failed：%s", p.Name, p.Version, err.Error())
				hasError = true
				continue
			}
			requireKey := config.GetCacheKey()
			err = tools.SetStoreValue(*cof, p, conf, requireKey, time.Now().Format(time.RFC3339))
			if err != nil {
				ui.PrintWarningF("%s %s set cache time failed：%s", p.Name, p.Version, err.Error())
				hasError = true
				continue
			}
		}
	}
	if hasError {
		ui.PrintErrorF(`
This may be exception, view detailed exception in %s .
`, filepath.Join(*global.WORKSPACE, "logs"))
	}

	return errLogs, lockSlice, nil
}
