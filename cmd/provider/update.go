package provider

import (
	"github.com/selefra/selefra/cmd/fetch"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

func newCmdProviderUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update",
		Long:  "update",
		RunE: func(cmd *cobra.Command, args []string) error {

			wd, err := os.Getwd()
			*global.WORKSPACE = wd
			argsMap := make(map[string]bool)
			for i := range args {
				argsMap[args[i]] = true
			}
			ctx := cmd.Context()
			var cof = &config.SelefraConfig{}
			err = cof.GetConfig()
			if err != nil {
				return err
			}
			namespace, _, err := utils.Home()
			if err != nil {
				return err
			}
			provider := registry.NewProviderRegistry(namespace)
			for _, p := range cof.Selefra.Providers {
				prov := registry.ProviderBinary{
					Provider: registry.Provider{
						Name:    p.Name,
						Version: p.Version,
						Source:  "",
					},
					Filepath: p.Path,
				}
				if len(args) != 0 && !argsMap[p.Name] {
					break
				}

				pp, err := provider.CheckUpdate(ctx, prov)
				if err != nil {
					return err
				}
				p.Path = pp.Filepath
				p.Version = pp.Version
				err = fetch.Fetch(ctx, cof, p)
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
