package provider

import (
	"github.com/spf13/cobra"
)

func NewProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider [command]",
		Short: "Top-level command to interact with providers",
		Long:  "Top-level command to interact with providers",
	}

	cmd.AddCommand(newCmdProviderUpdate(), newCmdProviderRemove(), newCmdProviderRemove(), newCmdProviderList(), newCmdProviderInstall())

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}
