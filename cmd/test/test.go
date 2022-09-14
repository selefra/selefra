package test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func NewTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Check whether the configuration is valid",
		Long:  "Check whether the configuration is valid",
		RunE:  testFunc,
	}

	cmd.PersistentFlags().StringP("dir", "d", ".", "the directory to initialize in")

	cmd.SetHelpFunc(cmd.HelpFunc())

	return cmd
}

func testFunc(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	wd, err := os.Getwd()
	dirname, _ := cmd.PersistentFlags().GetString("dir")
	*global.WORKSPACE = filepath.Join(wd, dirname)
	s := config.SelefraConfig{}
	err = s.GetConfigByNode()
	if err != nil {
		ui.PrintErrorF(err.Error())
		return nil
	}
	cof, err := s.GetConfigWithViper()
	if err != nil {
		ui.PrintErrorF("Profile deserialization exception:%s", err.Error())
		return nil
	}
	err = checkConfig(ctx, s)
	if err != nil {
		ui.PrintErrorF("selefra configuration exception:%s", err.Error())
		return nil
	}
	ui.PrintSuccessF("Client Verification Success")
	for _, p := range s.Selefra.Providers {
		if p.Path == "" {
			p.Path = utils.GetPathBySource(*p.Source)
		}
		plug, err := plugin.NewManagedPlugin(p.Path, p.Name, p.Version, "", nil)
		if err != nil {

			ui.PrintErrorF("%s %s verification failed ：%s", p.Name, p.Version, err.Error())
			continue
		}
		conf, err := yaml.Marshal(cof.Get("providers." + p.Name))
		if err != nil {
			ui.PrintErrorLn(err.Error())
			continue
		}

		storage := postgresql_storage.NewPostgresqlStorageOptions(s.Selefra.GetDSN())
		opt, err := json.Marshal(storage)

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
			ui.PrintErrorF("%s %s verification failed ：%s", p.Name, p.Version, err.Error())
			continue
		} else {
			if initRes.Diagnostics != nil && initRes.Diagnostics.HasError() {
				ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
				continue
			}
		}

		res, err := provider.SetProviderConfig(ctx, &shard.SetProviderConfigRequest{
			Storage: &shard.Storage{
				Type:           0,
				StorageOptions: opt,
			},
			ProviderConfig: pointer.ToStringPointer(string(conf)),
		})
		if err != nil {
			ui.PrintErrorLn(err.Error())
			continue
		} else {
			if res.Diagnostics != nil && res.Diagnostics.HasError() {
				ui.PrintDiagnostic(res.Diagnostics.GetDiagnosticSlice())
				continue
			}
		}
		ui.PrintSuccessF("%s %s check successfully", p.Name, p.Version)
	}

	ui.PrintSuccessF("Providers verification completed")
	ui.PrintSuccessF("Profile verification complete")
	return nil
}

func checkConfig(ctx context.Context, c config.SelefraConfig) error {
	var err error
	if c.Selefra.CliVersion == "" {
		err = errors.New("cliVersion is empty")
		return err
	}
	uid, _ := uuid.NewUUID()
	_, e := client.CreateClientFromConfig(ctx, &c.Selefra, uid)
	if e != nil {
		return e
	}
	return nil
}
