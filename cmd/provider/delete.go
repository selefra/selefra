package provider

import (
	"encoding/json"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

func newCmdProviderRemove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove one or more plugins from the download cache",
		Long:  "Remove one or more plugins from the download cache",
		RunE: func(cmd *cobra.Command, names []string) error {
			argsMap := make(map[string]bool)
			for i := range names {
				argsMap[names[i]] = true
			}

			wd, err := os.Getwd()
			*global.WORKSPACE = wd
			var providers []*config.ProviderRequired

			var cof = &config.SelefraConfig{}

			namespace, _, err := utils.Home()
			if err != nil {
				return err
			}
			provider := registry.NewProviderRegistry(namespace)
			err = cof.GetConfig()
			if err != nil {
				return err
			}
			for _, p := range cof.Selefra.Providers {
				name := utils.GetNameBySource(*p.Source)
				path := utils.GetPathBySource(*p.Source)
				prov := registry.ProviderBinary{
					Provider: registry.Provider{
						Name:    name,
						Version: p.Version,
						Source:  "",
					},
					Filepath: path,
				}
				if !argsMap[p.Name] {
					providers = append(providers, p)
					break
				}

				err := provider.DeleteProvider(prov)
				if err != nil {
					return err
				}
				_, jsonPath, err := utils.Home()
				if err != nil {
					return err
				}
				c, err := os.ReadFile(jsonPath)
				if err != nil {
					return err
				}
				var configMap = make(map[string]string)
				err = json.Unmarshal(c, &configMap)
				if err != nil {
					return err
				}

				delete(configMap, *p.Source)

				c, err = json.Marshal(configMap)
				if err != nil {
					return err
				}
				err = os.Remove(jsonPath)
				if err != nil {
					return err
				}
				err = os.WriteFile(jsonPath, c, 0644)
				if err != nil {
					return err
				}

				for i := range cof.Selefra.Providers {
					if cof.Selefra.Providers[i].Name == p.Name {
						cof.Selefra.Providers = append(cof.Selefra.Providers[:i], cof.Selefra.Providers[i+1:]...)
					}
				}

				for i := 0; i < len(cof.Providers.Content); i++ {
					for ii := range cof.Providers.Content[i].Content {
						if cof.Providers.Content[i].Content[ii].Kind == yaml.ScalarNode && cof.Providers.Content[0].Content[i].Value == "name" && cof.Providers.Content[0].Content[i+1].Value == p.Name {
							if len(cof.Providers.Content) == 1 {
								cof.Providers.Content = nil
								i--
								break
							} else {
								cof.Providers.Content = append(cof.Providers.Content[:i], cof.Providers.Content[i+1:]...)
								i--
								break
							}
						}
					}
				}

				ui.PrintSuccessF("Removed %s success", *p.Source)

			}
			configPath, err := config.GetConfigPath()
			if err != nil {
				return err
			}
			str, err := yaml.Marshal(cof)
			if err != nil {
				return err
			}
			err = os.Remove(configPath)
			if err != nil {
				return err
			}
			return os.WriteFile(configPath, str, 0644)
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}
