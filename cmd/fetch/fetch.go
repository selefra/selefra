package fetch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/progress"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

func NewFetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch resources from configured providers",
		Long:  "Fetch resources from configured providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			global.CMD = "fetch"
			ctx := cmd.Context()
			var cof = &config.SelefraConfig{}

			wd, err := os.Getwd()
			*global.WORKSPACE = wd
			err = cof.GetConfig()
			if err != nil {
				return err
			}
			ui.PrintSuccessF("Selefra start fetch")
			for _, p := range cof.Selefra.Providers {
				confs, err := tools.GetProviders(cof, p.Name)
				if err != nil {
					return err
				}
				for i := range confs {
					err = Fetch(ctx, cof, p, confs[i])
					if err != nil {
						return err
					}
				}
			}

			ui.PrintErrorF(`
This may be exception, view detailed exception in %s.`,
				filepath.Join(*global.WORKSPACE, "logs"))

			return nil
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func Fetch(ctx context.Context, cof *config.SelefraConfig, p *config.ProviderRequired, conf string) error {
	var cp config.CliProviders
	err := yaml.Unmarshal([]byte(conf), &cp)
	if err != nil {
		return err
	}

	if p.Path == "" {
		p.Path = utils.GetPathBySource(*p.Source)
	}
	var providersName = utils.GetNameBySource(*p.Source)
	ui.PrintSuccessF("%s %s@%s pull infrastructure data:\n", cp.Name, providersName, p.Version)
	ui.PrintCustomizeLnNotShow(fmt.Sprintf("Pulling %s@%s Please wait for resource information ...", providersName, p.Version))
	plug, err := plugin.NewManagedPlugin(p.Path, providersName, p.Version, "", nil)
	if err != nil {
		return err
	}

	storageOpt := postgresql_storage.NewPostgresqlStorageOptions(cof.Selefra.GetDSN())

	schema := config.GetSchemaKey(p, cp)
	storageOpt.SearchPath = schema
	opt, err := json.Marshal(storageOpt)
	if err != nil {
		return err
	}
	err = cof.GetConfig()

	provider := plug.Provider()
	initRes, err := provider.Init(ctx, &shard.ProviderInitRequest{
		Workspace: global.WORKSPACE,
		Storage: &shard.Storage{
			Type:           0,
			StorageOptions: opt,
		},
		IsInstallInit:  pointer.FalsePointer(),
		ProviderConfig: pointer.ToStringPointer(string(conf)),
	})
	if err != nil {
		return err
	} else {
		if initRes.Diagnostics != nil {
			ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
			return errors.New("fetch provider init error")
		}
	}

	defer plug.Close()
	dropRes, err := provider.DropTableAll(ctx, &shard.ProviderDropTableAllRequest{})
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}
	if dropRes.Diagnostics != nil {
		err := ui.PrintDiagnostic(dropRes.Diagnostics.GetDiagnosticSlice())
		if err != nil {
			return errors.New("fetch provider drop table error")
		}
	}

	createRes, err := provider.CreateAllTables(ctx, &shard.ProviderCreateAllTablesRequest{})
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}
	if createRes.Diagnostics != nil {
		err := ui.PrintDiagnostic(createRes.Diagnostics.GetDiagnosticSlice())
		if err != nil {
			return errors.New("fetch provider create table error")
		}
	}
	var tables []string
	resources := cp.Resources
	if len(resources) == 0 {
		tables = append(tables, "*")
	} else {
		for i := range resources {
			tables = append(tables, resources[i])
		}
	}
	var maxGoroutines uint64 = 100
	if cp.MaxGoroutines > 0 {
		maxGoroutines = cp.MaxGoroutines
	}
	recv, err := provider.PullTables(ctx, &shard.PullTablesRequest{
		Tables:        tables,
		MaxGoroutines: maxGoroutines,
		Timeout:       0,
	})
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}
	progbar := progress.CreateProgress()
	progbar.Add(p.Name+"@"+p.Version, -1)
	success := 0
	errorsN := 0
	var total int64
	for {
		res, err := recv.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				progbar.Current(p.Name+"@"+p.Version, total, "done")
				progbar.Done(p.Name + "@" + p.Version)
				break
			}
			return err
		}
		progbar.SetTotal(p.Name+"@"+p.Version, int64(res.TableCount))
		progbar.Current(p.Name+"@"+p.Version, int64(len(res.FinishedTables)), res.Table)
		total = int64(res.TableCount)
		if res.Diagnostics != nil {
			if res.Diagnostics.HasError() {
				errorsN++
				ui.SaveLogToDiagnostic(res.Diagnostics.GetDiagnosticSlice())
			} else {
				success++
			}
		}
	}
	progbar.Wait(p.Name + "@" + p.Version)
	if errorsN > 0 {
		ui.PrintErrorF("\nPull complete! Total Resources pulled:%d        Errors: %d\n", success, errorsN)
		return nil
	}
	ui.PrintSuccessF("\nPull complete! Total Resources pulled:%d        Errors: %d\n", success, errorsN)
	return nil
}
