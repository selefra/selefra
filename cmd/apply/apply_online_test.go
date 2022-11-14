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
	global.SERVER = "sfc.shubo6.cn:8443"
	global.LOGINTOKEN = "eff6532520a9919cd785669aa1c6d735"
	*global.WORKSPACE = "../../tests/workspace/online"
	err := Apply(context.Background())
	if err != nil {
		t.Error(err)
	}
}
