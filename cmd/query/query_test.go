package query

import (
	"context"
	"github.com/google/uuid"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/client"
	"testing"
)

func TestNewQueryClient(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	var cof = &config.SelefraConfig{}
	err := cof.GetConfig()
	if err != nil {
		ui.PrintErrorLn(err)
		return
	}
	uid, _ := uuid.NewUUID()
	ctx := context.Background()
	c, e := client.CreateClientFromConfig(ctx, &cof.Selefra, uid, nil, config.CliProviders{})
	if e != nil {
		ui.PrintErrorLn(e)
		return
	}

	queryClient := NewQueryClient(ctx, c)
	if queryClient == nil {
		t.Error("queryClient is nil")
	}
}

func TestNewQueryClientOnline(t *testing.T) {
	global.SERVER = "dev-api.selefra.io"
	global.LOGINTOKEN = "4fe8ed36488c479d0ba7292fe09a4132"
	*global.WORKSPACE = "../../tests/workspace/online"
	var cof = &config.SelefraConfig{}
	err := cof.GetConfig()
	if err != nil {
		ui.PrintErrorLn(err)
		return
	}
	uid, _ := uuid.NewUUID()
	ctx := context.Background()
	c, e := client.CreateClientFromConfig(ctx, &cof.Selefra, uid, nil, config.CliProviders{})
	if e != nil {
		ui.PrintErrorLn(e)
		return
	}

	queryClient := NewQueryClient(ctx, c)
	if queryClient == nil {
		t.Error("queryClient is nil")
	}
}
