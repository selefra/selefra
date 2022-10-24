package fetch

import (
	"context"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestFetchOnline(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/online"
	global.SERVER = "dev-api.selefra.io"
	global.LOGINTOKEN = "4fe8ed36488c479d0ba7292fe09a4132"
	ctx := context.Background()
	var cof = &config.SelefraConfig{}
	err := cof.GetConfig()
	for _, p := range cof.Selefra.Providers {
		err = Fetch(ctx, cof, p)
		if err != nil {
			t.Error(err)
		}
	}
}
