package provider

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestUpdate(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	ctx := context.Background()
	arg := []string{"aws"}
	err := update(ctx, arg)
	if err != nil {
		t.Error(err)
	}
}
