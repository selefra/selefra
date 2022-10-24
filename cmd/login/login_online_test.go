package login

import (
	"github.com/spf13/cobra"
	"testing"
)

func TestRunFunc(t *testing.T) {
	var cmd *cobra.Command
	err := RunFunc(cmd, []string{"4fe8ed36488c479d0ba7292fe09a4132"})
	if err != nil {
		t.Error(err)
	}
}
