package tools

import (
	"encoding/json"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"gopkg.in/yaml.v3"
	"os"
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

func SetSelefraProvider(provider registry.ProviderBinary, selefraConfig *config.SelefraConfig) error {
	source := utils.CreateSource(provider.Name, provider.Version)
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
