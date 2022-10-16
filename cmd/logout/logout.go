package logout

import (
	"github.com/selefra/selefra/pkg/httpClient"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout to selefra",
		Long:  "Logout to selefra",
		RunE:  RunFunc,
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func RunFunc(cmd *cobra.Command, args []string) error {
	path, err := utils.GetCredentialsPath()
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
	}
	token, err := utils.GetCredentialsToken()
	if token != "" && err == nil {
		err := httpClient.Logout(token)
		if err != nil {
			ui.PrintErrorLn("Logout error:" + err.Error())
			return nil
		}
		err = utils.SetCredentials("")
		if err != nil {
			ui.PrintErrorLn(err.Error())
		}
	}

	ui.PrintSuccessF(`Removing the stored credentials for app.selefra.io from the following file:
    %s

Success! Selefra has removed the stored Access Token for app.selefra.io.`, path)
	return nil
}
