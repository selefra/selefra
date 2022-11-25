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
	global.LOGINTOKEN = "fc656fd36296a4c2f61dee25aeedfd0f"
	*global.WORKSPACE = "../../tests/workspace/online"
	err := Apply(context.Background())
	if err != nil {
		t.Error(err)
	}
}
