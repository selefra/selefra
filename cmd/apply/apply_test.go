package apply

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestApply(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	err := Apply(context.Background())
	if err != nil {
		t.Error(err)
	}
}
