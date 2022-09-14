package init

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	utils2 "github.com/selefra/selefra/cmd/utils"
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
	cmd.PersistentFlags().StringP("dir", "d", ".", "the directory to initialize in")
	cmd.PersistentFlags().BoolP("force", "f", false, "force overwriting the directory if it is not empty")

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func CreateYaml(ctx context.Context) error {
	configYaml := config.SelefraConfig{}
	configYaml.Selefra.CliVersion = version.Version
	configYaml.Selefra.Connection = &config.DB{
		Driver:   "",
		Type:     "postgres",
		Username: "postgres",
		Password: "pass",
		Host:     "localhost",
		Port:     "5432",
		Database: "postgres",
		SSLMode:  "disable",
		Extras:   nil,
	}

	storage := postgresql_storage.NewPostgresqlStorageOptions(configYaml.Selefra.GetDSN())

	_, diag := postgresql_storage.NewPostgresqlStorage(ctx, storage)
	if diag != nil && diag.HasError() {
		ui.PrintDiagnostic(diag.GetDiagnosticSlice())

		ui.PrintErrorLn("The database maybe not ready.")
		ui.PrintErrorLn("You can execute the following command to install the official database image.")
		ui.PrintErrorLn("docker run --name selefra_postgres -p 5432:5432 -e POSTGRES_PASSWORD=pass -d postgres\n")
		return nil
	}
	var prov []string
	ui.PrintInfoLn("Getting provider list...")
	req, _ := http.NewRequest("GET", "https://github.com/selefra/registry/file-list/main/provider", nil)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		ui.PrintErrorF("Error: %s", err.Error())
		return nil
	}
	d, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		ui.PrintErrorF("Error: %s", err.Error())
		return nil
	}
	d.Find(".js-navigation-open.Link--primary").Each(func(i int, s *goquery.Selection) {
		if s.Text() != "template" {
			prov = append(prov, s.Text())
		}
	})
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
	}
	provs := term.SelectProviders(prov)
	ui.PrintInfoLn("Initializing Selefra provider plugin...")
	namespace, _, err := utils.Home()
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
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
			ui.PrintErrorF("Installed %s@%s failed：%s", p.Name, p.Version, err.Error())
			return nil
		} else {
			ui.PrintSuccessF("Installed %s@%s verified", p.Name, p.Version)
		}
		ui.PrintInfoF("Synchronization %s@%s's config...", p.Name, p.Version)
		plug, err := plugin.NewManagedPlugin(p.Filepath, p.Name, p.Version, "", nil)
		if err != nil {
			ui.PrintErrorF("Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
			return nil
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
			ui.PrintErrorLn(err.Error())
			return nil
		}

		if initRes != nil && initRes.Diagnostics != nil && initRes.Diagnostics.HasError() {
			ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
			return nil
		}

		res, err := plugProvider.GetProviderInformation(ctx, &shard.GetProviderInformationRequest{})
		if err != nil {
			ui.PrintErrorF("Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
			return nil
		}
		ui.PrintSuccessF("Synchronization %s@%s's config successful", p.Name, p.Version)
		utils2.SetSelefraProvider(p, &configYaml)
		err = utils2.SetProviders(res.DefaultConfigTemplate, p, &configYaml)
		if err != nil {
			ui.PrintErrorF("set %s@%s's config failed：%s", p.Name, p.Version, err.Error())
			return nil
		}
	}
	str, err := yaml.Marshal(configYaml)
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return nil
	}

	rulePath := filepath.Join(*global.WORKSPACE, "rules")
	_, err = os.Stat(rulePath)
	if err != nil {
		if os.IsNotExist(err) {
			mkErr := os.Mkdir(rulePath, 0755)
			if mkErr != nil {
				return mkErr
			}
		}
	}

	err = os.WriteFile(filepath.Join(rulePath, "default.yaml"), []byte(strings.TrimSpace(ruleComment)), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(*global.WORKSPACE, "module.yaml"), []byte(strings.TrimSpace(moduleComment)), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(*global.WORKSPACE, "selefra.yaml"), str, 0644)

	ui.PrintSuccessF(`
Selefra has been successfully initialized! 
	
Configuration generated successfully to %s
	`, *global.WORKSPACE)

	return nil
}

func initFunc(cmd *cobra.Command, args []string) error {

	wd, err := os.Getwd()
	force, _ := cmd.PersistentFlags().GetBool("force")
	dirname, _ := cmd.PersistentFlags().GetString("dir")
	*global.WORKSPACE = filepath.Join(wd, dirname)

	_, err = os.Stat(*global.WORKSPACE)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(*global.WORKSPACE, 0755)
		if err != nil {
			return nil
		}
	}
	dir, _ := os.ReadDir(*global.WORKSPACE)
	if len(dir) != 0 && !force {
		errStr := fmt.Sprintf("%s is not empty; rerun in an empty directory, pass the path to an empty directory to --dir, or use --force\n", *global.WORKSPACE)
		return errors.New(errStr)
	}
	_, clientErr := config.GetClientStr()
	if !errors.Is(clientErr, config.NoClient) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("error:%s is already init. Continue and overwrite it?[Y/N]\n", *global.WORKSPACE)
		text, err := reader.ReadString('\n')
		if err != nil {
			return nil
		}
		if text == "Y\n" || text == "y\n" {
			err = CreateYaml(cmd.Context())
			if err != nil {
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
	CreateYaml(cmd.Context())
	return nil
}
