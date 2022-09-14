package utils

import (
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/pkg/registry"
	"gopkg.in/yaml.v3"
)

func SetProviders(DefaultConfigTemplate string, provider registry.ProviderBinary, config *config.SelefraConfig) error {

	config.Providers.Kind = yaml.MappingNode
	config.Providers.HeadComment = "provider configurations"
	var node yaml.Node
	yaml.Unmarshal([]byte(DefaultConfigTemplate), &node)

	var provNode yaml.Node
	for _, Node := range config.Providers.Content {
		if Node.Kind == yaml.ScalarNode && Node.Value == provider.Name {
			return nil
		}
	}
	provNode.Content = append([]*yaml.Node{
		{
			Kind:  yaml.ScalarNode,
			Value: provider.Name,
		},
		{
			Kind:    yaml.MappingNode,
			Content: node.Content[0].Content,
		},
	})

	config.Providers.Content = append(config.Providers.Content, provNode.Content...)

	return nil
}
