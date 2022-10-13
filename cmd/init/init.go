package init

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/selefra/selefra/cmd/login"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/pkg/httpClient"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"

	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cmd/version"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/term"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Prepare your working directory for other commands",
		Long:  "Prepare your working directory for other commands",
		RunE:  initFunc,
	}
	cmd.PersistentFlags().BoolP("force", "f", false, "force overwriting the directory if it is not empty")
	cmd.PersistentFlags().StringP("relevance", "r", "", "force overwriting the directory if it is not empty")

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func CreateYaml(cmd *cobra.Command) (*config.SelefraConfig, error) {
	ctx := cmd.Context()
	configYaml := config.SelefraConfig{}
	configYaml.Selefra.CliVersion = version.Version
	storage := postgresql_storage.NewPostgresqlStorageOptions(configYaml.Selefra.GetDSN())

	_, diag := postgresql_storage.NewPostgresqlStorage(ctx, storage)
	if diag != nil && diag.HasError() {
		ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		return &configYaml, errors.New(`The database maybe not ready.
		You can execute the following command to install the official database image.
		docker run --name selefra_postgres -p 5432:5432 -e POSTGRES_PASSWORD=pass -d postgres\n`)
	}
	var prov []string
	ui.PrintInfoLn("Getting provider list...")
	req, _ := http.NewRequest("GET", "https://github.com/selefra/registry/file-list/main/provider", nil)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		ui.PrintErrorF("Error: %s", err.Error())
		return &configYaml, err
	}
	d, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		ui.PrintErrorF("Error: %s", err.Error())
		return &configYaml, err
	}
	d.Find(".js-navigation-open.Link--primary").Each(func(i int, s *goquery.Selection) {
		if s.Text() != "template" {
			prov = append(prov, s.Text())
		}
	})
	if err != nil {
		return &configYaml, err
	}
	provs := term.SelectProviders(prov)
	if len(provs) == 0 {
		return &configYaml, errors.New("No provider selected or user canceled.")
	}
	initHeaderOutput(provs)
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("project name:(%s)", filepath.Base(*global.WORKSPACE))

	projectName, err := reader.ReadString('\n')
	if err != nil {
		return &configYaml, err
	}
	projectName = strings.TrimSpace(strings.Replace(projectName, "\n", "", -1))
	if projectName == "" {
		projectName = filepath.Base(*global.WORKSPACE)
	}

	token, err := utils.GetCredentialsToken()
	if token != "" && err == nil {
		configYaml.Selefra.Cloud = new(config.Cloud)
		relevance, _ := cmd.PersistentFlags().GetString("relevance")
		if relevance != "" {
			configYaml.Selefra.Cloud.Project = relevance
		} else {
			configYaml.Selefra.Cloud.Project = projectName
		}
		orgName, err := httpClient.CreateProject(token, configYaml.Selefra.Cloud.Project)
		configYaml.Selefra.Cloud.Organization = orgName
		if err != nil {
			return &configYaml, err
		}
		err = httpClient.SetupStag(token, configYaml.Selefra.Cloud.Project, httpClient.Creating)
		if err != nil {
			return &configYaml, err
		}
	}
	configYaml.Selefra.Name = projectName
	ui.PrintInfoLn("Initializing Selefra provider plugin...")

	namespace, _, err := utils.Home()
	if err != nil {
		return &configYaml, err
	}
	provider := registry.NewProviderRegistry(namespace)

	for _, s := range provs {
		pr := registry.Provider{
			Name:    s,
			Version: "latest",
			Source:  "",
		}
		p, err := provider.Download(ctx, pr, true)
		if err != nil {
			return &configYaml, fmt.Errorf("	Installed %s@%s failed：%s", p.Name, p.Version, err.Error())
		} else {
			if token != "" && err == nil {
				_ = httpClient.SetupStag(token, configYaml.Selefra.Cloud.Project, httpClient.Failed)
			}
			ui.PrintSuccessF("	Installed %s@%s verified", p.Name, p.Version)
		}
		ui.PrintInfoF("	Synchronization %s@%s's config...", p.Name, p.Version)
		plug, err := plugin.NewManagedPlugin(p.Filepath, p.Name, p.Version, "", nil)
		if err != nil {
			return &configYaml, fmt.Errorf("	Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
		}

		plugProvider := plug.Provider()
		opt, err := json.Marshal(storage)
		initRes, err := plugProvider.Init(ctx, &shard.ProviderInitRequest{
			Workspace: global.WORKSPACE,
			Storage: &shard.Storage{
				Type:           0,
				StorageOptions: opt,
			},
			IsInstallInit:  pointer.TruePointer(),
			ProviderConfig: pointer.ToStringPointer(""),
		})

		if err != nil {
			return &configYaml, err
		}
		if initRes != nil && initRes.Diagnostics != nil && initRes.Diagnostics.HasError() {
			return &configYaml, ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
		}

		res, err := plugProvider.GetProviderInformation(ctx, &shard.GetProviderInformationRequest{})
		if err != nil {
			return &configYaml, fmt.Errorf("	Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
		}
		ui.PrintSuccessF("	Synchronization %s@%s's config successful", p.Name, p.Version)
		tools.SetSelefraProvider(p, &configYaml)
		err = tools.SetProviders(res.DefaultConfigTemplate, p, &configYaml)
		if err != nil {
			return &configYaml, fmt.Errorf("set %s@%s's config failed：%s", p.Name, p.Version, err.Error())
		}
	}
	waitStr, err := yaml.Marshal(configYaml)
	if err != nil {
		return &configYaml, err
	}
	var str []byte
	if token != "" {
		var initConfigYaml config.SelefraConfigInitWithLogin
		err = yaml.Unmarshal(waitStr, &initConfigYaml)
		if err != nil {
			return &configYaml, err
		}

		str, err = yaml.Marshal(initConfigYaml)
	} else {
		var initConfigYaml config.SelefraConfigInit
		err = yaml.Unmarshal(waitStr, &initConfigYaml)
		if err != nil {
			return &configYaml, err
		}

		str, err = yaml.Marshal(initConfigYaml)
	}

	if err != nil {
		return &configYaml, err
	}

	rulePath := filepath.Join(*global.WORKSPACE, "rules")
	_, err = os.Stat(rulePath)
	if err != nil {
		if os.IsNotExist(err) {
			mkErr := os.Mkdir(rulePath, 0755)
			if mkErr != nil {
				return &configYaml, mkErr
			}
		}
	}

	err = os.WriteFile(filepath.Join(rulePath, "iam_mfa.yaml"), []byte(strings.TrimSpace(ruleComment)), 0644)
	if err != nil {
		return &configYaml, err
	}

	err = os.WriteFile(filepath.Join(*global.WORKSPACE, "module.yaml"), []byte(strings.TrimSpace(moduleComment)), 0644)
	if err != nil {
		return &configYaml, err
	}
	err = os.WriteFile(filepath.Join(*global.WORKSPACE, "selefra.yaml"), str, 0644)

	ui.PrintSuccessF(`
Selefra has been successfully initialized! 
	
Your new Selefra project "%s" was created!

To perform an initial analysis, run selefra apply
	`, projectName)

	return &configYaml, nil
}

func initFunc(cmd *cobra.Command, args []string) error {
	global.CMD = "init"
	wd, err := os.Getwd()
	force, _ := cmd.PersistentFlags().GetBool("force")
	dirname := "."
	if len(args) > 0 {
		dirname = args[0]
	}
	*global.WORKSPACE = filepath.Join(wd, dirname)

	token, err := utils.GetCredentialsToken()
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
	}
	if token != "" {
		err := login.CliLogin(token)
		if err != nil {
			ui.PrintErrorLn("The token is invalid. Please execute selefra to log out or log in again")
			return nil
		}
	}
	relevance, _ := cmd.PersistentFlags().GetString("relevance")
	if token == "" && relevance != "" {
		err := login.RunFunc(cmd, []string{})
		if err != nil {
			return nil
		}
	}

	_, err = os.Stat(*global.WORKSPACE)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(*global.WORKSPACE, 0755)
		if err != nil {
			return nil
		}
	}
	dir, _ := os.ReadDir(*global.WORKSPACE)
	if len(dir) != 0 && !force {
		return fmt.Errorf("%s is not empty; Rerun in an empty directory, or use -- force/-f to force overwriting in the current directory\n", *global.WORKSPACE)
	}
	_, clientErr := config.GetClientStr()
	if !errors.Is(clientErr, config.NoClient) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Error:%s is already init. Continue and overwrite it?[Y/N]\n", *global.WORKSPACE)
		text, err := reader.ReadString('\n')
		text = strings.TrimSpace(strings.ToLower(text))
		if err != nil {
			return nil
		}
		if text == "y" {
			cof, err := CreateYaml(cmd)
			if err != nil {
				if global.LOGINTOKEN != "" && err == nil {
					_ = httpClient.SetupStag(global.LOGINTOKEN, cof.Selefra.Cloud.Project, httpClient.Failed)
				}
				ui.PrintErrorLn(err.Error())
			}
			return nil
		}
		return errors.New("config file already exists")
	}
	if err != nil {
		ui.PrintWarningF("Error: %s\n", err.Error())
		return nil
	}
	cof, err := CreateYaml(cmd)
	if err != nil {
		if global.LOGINTOKEN != "" && err == nil {
			_ = httpClient.SetupStag(global.LOGINTOKEN, cof.Selefra.Cloud.Project, httpClient.Failed)
		}
		ui.PrintErrorLn("Error:", err.Error())
	}
	return nil
}
