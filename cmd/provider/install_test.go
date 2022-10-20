package provider

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestInstall(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	ctx := context.Background()
	err := install(ctx, []string{"aws@v0.0.4"})
	if err != nil {
		t.Error(err)
	}
}
