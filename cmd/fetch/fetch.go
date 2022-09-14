package fetch

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

func NewFetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "fetch",
		Long:  "fetch",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var cof = &config.SelefraConfig{}

			wd, err := os.Getwd()
			dirname, _ := cmd.PersistentFlags().GetString("dir")
			*global.WORKSPACE = filepath.Join(wd, dirname)
			err = cof.GetConfig()
			if err != nil {
				return err
			}
			for _, p := range cof.Selefra.Providers {
				err = Fetch(ctx, cof, p)
				if err != nil {
					return err
				}
			}
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

	plug, err := plugin.NewManagedPlugin(p.Path, p.Name, p.Version, "", nil)
	if err != nil {
		return err
	}

	storage := postgresql_storage.NewPostgresqlStorageOptions(cof.Selefra.GetDSN())
	opt, err := json.Marshal(storage)
	if err != nil {
		return err
	}
	v, err := cof.GetConfigWithViper()

	conf, err := yaml.Marshal(v.Get("providers." + p.Name))

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
	for {
		current := 0
		res, err := recv.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if res.Diagnostics != nil && res.Diagnostics.HasError() {
			_ = ui.PrintDiagnostic(res.Diagnostics.GetDiagnosticSlice())
		}

		if res != nil {
			for s := range res.FinishedTables {
				if res.FinishedTables[s] {
					current++
				}
			}
		}
	}
	return nil
}
