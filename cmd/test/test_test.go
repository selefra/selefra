package test

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestTestFunc(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	ctx := context.Background()
	err := TestFunc(ctx)
	if err != nil {
		t.Error(err)
	}
}
