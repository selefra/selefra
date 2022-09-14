package init

import (
	"encoding/json"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"os"
)

func SetSelefraProvider(provider registry.ProviderBinary, selefraConfig *config.SelefraConfig) {

	source := utils.CreateSource(provider.Name, provider.Version)

	_, configPath, err := utils.Home()
	if err != nil {
		ui.PrintErrorLn("SetSelefraProviderError: " + err.Error())
	}
	var pathMap = make(map[string]string)
	file, err := os.ReadFile(configPath)
	if err != nil {
		ui.PrintErrorLn("SetSelefraProviderError: " + err.Error())
	}

	json.Unmarshal(file, &pathMap)

	pathMap[source] = provider.Filepath

	pathMapJson, err := json.Marshal(pathMap)

	if err != nil {
		ui.PrintErrorLn("SetSelefraProviderError: " + err.Error())
	}

	err = os.WriteFile(configPath, pathMapJson, 0644)

	selefraConfig.Selefra.Providers = append(selefraConfig.Selefra.Providers, &config.ProviderRequired{
		Name:    provider.Name,
		Source:  &source,
		Version: provider.Version,
	})
}
