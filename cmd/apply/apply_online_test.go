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
	global.LOGINTOKEN = "e21f9cfb7dd3ae3e85a2c96a04062a4a"
	*global.WORKSPACE = "../../tests/workspace/online"
	err := Apply(context.Background())
	if err != nil {
		t.Error(err)
	}
}
