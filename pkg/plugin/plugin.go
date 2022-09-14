package plugin

import (
	"fmt"
	"github.com/hashicorp/go-plugin"
	"github.com/selefra/selefra-provider-sdk/grpc/serve"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra/pkg/logger"
	"os"
	"os/exec"
)

const (
	prefixManaged   = "managed"
	prefixUnmanaged = "unmanaged"
	defaultAlias    = "default"
)

type Plugin interface {
	Name() string
	Version() string
	ProtocolVersion() int
	Provider() shard.ProviderClient
	Close()
}

type pluginBase struct {
	name     string
	version  string
	client   *plugin.Client
	provider shard.ProviderClient
}

func (p pluginBase) Name() string {
	return p.name
}

func (p pluginBase) Provider() shard.ProviderClient {
	return p.provider
}

func (p pluginBase) Version() string {
	return p.version
}

type managedPlugin struct {
	pluginBase
}

func (m managedPlugin) ProtocolVersion() int {
	return m.client.NegotiatedVersion()
}

func (m managedPlugin) Close() {
	if m.client == nil {
		return
	}
	m.client.Kill()
}

type unmanagedPlugin struct {
	config *plugin.ReattachConfig
	pluginBase
}

func (u unmanagedPlugin) ProtocolVersion() int {
	return -1
}

func (u unmanagedPlugin) Close() {}

//type Plugins map[string]Plugin
//
//func (p Plugins) Get(alias string, name string, version string) Plugin {
//	alias = checkAlias(alias)
//
//	// 1. unmanagedPlugin
//	if v, ok := p[fmt.Sprintf(unmanagedFormat, alias, name, version)]; ok {
//		return v
//	}
//	// 2. managedPlugin
//	if v, ok := p[fmt.Sprintf(managedFormat, alias, name, version)]; ok {
//		return v
//	}
//	return nil
//}

func getProvider(name string, client *plugin.Client) (shard.ProviderClient, error) {
	grpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, err
	}
	raw, err := grpcClient.Dispense("provider")
	if err != nil {
		client.Kill()
		return nil, err
	}

	provider, ok := raw.(shard.ProviderClient)
	if !ok {
		client.Kill()
		return nil, fmt.Errorf("plugin %s is not a provider", name)
	}
	return provider, nil
}

func checkAlias(alias string) string {
	if alias == "" {
		return defaultAlias
	}
	return alias
}

func NewManagedPlugin(filepath string, name string, version string, alias string, env []string) (Plugin, error) {
	// managedFormat prefixManaged:alias:name:version (e.g. managed:alias:foo:1.0.0)
	managedFormat := fmt.Sprintf("%s:%%s:%%s:%%s", prefixManaged)

	defaultLogger, _ := logger.NewLogger(logger.Config{
		FileLogEnabled:    true,
		ConsoleLogEnabled: false,
		EncodeLogsAsJson:  true,
		ConsoleNoColor:    true,
		Source:            "plugin",
		Directory:         "logs",
		Level:             "debug",
	})

	alias = checkAlias(alias)
	cmd := exec.Command(filepath)
	cmd.Env = append(cmd.Env, env...)
	client := plugin.NewClient(&plugin.ClientConfig{
		SyncStdout:       os.Stdout,
		SyncStderr:       os.Stderr,
		HandshakeConfig:  serve.HandSharkConfig,
		VersionedPlugins: shard.VersionPluginMap,
		Managed:          true,
		Cmd:              cmd,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           defaultLogger,
	})
	provider, err := getProvider(name, client)
	if err != nil {
		return nil, err
	}

	return &managedPlugin{
		pluginBase: pluginBase{
			name:     fmt.Sprintf(managedFormat, alias, name, version),
			client:   client,
			provider: provider,
			version:  version,
		},
	}, nil
}

func NewUnmanagedPlugin(alias string, name string, version string, config *plugin.ReattachConfig) (Plugin, error) {
	alias = checkAlias(alias)
	// unmanagedFormat prefixUnmanaged:alias:name:version (e.g. unmanaged:alias:foo:1.0.0)
	unmanagedFormat := fmt.Sprintf("%s:%%s:%%s:%%s", prefixUnmanaged)
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  serve.HandSharkConfig,
		Plugins:          shard.PluginMap,
		Reattach:         config,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		//SyncStderr:       os.Stderr,
		//SyncStdout:       os.Stdout,
	})
	provider, err := getProvider(name, client)
	if err != nil {
		return nil, err
	}
	return &unmanagedPlugin{
		config: config,
		pluginBase: pluginBase{
			name:     fmt.Sprintf(unmanagedFormat, alias, name, version),
			client:   client,
			provider: provider,
			version:  version,
		},
	}, nil
}
