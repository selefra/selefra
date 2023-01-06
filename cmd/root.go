package cmd

import (
	"fmt"
	"github.com/selefra/selefra/cmd/apply"
	"github.com/selefra/selefra/cmd/fetch"
	initCmd "github.com/selefra/selefra/cmd/init"
	"github.com/selefra/selefra/cmd/login"
	"github.com/selefra/selefra/cmd/logout"
	"github.com/selefra/selefra/cmd/provider"
	"github.com/selefra/selefra/cmd/query"
	"github.com/selefra/selefra/cmd/test"
	"github.com/selefra/selefra/cmd/version"
	"github.com/selefra/selefra/global"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var group = make(map[string][]*cobra.Command)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "selefra",
	Short: "Selefra - Infrastructure data as Code",
	Long: `
Selefra - Infrastructure data as Code

For details see the selefra document https://selefra.io/docs
If you like selefra, give us a star https://github.com/selefra/selefra
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level, _ := cmd.Flags().GetString("loglevel")
		global.ChangeLevel(level)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error occurred in Execute: %+v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("loglevel", "l", "debug", "log level")
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")
	group["main"] = []*cobra.Command{
		initCmd.NewInitCmd(),
		test.NewTestCmd(),
		apply.NewApplyCmd(),
		login.NewLoginCmd(),
		logout.NewLogoutCmd(),
	}

	group["other"] = []*cobra.Command{
		fetch.NewFetchCmd(),
		provider.NewProviderCmd(),
		query.NewQueryCmd(),
		version.NewVersionCmd(),
	}

	rootCmd.AddCommand(group["main"]...)
	rootCmd.AddCommand(group["other"]...)

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.TrimSpace(cmd.Long))

		fmt.Println("\nUsage:")
		fmt.Printf("  %-13s", "selefra [command]\n\n")

		fmt.Println("Main commands:")
		for _, c := range group["main"] {
			fmt.Printf("  %-13s%s\n", c.Name(), c.Short)
		}
		fmt.Println()
		fmt.Println("All other commands:")
		for _, c := range group["other"] {
			fmt.Printf("  %-13s%s\n", c.Name(), c.Short)
		}
		fmt.Println()

		fmt.Println("Flags")
		fmt.Println(cmd.Flags().FlagUsages())

		fmt.Println(`Use "selefra [command] --help" for more information about a command.`)
	})

}
