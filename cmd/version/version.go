package version

import (
	"fmt"
	"github.com/selefra/selefra/global"
	"github.com/spf13/cobra"
)

const Version = "{{version}}"

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print Selefra's version number",
		Long:  "Print Selefra's version number",
		Run: func(cmd *cobra.Command, args []string) {
			version()
		},
	}
	return cmd
}

func version() {
	global.CMD = "version"
	fmt.Println(Version)
}
