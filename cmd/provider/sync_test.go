package provider

import (
	"github.com/selefra/selefra/global"
	"testing"
)

func TestSync(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	errLogs, err := Sync()
	if err != nil {
		t.Error(err)
	}
	if len(errLogs) != 0 {
		t.Error(errLogs)
	}
}
