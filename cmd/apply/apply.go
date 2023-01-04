package apply

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cmd/provider"
	"github.com/selefra/selefra/cmd/test"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/grpcClient"
	"github.com/selefra/selefra/pkg/grpcClient/proto/issue"
	"github.com/selefra/selefra/pkg/httpClient"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/client"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
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
		taskRes, err := httpClient.CreateTask(token, s.Selefra.Cloud.Project)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return nil
		}
		err = grpcClient.Cli.NewConn(token, taskRes.Data.TaskUUID)
		if err != nil {
			ui.PrintErrorLn(err.Error())
		}
	}
	uid, _ := uuid.NewUUID()
	global.STAG = "initializing"

	err = test.CheckSelefraConfig(ctx, s)
	if err != nil {
		ui.PrintErrorLn(err.Error())
		if token != "" && s.Selefra.Cloud != nil && err == nil {
			_ = httpClient.SetupStag(token, s.Selefra.Cloud.Project, httpClient.Failed)
		}
		return nil
	}

	_, lockArr, err := provider.Sync()
	defer func() {
		for _, item := range lockArr {
			err := item.Storage.UnLock(context.Background(), item.SchemaKey, item.Uuid)
			if err != nil {
				ui.PrintErrorLn(err.Error())
			}
		}
	}()
	if err != nil {
		if token != "" && s.Selefra.Cloud != nil && err == nil {
			_ = httpClient.SetupStag(token, s.Selefra.Cloud.Project, httpClient.Failed)
		}
		ui.PrintErrorLn(err.Error())
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

	var project string
	if token != "" && s.Selefra.Cloud != nil {
		project = s.Selefra.Cloud.Project
	} else {
		project = ""
	}
	_, err = grpcClient.Cli.UploadLogStatus()
	if err != nil {
		ui.PrintErrorLn(err.Error())
	}
	global.STAG = "infrastructure"
	for i := range s.Selefra.Providers {
		confs, err := tools.GetProviders(&s, s.Selefra.Providers[i].Name)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return nil
		}
		for _, conf := range confs {
			var cp config.CliProviders
			err := yaml.Unmarshal([]byte(conf), &cp)
			if err != nil {
				ui.PrintErrorLn(err.Error())
				continue
			}
			c, e := client.CreateClientFromConfig(ctx, &s.Selefra, uid, s.Selefra.Providers[i], cp)
			if e != nil {
				if token != "" && s.Selefra.Cloud != nil && err == nil {
					_ = httpClient.SetupStag(token, s.Selefra.Cloud.Project, httpClient.Failed)
				}
				ui.PrintErrorLn("Client creation error:" + e.Error())
				return nil
			}
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

			schema := config.GetSchemaKey(s.Selefra.Providers[i], cp)
			err = RunRules(ctx, s, c, project, mRules, schema)
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return nil
			}
		}
	}
	if token != "" && s.Selefra.Cloud != nil {
		_, err = grpcClient.Cli.UploadLogStatus()
		if err != nil {
			ui.PrintErrorLn(err.Error())
		}
		err = UploadWorkspace(project)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			sErr := httpClient.SetupStag(token, project, httpClient.Failed)
			if sErr != nil {
				ui.PrintErrorLn(sErr.Error())
			}
			return nil
		}
	}
	return nil
}

func UploadWorkspace(project string) error {
	fileMap, err := config.GetAllConfig(*global.WORKSPACE, nil)
	if err != nil {
		return err
	}
	err = httpClient.UploadWorkplace(global.LOGINTOKEN, project, fileMap)
	if err != nil {
		return err
	}
	return nil
}

func getTableMap(tableMap map[string]bool, schemaTable []*schema.Table) {
	for i := range schemaTable {
		tableMap[schemaTable[i].TableName] = true
		if len(schemaTable[i].SubTables) > 0 {
			getTableMap(tableMap, schemaTable[i].SubTables)
		}
	}
}

func match(s string, whitelistWordSet map[string]bool) []string {
	var matchResultSet []string
	inWord := false
	lastIndex := 0
	for index, c := range s {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c >= '0' && c <= '9' {
			if !inWord {
				inWord = true
				lastIndex = index
			}
		} else {
			if inWord {
				word := s[lastIndex:index]
				if _, exists := whitelistWordSet[word]; exists {
					matchResultSet = append(matchResultSet, word)
				}
				inWord = false
			}
		}
	}
	return matchResultSet
}

func getSqlTables(sql string, tableMap map[string]bool) (tables []string) {
	nonStr := strings.Replace(sql, "\n", "", -1)
	tables = match(nonStr, tableMap)
	return tables
}

func UploadIssueFunc(ctx context.Context, IssueReq <-chan *issue.Req) {
	inClient := grpcClient.Cli.GetIssueUploadIssueStreamClient()
	for {
		select {
		case req := <-IssueReq:
			if inClient == nil {
				continue
			}
			err := inClient.Send(req)
			if err != nil {
				ui.PrintErrorF("send issue to server error: %s", err.Error())
				return
			}
		case <-ctx.Done():
			if inClient != nil {
				inClient.CloseSend()
			}
			return
		}
	}
}

func RunRules(ctx context.Context, s config.SelefraConfig, c *client.Client, project string, rules []config.Rule, schema string) error {
	issueCtx, issueCancel := context.WithCancel(context.Background())
	defer issueCancel()
	issueChan := make(chan *issue.Req, 100)
	go UploadIssueFunc(issueCtx, issueChan)
	for _, rule := range rules {
		var variablesMap = make(map[string]interface{})
		for i := range s.Variables {
			variablesMap[s.Variables[i].Key] = s.Variables[i].Default
		}
		queryStr, err := fmtTemplate(rule.Query, variablesMap)
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
		if len(rows) == 0 {
			continue
		}
		ui.PrintSuccessF("%s - Rule \"%s\"\n", rule.Path, rule.Name)
		ui.PrintSuccessLn("Schema:")
		ui.PrintSuccessLn(schema + "\n")
		ui.PrintSuccessLn("Description:")

		desc, err := fmtTemplate(rule.Metadata.Description, variablesMap)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return err
		}
		ui.PrintSuccessLn("	" + desc)

		ui.PrintSuccessLn("Policy:")
		schemaTables, schemaDiag := c.Storage.TableList(ctx, schema)
		if schemaDiag != nil {
			if schemaDiag.HasError() {
				return ui.PrintDiagnostic(schemaDiag.GetDiagnosticSlice())
			}
		}
		var tableMap = make(map[string]bool)
		getTableMap(tableMap, schemaTables)

		uploadTables := getSqlTables(queryStr, tableMap)
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return err
		}
		ui.PrintSuccessLn("	" + queryStr)

		ui.PrintSuccessLn("Output")
		for _, row := range rows {
			var outMetaData issue.Metadata
			var baseRow = make(map[string]interface{})
			var outPut = rule.Output
			var outMap = make(map[string]interface{})
			for index, value := range row {
				key := column[index]
				outMap[key] = value
			}
			baseRow = outMap
			baseRowStr, err := json.Marshal(baseRow)
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return err
			}
			var outByte bytes.Buffer
			err = json.Indent(&outByte, baseRowStr, "", "\t")
			out, err := fmtTemplate(outPut, outMap)
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return err
			}
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
			outMetaData = issue.Metadata{
				Id:           rule.Metadata.Id,
				Severity:     rule.Metadata.Severity,
				Remediation:  remediation,
				Tags:         rule.Metadata.Tags,
				SrcTableName: uploadTables,
				Provider:     rule.Metadata.Provider,
				Title:        rule.Metadata.Title,
				Author:       rule.Metadata.Author,
				Description:  desc,
				Output:       outByte.String(),
			}

			ui.PrintSuccessLn("	" + out)

			var outLabel = make(map[string]string)
			for key := range rule.Labels {
				switch rule.Labels[key].(type) {
				case string:
					outStr, _ := fmtTemplate(rule.Labels[key].(string), baseRow)
					outLabel[key] = outStr
				case []string:
					var out []string
					for _, v := range rule.Labels[key].([]string) {
						s, _ := fmtTemplate(v, baseRow)
						out = append(out, s)
					}
					outLabel[key] = strings.Join(out, ",")
				}
			}

			if global.LOGINTOKEN != "" {
				reqs := issue.Req{
					Name:        rule.Name,
					Query:       rule.Query,
					Metadata:    &outMetaData,
					Labels:      outLabel,
					ProjectName: project,
					TaskUUID:    grpcClient.Cli.GetTaskID(),
					Token:       grpcClient.Cli.GetToken(),
				}
				issueChan <- &reqs
			}
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
	var b []byte
	var err error
	for _, use := range module.Uses {
		var usePath string
		if path.IsAbs(use) || strings.Index(use, "://") > -1 {
			usePath = use
		} else {
			usePath = filepath.Join(*global.WORKSPACE, use)
		}
		if strings.Index(usePath, "://") > -1 {
			d := config.Downloader{Url: usePath}
			b, err = d.Get()
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return nil
			}
		} else {
			b, err = os.ReadFile(usePath)
			if err != nil {
				ui.PrintErrorLn(err.Error())
				return nil
			}
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
