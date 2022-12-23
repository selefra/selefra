package tools

import (
	"context"
	"encoding/json"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetProviders(config *config.SelefraConfig, key string) (string, error) {
	var seleferMap = make(map[string][]*yaml.Node)
	for _, group := range config.Providers.Content {
		for i, node := range group.Content {
			if node.Kind == yaml.ScalarNode && node.Value == "name" && group.Content[i+1].Value == key {
				seleferMap["providers"] = append(seleferMap["providers"], group)
			}
		}
	}
	b, err := yaml.Marshal(seleferMap)
	return string(b), err
}

func SetProviders(DefaultConfigTemplate string, provider registry.ProviderBinary, config *config.SelefraConfig) error {
	if config.Providers.Kind != yaml.SequenceNode {
		config.Providers.Kind = yaml.SequenceNode
		config.Providers.Tag = "!!seq"
		config.Providers.Value = ""
		config.Providers.Content = []*yaml.Node{}
	}

	var node yaml.Node

	yaml.Unmarshal([]byte(DefaultConfigTemplate), &node)
	var provNode yaml.Node
	if node.Content == nil {
		provNode.Content = []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Tag:  "!!map",
				Content: append([]*yaml.Node{
					{
						Kind:        yaml.ScalarNode,
						Value:       "name",
						FootComment: DefaultConfigTemplate,
					},
					{
						Kind:  yaml.ScalarNode,
						Value: provider.Name,
					},
				}),
			},
		}
	} else {
		provNode.Content = []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Tag:  "!!map",
				Content: append([]*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Value: "name",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: provider.Name,
					},
				}),
			},
		}
	}

	config.Providers.Content = append(config.Providers.Content, provNode.Content...)

	return nil
}

func SetSelefraProvider(provider registry.ProviderBinary, selefraConfig *config.SelefraConfig, configVersion string) error {
	source, latestSource := utils.CreateSource(provider.Name, provider.Version, configVersion)
	_, configPath, err := utils.Home()
	if err != nil {
		ui.PrintErrorLn("SetSelefraProviderError: " + err.Error())
		return err
	}
	var pathMap = make(map[string]string)
	file, err := os.ReadFile(configPath)
	if err != nil {
		ui.PrintErrorLn("SetSelefraProviderError: " + err.Error())
		return err
	}
	json.Unmarshal(file, &pathMap)
	if latestSource != "" {
		pathMap[latestSource] = provider.Filepath
	}
	pathMap[source] = provider.Filepath

	pathMapJson, err := json.Marshal(pathMap)

	if err != nil {
		ui.PrintErrorLn("SetSelefraProviderError: " + err.Error())
	}

	err = os.WriteFile(configPath, pathMapJson, 0644)
	if selefraConfig != nil {
		selefraConfig.Selefra.Providers = append(selefraConfig.Selefra.Providers, &config.ProviderRequired{
			Name:    provider.Name,
			Source:  &source,
			Version: provider.Version,
		})
	}
	return nil
}

func GetStoreValue(cof config.SelefraConfig, provider *config.ProviderRequired, key string) (string, error) {
	storageOpt := postgresql_storage.NewPostgresqlStorageOptions(cof.Selefra.GetDSN())
	storageOpt.SearchPath = config.GetSchemaKey(provider)
	store, diag := postgresql_storage.NewPostgresqlStorage(context.Background(), storageOpt)
	if diag != nil && diag.HasError() {
		err := ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		return "", err
	}
	stoLogger, err := ui.StoLogger()
	if err != nil {
		return "", err
	}
	meta := &schema.ClientMeta{ClientLogger: stoLogger}
	store.SetClientMeta(meta)
	t, diag := store.GetValue(context.Background(), key)
	if diag != nil && diag.HasError() {
		err := ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		return "", err
	}
	return t, nil
}

func SetStoreValue(cof config.SelefraConfig, provider *config.ProviderRequired, key, value string) error {
	storageOpt := postgresql_storage.NewPostgresqlStorageOptions(cof.Selefra.GetDSN())
	storageOpt.SearchPath = config.GetSchemaKey(provider)
	store, diag := postgresql_storage.NewPostgresqlStorage(context.Background(), storageOpt)
	if diag != nil && diag.HasError() {
		err := ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		return err
	}

	stoLogger, err := ui.StoLogger()
	if err != nil {
		return err
	}
	meta := &schema.ClientMeta{ClientLogger: stoLogger}
	store.SetClientMeta(meta)
	diag = store.SetKey(context.Background(), key, value)
	if diag != nil && diag.HasError() {
		err := ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		return err
	}
	return nil
}

func NeedFetch(required config.ProviderRequired, cof config.SelefraConfig) (bool, error) {
	requireKey := config.GetCacheKey()
	t, err := GetStoreValue(cof, &required, requireKey)
	if err != nil {
		return true, err
	}
	fetchTime, err := time.ParseInLocation(time.RFC3339, t, time.Local)
	if err != nil {
		return true, err
	}
	cp, err := cof.GetProvider(required.Name)
	if err != nil {
		return true, err
	}
	duration, err := parseDuration(cp.Cache)
	if err != nil || duration == 0 {
		return true, err
	}
	if time.Now().Sub(fetchTime) > duration {
		return true, nil
	}
	return false, nil
}

func Lock(ctx context.Context, provider config.ProviderRequired, cof config.SelefraConfig) error {
	requireKey := config.GetLockKey()
	var err error
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			Unlock(provider, cof)
			return ctx.Err()
		case <-ticker.C:
			err = SetStoreValue(cof, &provider, requireKey, time.Now().Format(time.RFC3339))
			if err != nil {
				return err
			}
		}
	}
}

func Unlock(provider config.ProviderRequired, cof config.SelefraConfig) error {
	requireKey := config.GetLockKey()
	return SetStoreValue(cof, &provider, requireKey, "")
}

func Locked(required config.ProviderRequired, cof config.SelefraConfig) (bool, error) {
	requireKey := config.GetLockKey()
	t, err := GetStoreValue(cof, &required, requireKey)
	if err != nil {
		return false, err
	}
	lockingTime, err := time.ParseInLocation(time.RFC3339, t, time.Local)
	if err != nil {
		return false, err
	}
	lockingStep, err := parseDuration("1m")
	if time.Now().Sub(lockingTime) < lockingStep {
		return true, nil
	}
	return false, nil
}

func parseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)
	dr, err := time.ParseDuration(d)
	if err == nil {
		return dr, nil
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")

		hour, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(hour)
		ndr, err := time.ParseDuration(d[index+1:])
		if err != nil {
			return dr, nil
		}
		return dr + ndr, nil
	}

	dv, err := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv), err
}
