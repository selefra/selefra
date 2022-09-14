package provider

import (
	"fmt"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func newCmdProviderList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				ui.PrintErrorLn("Error:" + err.Error())
				return nil
			}
			dirname, _ := cmd.PersistentFlags().GetString("dir")
			*global.WORKSPACE = filepath.Join(wd, dirname)

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
			fmt.Printf("  %-13s%s\n", "Name", "Version")
			for _, provider := range configYaml.Selefra.Providers {
				fmt.Printf("  %-13s%s\n", provider.Name, provider.Version)
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringP("dir", "d", ".", "the directory to initialize in")
	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}
