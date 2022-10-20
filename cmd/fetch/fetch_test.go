package fetch

import (
	"context"
	"github.com/selefra/selefra/config"
	"testing"
)

func TestFetch(t *testing.T) {
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
