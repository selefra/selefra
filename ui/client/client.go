package client

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/selefra/selefra-provider-sdk/storage"
	postgres "github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-provider-sdk/storage_factory"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/ui"
)

type Client struct {
	//downloadProgress ui.Progress
	cfg           *config.Config
	Providers     registry.Providers
	Registry      interface{}
	PluginManager interface{}
	Storage       storage.Storage
	instanceId    uuid.UUID
}

func CreateClientFromConfig(ctx context.Context, cfg *config.Config, instanceId uuid.UUID, provider *config.ProviderRequired) (*Client, error) {

	hub := new(interface{})
	pm := new(interface{})

	c := &Client{
		Storage:       nil,
		cfg:           cfg,
		Registry:      hub,
		PluginManager: pm,
		instanceId:    instanceId,
	}
	if cfg.GetDSN() != "" {
		options := postgres.NewPostgresqlStorageOptions(cfg.GetDSN())
		schema := config.GetSchema(provider)
		options.SearchPath = schema
		sto, err := storage_factory.NewStorage(ctx, storage_factory.StorageTypePostgresql, options)
		if err != nil && err.HasError() {
			ui.PrintDiagnostic(err.GetDiagnosticSlice())
			return nil, errors.New("failed to create storage")
		}
		c.Storage = sto
	}
	c.Providers = registry.Providers{}
	for _, rp := range cfg.Providers {
		c.Providers.Set(registry.Provider{Name: rp.Name, Version: rp.Version})
	}

	return c, nil
}
