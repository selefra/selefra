package provider

import (
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

func newCmdProviderDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete",
		Long:  "delete",
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
				prov := registry.ProviderBinary{
					Provider: registry.Provider{
						Name:    p.Name,
						Version: p.Version,
						Source:  "",
					},
					Filepath: p.Path,
				}
				if !argsMap[p.Name] {
					providers = append(providers, p)
					break
				}

				err := provider.DeleteProvider(prov)
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
