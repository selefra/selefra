package provider

import (
	"context"
	"github.com/selefra/selefra/cmd/fetch"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"os"
)

func newCmdProviderUpdate() *cobra.Command {
	global.CMD = "provider update"
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update one or more plugins",
		Long:  "Update one or more plugins",
		RunE: func(cmd *cobra.Command, args []string) error {

			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			*global.WORKSPACE = wd
			err = update(cmd.Context(), args)
			return nil
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func update(ctx context.Context, args []string) error {
	err := config.IsSelefra()
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return err
	}
	argsMap := make(map[string]bool)
	for i := range args {
		argsMap[args[i]] = true
	}
	var cof = &config.SelefraConfig{}
	err = cof.GetConfig()
	if err != nil {
		return err
	}
	namespace, _, err := utils.Home()
	if err != nil {
		return err
	}
	provider := registry.NewProviderRegistry(namespace)
	for _, p := range cof.Selefra.Providers {
		prov := registry.ProviderBinary{
			Provider: registry.Provider{
				Name:    p.Name,
				Version: p.Version,
				Source:  "",
			},
			Filepath: p.Path,
		}
		if len(args) != 0 && !argsMap[p.Name] {
			break
		}

		pp, err := provider.CheckUpdate(ctx, prov)
		if err != nil {
			return err
		}
		p.Path = pp.Filepath
		p.Version = pp.Version
		confs, err := tools.GetProviders(cof, p.Name)
		if err != nil {
			return err
		}
		for _, c := range confs {
			err = fetch.Fetch(ctx, cof, p, c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
