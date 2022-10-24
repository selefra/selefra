package provider

import (
	"github.com/selefra/selefra/global"
	"testing"
)

func TestListOnline(t *testing.T) {
	global.SERVER = "dev-api.selefra.io"
	global.LOGINTOKEN = "4fe8ed36488c479d0ba7292fe09a4132"
	*global.WORKSPACE = "../../tests/workspace/online"
	err := list()
	if err != nil {
		t.Error(err)
	}
}
