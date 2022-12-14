package apply

import (
	"context"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/grpcClient"
	"log"
	"testing"
)

func TestApplyOnLine(t *testing.T) {
	defer func() {
		logCli := grpcClient.Cli.GetLogUploadLogStreamClient()
		conn := grpcClient.Cli.GetConn()
		if logCli != nil {
			err := logCli.CloseSend()
			if err != nil {
				log.Fatalf("fail to close log stream:%s", err.Error())
			}
		}
		if conn != nil {
			err := conn.Close()
			if err != nil {
				log.Fatalf("fail to close grpc conn:%s", err.Error())
			}
		}
	}()
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}
	global.LOGINTOKEN = "f4e13ad351c2cdb1e1de3b8c2cf83a32"
	*global.WORKSPACE = "../../tests/workspace/online"
	err := Apply(context.Background())
	if err != nil {
		t.Error(err)
	}
}
