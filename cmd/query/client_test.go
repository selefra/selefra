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

func createCtxAndClient() (context.Context, *client.Client, error) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	var cof = &config.SelefraConfig{}
	err := cof.GetConfig()
	if err != nil {
		ui.PrintErrorLn(err)
		return nil, nil, err
	}
	uid, _ := uuid.NewUUID()
	ctx := context.Background()
	c, e := client.CreateClientFromConfig(ctx, &cof.Selefra, uid)
	if e != nil {
		ui.PrintErrorLn(e)
		return nil, nil, e
	}
	return ctx, c, nil
}

func TestCreateColumnsSuggest(t *testing.T) {
	ctx, c, err := createCtxAndClient()
	if err != nil {
		t.Error(err)
	}
	columns := CreateColumnsSuggest(ctx, c)
	if columns == nil {
		t.Error("Columns is nil")
	}
}

func TestCreateTablesSuggest(t *testing.T) {
	ctx, c, err := createCtxAndClient()
	if err != nil {
		t.Error(err)
	}
	tables := CreateTablesSuggest(ctx, c)
	if tables == nil {
		t.Error("Tables is nil")
	}
}
