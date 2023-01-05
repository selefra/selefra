package query

import (
	"context"
	"github.com/c-bata/go-prompt"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/client"
	"strings"
)

type QueryClient struct {
	Ctx     context.Context
	Client  *client.Client
	Tables  []prompt.Suggest
	Columns []prompt.Suggest
}

func NewQueryClient(ctx context.Context, c *client.Client) *QueryClient {
	tables := CreateTablesSuggest(ctx, c)
	columns := CreateColumnsSuggest(ctx, c)
	return &QueryClient{
		ctx, c, tables, columns,
	}
}

// if there are no spaces this is the first word
func (q *QueryClient) isFirstWord(text string) bool {
	return strings.LastIndex(text, " ") == -1
}

func (q *QueryClient) formatSuggest(text string, before string) []prompt.Suggest {
	var s []prompt.Suggest
	if q.isFirstWord(text) {
		if text != "" {
			s = []prompt.Suggest{
				{Text: "SELECT"},
				{Text: "WITH"},
			}
		}
	} else {
		texts := strings.Split(before, " ")
		if strings.ToLower(texts[len(texts)-2]) == "from" {
			s = q.Tables
		}
		if strings.ToLower(texts[len(texts)-2]) == "select" {
			s = q.Columns
		}
		if strings.ToLower(texts[len(texts)-2]) == "," {
			s = q.Columns
		}
	}
	return s
}

func (q *QueryClient) completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	s := q.formatSuggest(d.Text, text)
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func CreateTablesSuggest(ctx context.Context, c *client.Client) []prompt.Suggest {
	res, diag := c.Storage.Query(ctx, TABLESQL)
	tables := []prompt.Suggest{}
	if diag != nil {
		_ = ui.PrintDiagnostic(diag.GetDiagnosticSlice())
	} else {
		rows, diag := res.ReadRows(-1)
		if diag != nil {
			_ = ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		}
		for i := range rows.GetMatrix() {
			tableName := rows.GetMatrix()[i][0].(string)
			tables = append(tables, prompt.Suggest{Text: tableName})
		}
	}
	return tables
}

func CreateColumnsSuggest(ctx context.Context, c *client.Client) []prompt.Suggest {
	res, err := c.Storage.Query(ctx, COLUMNSQL)
	columns := []prompt.Suggest{}
	if err != nil {
		_ = ui.PrintDiagnostic(err.GetDiagnosticSlice())
	} else {
		rows, err := res.ReadRows(-1)
		if err != nil {
			_ = ui.PrintDiagnostic(err.GetDiagnosticSlice())
		}
		for i := range rows.GetMatrix() {
			schemaName := rows.GetMatrix()[i][0].(string)
			tableName := rows.GetMatrix()[i][1].(string)
			columnName := rows.GetMatrix()[i][2].(string)
			columns = append(columns, prompt.Suggest{Text: columnName, Description: schemaName + "." + tableName})
		}
	}
	return columns
}
