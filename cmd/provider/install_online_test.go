package provider

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestInstallOnline(t *testing.T) {
	global.SERVER = "dev-api.selefra.io"
	global.LOGINTOKEN = "4fe8ed36488c479d0ba7292fe09a4132"
	*global.WORKSPACE = "../../tests/workspace/online"
	ctx := context.Background()
	err := install(ctx, []string{"aws@v0.0.4"})
	if err != nil {
		t.Error(err)
	}
}
