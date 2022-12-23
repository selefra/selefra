package provider

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestSync(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	ctx := context.Background()
	errLogs, err := Sync(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(errLogs) != 0 {
		t.Error(errLogs)
	}
}
