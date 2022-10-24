package apply

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/selefra/selefra/cmd/provider"
	"github.com/selefra/selefra/cmd/test"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/httpClient"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/pkg/ws"
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
	global.CMD = "apply"
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	*global.WORKSPACE = wd
	return Apply(cmd.Context())
}

func Apply(ctx context.Context) error {
	err := config.IsSelefra()
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}
	s := config.SelefraConfig{}
	err = s.GetConfig()
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}
	token, err := utils.GetCredentialsToken()
	if token != "" && s.Selefra.Cloud != nil && err == nil {
		if err != nil {
			ui.PrintErrorLn("The token is invalid. Please execute selefra to log out or log in again")
			return nil
		}
		if global.LOGINTOKEN == "" {
			global.LOGINTOKEN = token
		}
		_, err := httpClient.CreateProject(token, s.Selefra.Cloud.Project)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return nil
		}
		res, err := httpClient.CreateTask(token, s.Selefra.Cloud.Project)
		if err == nil {
			err := ws.Regis(token, res.Data.TaskId)
			if err != nil {
				ui.PrintWarningLn(err.Error())
			}
		}
	}
	uid, _ := uuid.NewUUID()
	global.STAG = "initializing"
	_, err = provider.Sync()
	if err != nil {
		if token != "" && s.Selefra.Cloud != nil && err == nil {
			_ = httpClient.SetupStag(token, s.Selefra.Cloud.Project, httpClient.Failed)
		}
		ui.PrintErrorLn(err.Error())
		return nil
	}
	err = test.CheckSelefraConfig(ctx, s)
	if err != nil {
		ui.PrintErrorLn(err.Error())
		if token != "" && s.Selefra.Cloud != nil && err == nil {
			_ = httpClient.SetupStag(token, s.Selefra.Cloud.Project, httpClient.Failed)
		}
		return nil
	}
	err = s.GetConfig()
	if err != nil {
		if token != "" && s.Selefra.Cloud != nil && err == nil {
			_ = httpClient.SetupStag(token, s.Selefra.Cloud.Project, httpClient.Failed)
		}
		ui.PrintErrorLn("Client creation error:" + err.Error())
		return nil
	}
	c, e := client.CreateClientFromConfig(ctx, &s.Selefra, uid)
	if e != nil {
		if token != "" && s.Selefra.Cloud != nil && err == nil {
			_ = httpClient.SetupStag(token, s.Selefra.Cloud.Project, httpClient.Failed)
		}
		ui.PrintErrorLn("Client creation error:" + e.Error())
		return nil
	}
	global.STAG = "infrastructure"
	modules, err := config.GetModulesByPath()
	if err != nil {
		if token != "" && s.Selefra.Cloud != nil && err == nil {
			err = httpClient.SetupStag(token, s.Selefra.Cloud.Project, httpClient.Failed)
		}
		ui.PrintErrorLn("Client creation error:" + err.Error())
		return nil
	}
	var mRules []config.Rule
	ui.PrintSuccessLn(`----------------------------------------------------------------------------------------------

Loading Selefra analysis code ...
`)
	if len(modules) == 0 {
		mRules = *RunRulesWithoutModule()
	} else {
		mRules = CreateRulesByModule(modules)
	}

	ui.PrintSuccessF("\n---------------------------------- Result for rules  ----------------------------------------\n")
	var project string
	if token != "" && s.Selefra.Cloud != nil {
		project = s.Selefra.Cloud.Project
	} else {
		project = ""
	}
	err = RunRules(ctx, c, project, mRules)
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
	}
	if token != "" && s.Selefra.Cloud != nil {
		err = UploadWorkspace(project)
		if err != nil {
			err = httpClient.SetupStag(token, project, httpClient.Failed)
			ui.PrintErrorLn(err.Error())
			return nil
		}
	}
	return nil
}

func UploadWorkspace(project string) error {
	fileMap, err := config.GetAllConfig(".", nil)
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}
	err = httpClient.UploadWorkplace(global.LOGINTOKEN, project, fileMap)
	if err != nil {
		ui.PrintErrorLn(err)
		return err
	}
	return nil
}

func RunRules(ctx context.Context, c *client.Client, project string, rules []config.Rule) error {
	var outputReq []httpClient.OutputReq
	for _, rule := range rules {
		var params = make(map[string]interface{})
		for key, input := range rule.Input {
			params[key] = input["default"]
		}
		ui.PrintSuccessF("%s - Rule \"%s\"\n", rule.Path, rule.Name)

		ui.PrintSuccessLn("Description:")
		desc, err := fmtTemplate(rule.Metadata.Description, params)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return err
		}
		ui.PrintSuccessLn("	" + desc)

		ui.PrintSuccessLn("Policy:")
		queryStr, err := fmtTemplate(rule.Query, params)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return err
		}
		ui.PrintSuccessLn("	" + queryStr)

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
		ui.PrintSuccessLn("Output")
		var outMetaData []httpClient.Metadata
		for _, row := range rows {
			var outPut = rule.Output
			var outMap = make(map[string]interface{})
			for index, value := range row {
				key := column[index]
				outMap[key] = value
			}
			out, err := fmtTemplate(outPut, outMap)
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return err
			}
			ResourceAccountId, _ := fmtTemplate(rule.Metadata.ResourceAccountId, outMap)
			ResourceId, _ := fmtTemplate(rule.Metadata.ResourceId, outMap)
			ResourceRegion, _ := fmtTemplate(rule.Metadata.ResourceRegion, outMap)
			var remediation string
			var remediationPath string
			if filepath.IsAbs(rule.Metadata.Remediation) {
				remediationPath = rule.Metadata.Remediation
			} else {
				remediationPath = filepath.Join(*global.WORKSPACE, rule.Metadata.Remediation)
			}
			remediationByte, err := os.ReadFile(remediationPath)
			remediation = string(remediationByte)
			if err != nil {
				remediation = err.Error()
			}
			outMetaData = append(outMetaData, httpClient.Metadata{
				Id:                rule.Metadata.Id,
				Severity:          rule.Metadata.Severity,
				ResourceType:      rule.Metadata.ResourceType,
				Remediation:       remediation,
				Provider:          rule.Metadata.Provider,
				ResourceAccountId: ResourceAccountId,
				ResourceId:        ResourceId,
				ResourceRegion:    ResourceRegion,
				Title:             rule.Metadata.Title,
				Description:       desc,
				Output:            out,
			})
			ui.PrintSuccessLn("	" + out)
		}

		if global.LOGINTOKEN != "" {

			var req httpClient.OutputReq
			req.Name = rule.Name
			req.Query = rule.Query
			req.Labels = rule.Labels
			req.Metadata = outMetaData
			outputReq = append(outputReq, req)
		}
	}
	if global.LOGINTOKEN != "" {
		err := httpClient.OutPut(global.LOGINTOKEN, project, outputReq)
		if err != nil {
			ui.PrintErrorLn(err)
		}
		err = ws.Completed()
		if err != nil {
			ui.PrintErrorLn(err)
		}
	}
	return nil
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

func RunRulesWithoutModule() *[]config.Rule {
	rules, _ := config.GetRules()
	for i := range rules.Rules {
		if strings.HasPrefix(rules.Rules[i].Query, ".") {
			sqlByte, err := os.ReadFile(filepath.Join(".", rules.Rules[i].Query))
			if err != nil {
				ui.PrintErrorF("sql open error:%s", err.Error())
				return nil
			}
			rules.Rules[i].Query = string(sqlByte)
		}
	}
	return &rules.Rules
}

func RunPathModule(module config.Module) *[]config.Rule {
	var resRule config.RulesConfig
	for _, use := range module.Uses {
		usePath := filepath.Join(*global.WORKSPACE, use)
		b, err := os.ReadFile(usePath)
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

		if err != nil {
			ui.PrintErrorLn(err.Error())
			return nil
		}
		var ruleConfig config.RulesConfig
		err = yaml.Unmarshal([]byte(string(b)), &ruleConfig)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return nil
		}
		for i := range ruleConfig.Rules {
			ruleConfig.Rules[i].Output = baseRule.Rules[i].Output
			ruleConfig.Rules[i].Query = baseRule.Rules[i].Query
			ruleConfig.Rules[i].Path = use
			_, err := os.Stat(filepath.Join(*global.WORKSPACE, ruleConfig.Rules[i].Query))
			if err == nil {
				var sqlPath string
				if filepath.IsAbs(ruleConfig.Rules[i].Query) {
					sqlPath = ruleConfig.Rules[i].Query
				} else {
					sqlPath = filepath.Join(*global.WORKSPACE, ruleConfig.Rules[i].Query)
				}
				sqlByte, err := os.ReadFile(sqlPath)
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
			ui.PrintSuccessF("	%s - Rule %s: loading ... ", use, baseRule.Rules[i].Name)
		}
		resRule.Rules = append(resRule.Rules, ruleConfig.Rules...)
	}

	return &resRule.Rules
}

func fmtTemplate(temp string, params map[string]interface{}) (string, error) {
	t, err := template.New("temp").Parse(temp)
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
