package query

import (
	"github.com/c-bata/go-prompt"
	"github.com/google/uuid"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/client"
	"github.com/selefra/selefra/ui/table"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query infrastructure data from storage",
		Long:  "Query infrastructure data from storage",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ui.PrintWarningLn("Please select table.")
			var cof = &config.SelefraConfig{}
			wd, err := os.Getwd()
			*global.WORKSPACE = wd
			err = cof.GetConfig()
			if err != nil {
				ui.PrintErrorLn(err)
				return
			}
			uid, _ := uuid.NewUUID()
			c, e := client.CreateClientFromConfig(ctx, &cof.Selefra, uid, nil, config.CliProviders{})
			if e != nil {
				ui.PrintErrorLn(e)
				return
			}

			queryClient := NewQueryClient(ctx, c)
			p := prompt.New(func(in string) {
				strArr := strings.Split(in, "/")
				s := strArr[0]

				res, err := c.Storage.Query(ctx, s)
				if err != nil {
					ui.PrintErrorLn(err)
				} else {
					tables, e := res.ReadRows(-1)
					if e != nil && e.HasError() {
						return
					}
					header := tables.GetColumnNames()
					body := tables.GetMatrix()
					var tableBody [][]string
					for i := range body {
						var row []string
						for j := range body[i] {
							row = append(row, utils.Strava(body[i][j]))
						}
						tableBody = append(tableBody, row)
					}

					if len(strArr) > 1 && strArr[1] == "g" {
						table.ShowRows(header, tableBody, []string{}, true)
					} else {
						table.ShowTable(header, tableBody, []string{}, true)
					}

				}
				if s == "exit;" || s == ".exit" {
					os.Exit(0)
				}
			}, queryClient.completer,
				prompt.OptionTitle("Table"),
				prompt.OptionPrefix("> "),
				prompt.OptionAddKeyBind(prompt.KeyBind{
					Key: prompt.ControlC,
					Fn: func(buffer *prompt.Buffer) {
						os.Exit(0)
					},
				}),
			)
			p.Run()
		},
	}
	return cmd
}
