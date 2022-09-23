package login

import (
	"bufio"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/httpClient"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to selefra using token",
		Long:  "Login to selefra using token",
		RunE:  RunFunc,
	}
	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func RunFunc(cmd *cobra.Command, args []string) error {
	s := config.SelefraConfig{}
	err := s.GetConfig()
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}
	var token string
	if len(args) == 0 {
		credentials, err := utils.GetCredentialsPath()
		if err != nil {
			return err
		}
		ui.PrintCustomizeFNotN(ui.InfoColor, `
Selefra will login for login app.selefra.io  using your browser.
If login is successful, Terraform will store the token in plain text in
the following file for use by subsequent commands:
	%s

	Enter your access token from https://app.selefra.io/settings/access_tokens
	or hit <ENTER> to log in using your browser:`, credentials)
		reader := bufio.NewReader(os.Stdin)
		token, err = reader.ReadString('\n')
		if err != nil {
			return nil
		}
		token = strings.TrimSpace(strings.Replace(token, "\n", "", -1))
		if token == "" {
			ui.PrintErrorLn("No token provided")
			return nil
		}
	} else {
		token = args[0]
	}
	err = CliLogin(token)
	if err != nil {
		ui.PrintErrorLn("The token is invalid. Please execute selefra to log out or log in again")
		return err
	}
	return nil
}

func CliLogin(token string) error {
	res, err := httpClient.Login(token)
	if err != nil {
		return err
	}
	Success(res.Data.OrgName, res.Data.TokenName, token)
	return nil
}

func Success(orgName, userName, token string) {
	err := utils.SetCredentials(token)
	if global.LOGINTOKEN == "" {
		global.LOGINTOKEN = token
	}
	global.ORGNAME = orgName
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return
	}
	ui.PrintSuccessF(`
Retrieved token for user: %s.

Welcome to Selefra Cloud!

Logged in to selefra as %s (https://app.selefra.io/%s)`, userName, orgName, orgName)
}
