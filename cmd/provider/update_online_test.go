package provider

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestUpdateOnline(t *testing.T) {
	global.SERVER = "dev-api.selefra.io"
	global.LOGINTOKEN = "4fe8ed36488c479d0ba7292fe09a4132"
	*global.WORKSPACE = "../../tests/workspace/online"
	ctx := context.Background()
	arg := []string{"aws"}
	err := update(ctx, arg)
	if err != nil {
		t.Error(err)
	}
}
