package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cmd/tools"
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

func TestFunc(ctx context.Context) error {
	err := config.IsSelefra()
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}

	if err != nil {
		ui.PrintErrorLn("GetWDError:" + err.Error())
	}
	s := config.SelefraConfig{}
	return CheckSelefraConfig(ctx, s)
}

func testFunc(cmd *cobra.Command, args []string) error {
	global.CMD = "test"
	ctx := cmd.Context()
	wd, err := os.Getwd()
	if err != nil {
		ui.PrintErrorLn("Error:" + err.Error())
		return nil
	}
	*global.WORKSPACE = wd
	return TestFunc(ctx)
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

func CheckSelefraConfig(ctx context.Context, s config.SelefraConfig) error {
	err := s.TestConfigByNode()
	if err != nil {
		return err
	}
	err = s.GetConfig()
	if err != nil {
		return errors.New(fmt.Sprintf("Profile deserialization exception:%s", err.Error()))
	}
	err = checkConfig(ctx, s)
	if err != nil {
		return errors.New(fmt.Sprintf("selefra configuration exception:%s", err.Error()))
	}
	ui.PrintSuccessF("Client Verification Success\n")
	hasError := false
	for _, p := range s.Selefra.Providers {
		if p.Path == "" {
			p.Path = utils.GetPathBySource(*p.Source)
		}
		var providersName = utils.GetNameBySource(*p.Source)
		plug, err := plugin.NewManagedPlugin(p.Path, providersName, p.Version, "", nil)
		if err != nil {
			hasError = true
			ui.PrintErrorF("%s@%s verification failed ：%s", providersName, p.Version, err.Error())
			continue
		}
		conf, err := tools.GetProviders(&s, p.Name)
		if err != nil {
			hasError = true
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
			hasError = true
			ui.PrintErrorF("%s@%s verification failed ：%s", providersName, p.Version, err.Error())
			continue
		} else {
			if initRes.Diagnostics != nil && initRes.Diagnostics.HasError() {
				ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
				hasError = true
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
			hasError = true
			continue
		} else {
			if res.Diagnostics != nil && res.Diagnostics.HasError() {
				ui.PrintDiagnostic(res.Diagnostics.GetDiagnosticSlice())
				hasError = true
				continue
			}
		}
		ui.PrintSuccessF("	%s@%s check successfully", providersName, p.Version)
	}

	ui.PrintSuccessF("\nProviders verification completed\n")
	ui.PrintSuccessF("Profile verification complete\n")
	if hasError {
		return errors.New("Need help? Know on Slack or open a Github Issue: https://github.com/selefra/selefra#community")
	}
	return nil
}
