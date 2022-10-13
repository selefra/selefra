package fetch

import (
	"context"
	"encoding/json"
	"errors"
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
				err = Fetch(ctx, cof, p)
				if err != nil {
					return err
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

func Fetch(ctx context.Context, cof *config.SelefraConfig, p *config.ProviderRequired) error {
	if p.Path == "" {
		p.Path = utils.GetPathBySource(*p.Source)
	}
	var providersName = utils.GetNameBySource(*p.Source)
	ui.PrintSuccessF("%s@%s pull infrastructure data:\n", providersName, p.Version)
	plug, err := plugin.NewManagedPlugin(p.Path, providersName, p.Version, "", nil)
	if err != nil {
		return err
	}

	storage := postgresql_storage.NewPostgresqlStorageOptions(cof.Selefra.GetDSN())
	opt, err := json.Marshal(storage)
	if err != nil {
		return err
	}
	err = cof.GetConfig()

	conf, err := tools.GetProviders(cof, providersName)

	if err != nil {
		return err
	}
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
		if initRes.Diagnostics != nil && initRes.Diagnostics.HasError() {
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
	if dropRes.Diagnostics != nil && dropRes.Diagnostics.HasError() {
		ui.PrintDiagnostic(dropRes.Diagnostics.GetDiagnosticSlice())
		return errors.New("fetch provider drop table error")
	}

	createRes, err := provider.CreateAllTables(ctx, &shard.ProviderCreateAllTablesRequest{})
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}
	if createRes.Diagnostics != nil && createRes.Diagnostics.HasError() {
		ui.PrintDiagnostic(createRes.Diagnostics.GetDiagnosticSlice())
		return errors.New("fetch provider create table error")
	}

	recv, err := provider.PullTables(ctx, &shard.PullTablesRequest{
		Tables:        []string{"*"},
		MaxGoroutines: 100,
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
	for {
		current := 0
		res, err := recv.Recv()

		if err != nil {
			if errors.Is(err, io.EOF) {
				progbar.Done(p.Name + "@" + p.Version)
				break
			}
			return err
		}
		successNum := 0
		errorsNum := 0
		for _, value := range res.FinishedTables {
			if value {
				successNum++
			} else {
				errorsNum++
			}
		}
		success = successNum
		errorsN = errorsNum
		progbar.SetTotal(p.Name+"@"+p.Version, int64(res.TableCount))
		progbar.Current(p.Name+"@"+p.Version, int64(len(res.FinishedTables)), res.Table)
		if res.Diagnostics != nil && res.Diagnostics.HasError() {
			_ = ui.SaveLogToDiagnostic(res.Diagnostics.GetDiagnosticSlice())
		}

		if res != nil {
			for s := range res.FinishedTables {
				if res.FinishedTables[s] {
					current++
				}
			}
		}
	}
	progbar.Wait(p.Name + "@" + p.Version)

	ui.PrintSuccessF("\nPull complete! Total Resources pulled:%d        Errors: %d\n", success, errorsN)
	return nil
}
