package provider

import (
	"github.com/spf13/cobra"
)

func NewProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider [command]",
		Short: "provider",
		Long:  "provider",
	}

	cmd.AddCommand(newCmdProviderUpdate(), newCmdProviderDelete(), newCmdProviderList(), newCmdProviderInstall())

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}
