package init

import (
	"github.com/selefra/selefra/cmd/version"
	"github.com/selefra/selefra/ui"
)

func initHeaderOutput(providers []string) {
	for i := range providers {
		ui.PrintSuccessLn(providers[i] + " [✔]\n")
	}
	ui.PrintSuccessF(`Running with selefra-cli %s

	This command will walk you through creating a new Selefra project

	Enter a value or leave blank to accept the (default), and press <ENTER>.
	Press ^C at any time to quit.`, version.Version)
}
