package provider

import (
	"fmt"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

func newCmdProviderList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List currently installed plugins",
		Long:  "List currently installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				ui.PrintErrorLn("Error:" + err.Error())
				return nil
			}
			*global.WORKSPACE = wd

			b, err := config.GetClientStr()
			if err != nil {
				ui.PrintErrorLn("Error:" + err.Error())
				return nil
			}
			var configYaml config.SelefraConfig
			err = yaml.Unmarshal(b, &configYaml)
			if err != nil {
				ui.PrintErrorLn("Error:" + err.Error())
				return nil
			}
			fmt.Printf("  %-13s %-26s %s\n", "Name", "Source", "Version")
			for _, provider := range configYaml.Selefra.Providers {
				fmt.Printf("  %-13s %-26s %s\n", provider.Name, *provider.Source, provider.Version)
			}
			return nil
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}
