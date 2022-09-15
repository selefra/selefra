package apply

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/selefra/selefra/cmd/provider"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func NewApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Create or update infrastructure",
		Long:  "Create or update infrastructure",
		RunE:  applyFunc,
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func applyFunc(cmd *cobra.Command, args []string) error {

	wd, err := os.Getwd()
	*global.WORKSPACE = wd
	ctx := cmd.Context()
	uid, _ := uuid.NewUUID()
	err = provider.Sync()
	if err != nil {
		ui.PrintErrorLn("Client creation error:" + err.Error())
		return nil
	}
	s := config.SelefraConfig{}
	_, err = s.GetConfigWithViper()
	if err != nil {
		ui.PrintErrorLn("Client creation error:" + err.Error())
		return nil
	}
	c, e := client.CreateClientFromConfig(ctx, &s.Selefra, uid)
	if e != nil {
		ui.PrintErrorLn("Client creation error:" + e.Error())
		return nil
	}

	modules, err := config.GetModulesByPath()
	if err != nil {
		ui.PrintErrorLn("Client creation error:" + err.Error())
		return nil
	}

	mRules := CreateRulesByModule(modules)
	RunRules(ctx, c, mRules)
	return nil
}

func RunRules(ctx context.Context, c *client.Client, rules []config.Rule) {
	for _, rule := range rules {
		var params = make(map[string]interface{})
		for key, input := range rule.Input {
			params[key] = input["default"]
		}
		query := rule.Query

		queryStr, err := fmtTemplate(query, params)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return
		}
		ui.PrintWarningLn(queryStr)
		res, diag := c.Storage.Query(ctx, queryStr)
		if diag != nil && diag.HasError() {
			ui.PrintDiagnostic(diag.GetDiagnosticSlice())
			continue
		}
		table, diag := res.ReadRows(-1)
		if diag != nil && diag.HasError() {
			ui.PrintDiagnostic(diag.GetDiagnosticSlice())
			continue
		}
		column := table.GetColumnNames()
		rows := table.GetMatrix()
		for _, row := range rows {
			var outPut = rule.Output
			var outMap = make(map[string]interface{})
			for index, value := range row {
				key := column[index]
				outMap[key] = value
			}
			ui.PrintSuccessLn(fmtTemplate(outPut, outMap))
		}
	}
}

func CreateRulesByModule(modules []config.Module) []config.Rule {
	var rules []config.Rule
	for _, module := range modules {
		if rule := RunPathModule(module); rule != nil {
			rules = append(rules, *rule...)
		}
	}
	return rules
}

func RunPathModule(module config.Module) *[]config.Rule {
	b, err := os.ReadFile(module.Uses)
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
	}
	var baseRule config.RulesConfig
	err = yaml.Unmarshal(b, &baseRule)
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
	}
	fmtRuleConfigStr, err := fmtTemplate(string(b), module.Input)

	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
	}
	var ruleConfig config.RulesConfig
	err = yaml.Unmarshal([]byte(fmtRuleConfigStr), &ruleConfig)
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
	}
	for i := range ruleConfig.Rules {
		ruleConfig.Rules[i].Output = baseRule.Rules[i].Output
		ruleConfig.Rules[i].Query = baseRule.Rules[i].Query

		if strings.HasPrefix(ruleConfig.Rules[i].Query, ".") {
			dir := filepath.Dir(module.Uses)
			sqlByte, err := os.ReadFile(filepath.Join(dir, ruleConfig.Rules[i].Query))
			if err != nil {
				ui.PrintErrorF("sql open error:%s", err.Error())
				return nil
			}
			ruleConfig.Rules[i].Query = string(sqlByte)
		}
		for key, input := range ruleConfig.Rules[i].Input {
			if module.Input[key] != nil {
				input["default"] = module.Input[key]
			}
		}
	}

	return &ruleConfig.Rules
}

func fmtTemplate(temp string, params map[string]interface{}) (string, error) {
	t, err := template.New("test").Parse(temp)
	if err != nil {
		ui.PrintErrorLn("Format rule template error:", err.Error())
		return "", err
	}
	b := bytes.Buffer{}
	err = t.Execute(&b, params)
	if err != nil {
		ui.PrintErrorLn("Format rule template error:", err.Error())
		return "", err
	}
	by, err := io.ReadAll(&b)
	if err != nil {
		ui.PrintErrorLn("Format rule template error:", err.Error())
		return "", err
	}
	return string(by), nil
}
