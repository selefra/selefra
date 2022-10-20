package provider

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestRemove(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	err := Remove([]string{"aws"})
	if err != nil {
		t.Error(err)
	}
	err = install(context.Background(), []string{"aws@v0.0.4"})
	if err != nil {
		t.Error(err)
	}
}
