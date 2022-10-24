package provider

import (
	"github.com/selefra/selefra/global"
	"testing"
)

func TestSyncOnline(t *testing.T) {
	global.SERVER = "dev-api.selefra.io"
	global.LOGINTOKEN = "4fe8ed36488c479d0ba7292fe09a4132"
	*global.WORKSPACE = "../../tests/workspace/online"
	errLogs, err := Sync()
	if err != nil {
		t.Error(err)
	}
	if len(errLogs) != 0 {
		t.Error(errLogs)
	}
}
