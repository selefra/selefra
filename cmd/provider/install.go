package provider

import (
	"encoding/json"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func newCmdProviderInstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install",
		Long:  "install",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			wd, err := os.Getwd()
			if err != nil {
				ui.PrintErrorLn("Error:" + err.Error())
				return nil
			}
			*global.WORKSPACE = wd
			namespace, _, err := utils.Home()
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return nil
			}

			var configYaml config.SelefraConfig
			configStr, err := config.GetClientStr()
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return nil
			}
			err = yaml.Unmarshal(configStr, &configYaml)
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return nil
			}
			provider := registry.NewProviderRegistry(namespace)
			for _, s := range args {
				splitArr := strings.Split(s, "@")
				var name string
				var version string
				if len(splitArr) > 1 {
					name = splitArr[0]
					version = splitArr[1]
				} else {
					name = splitArr[0]
					version = "latest"
				}
				pr := registry.Provider{
					Name:    name,
					Version: version,
					Source:  "",
				}
				p, err := provider.Download(ctx, pr, true)
				continueFlag := false
				for _, provider := range configYaml.Selefra.Providers {
					providerName := utils.GetNameBySource(*provider.Source)
					if strings.ToLower(providerName) == strings.ToLower(p.Name) && strings.ToLower(provider.Version) == strings.ToLower(p.Version) {
						continueFlag = true
						break
					}
				}
				if continueFlag {
					ui.PrintWarningLn(fmt.Sprintf("Provider %s@%s already installed", p.Name, p.Version))
					continue
				}
				if err != nil {
					ui.PrintErrorF("Installed %s@%s failed：%s", p.Name, p.Version, err.Error())
					return nil
				} else {
					ui.PrintSuccessF("Installed %s@%s verified", p.Name, p.Version)
				}
				ui.PrintInfoF("Synchronization %s@%s's config...", p.Name, p.Version)
				plug, err := plugin.NewManagedPlugin(p.Filepath, p.Name, p.Version, "", nil)
				if err != nil {
					ui.PrintErrorF("Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
					return nil
				}

				plugProvider := plug.Provider()
				storage := postgresql_storage.NewPostgresqlStorageOptions(configYaml.Selefra.GetDSN())
				opt, err := json.Marshal(storage)
				initRes, err := plugProvider.Init(ctx, &shard.ProviderInitRequest{
					Workspace: global.WORKSPACE,
					Storage: &shard.Storage{
						Type:           0,
						StorageOptions: opt,
					},
					IsInstallInit:  pointer.TruePointer(),
					ProviderConfig: pointer.ToStringPointer(""),
				})

				if err != nil {
					ui.PrintErrorLn(err.Error())
					return nil
				}

				if initRes != nil && initRes.Diagnostics != nil && initRes.Diagnostics.HasError() {
					ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
					return nil
				}

				res, err := plugProvider.GetProviderInformation(ctx, &shard.GetProviderInformationRequest{})
				if err != nil {
					ui.PrintErrorF("Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
					return nil
				}
				ui.PrintSuccessF("Synchronization %s@%s's config successful", p.Name, p.Version)
				tools.SetSelefraProvider(p, &configYaml)
				hasProvider := false
				for _, Node := range configYaml.Providers.Content {
					if Node.Kind == yaml.ScalarNode && Node.Value == p.Name {
						hasProvider = true
						break
					}
				}
				if !hasProvider {
					err = tools.SetProviders(res.DefaultConfigTemplate, p, &configYaml)
				}
				if err != nil {
					ui.PrintErrorF("set %s@%s's config failed：%s", p.Name, p.Version, err.Error())
					return nil
				}
			}
			str, err := yaml.Marshal(configYaml)
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return nil
			}
			path, err := config.GetConfigPath()
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return nil
			}
			err = os.WriteFile(path, str, 0644)
			return nil
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}
