package test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/pointer"
	utils2 "github.com/selefra/selefra/cmd/utils"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/client"
	"github.com/spf13/cobra"
	"os"
)

func NewTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Check whether the configuration is valid",
		Long:  "Check whether the configuration is valid",
		RunE:  testFunc,
	}

	cmd.SetHelpFunc(cmd.HelpFunc())

	return cmd
}

func testFunc(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	wd, err := os.Getwd()
	*global.WORKSPACE = wd
	s := config.SelefraConfig{}
	err = s.GetConfigByNode()
	if err != nil {
		ui.PrintErrorF(err.Error())
		return nil
	}
	err = s.GetConfig()
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
		var providersName = utils.GetNameBySource(*p.Source)
		plug, err := plugin.NewManagedPlugin(p.Path, providersName, p.Version, "", nil)
		if err != nil {

			ui.PrintErrorF("%s %s verification failed ：%s", p.Name, p.Version, err.Error())
			continue
		}
		conf, err := utils2.GetProviders(&s, p.Name)
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
	if c.Selefra.Name == "" {
		err = errors.New("name is empty")
		return err
	}
	uid, _ := uuid.NewUUID()
	_, e := client.CreateClientFromConfig(ctx, &c.Selefra, uid)
	if e != nil {
		return e
	}
	return nil
}
