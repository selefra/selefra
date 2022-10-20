package provider

import (
	"github.com/selefra/selefra/global"
	"testing"
)

func TestList(t *testing.T) {
	*global.WORKSPACE = "../../tests/workspace/offline"
	err := list()
	if err != nil {
		t.Error(err)
	}
}
