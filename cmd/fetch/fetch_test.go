package fetch

import (
	"context"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestFetch(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
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
