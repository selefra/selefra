package apply

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestApplyOnLine(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}
	global.LOGINTOKEN = "8ddf8931d70601c04e60f79995507851"
	*global.WORKSPACE = "../../tests/workspace/online"
	err := Apply(context.Background())
	if err != nil {
		t.Error(err)
	}
}
