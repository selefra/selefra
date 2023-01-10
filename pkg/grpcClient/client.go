package grpcClient

import (
	"context"
	"fmt"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/grpcClient/proto/issue"
	glog "github.com/selefra/selefra/pkg/grpcClient/proto/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

type grpcCli struct {
	Ctx                          context.Context
	conn                         *grpc.ClientConn
	issueUploadIssueStreamClient issue.Issue_UploadIssueStreamClient
	logUploadLogStreamClient     glog.Log_UploadLogStreamClient
	taskId                       string
	token                        string
	statusMap                    map[string]string
}

var Cli grpcCli

func (g *grpcCli) SetStatus(status string) {
	if g.statusMap[global.STAG] == "" {
		g.statusMap[global.STAG] = status
	}
}

func (g *grpcCli) getStatus() string {
	if g.statusMap[global.STAG] != "" {
		return g.statusMap[global.STAG]
	}
	return "success"
}

func (g *grpcCli) getDial() string {
	var dialMap = make(map[string]string)
	dialMap["dev-api.selefra.io"] = "selefra-cloud-api-svc.selefra-cloud-dev:1234"
	dialMap["main-api.selefra.io"] = "selefra-cloud-api-svc.selefra-cloud-main:1234"
	dialMap["pre-api.selefra.io"] = "selefra-cloud-api-svc.selefra-cloud-pre:1234"
	if dialMap[global.SERVER] != "" {
		return dialMap[global.SERVER]
	}
	arr := strings.Split(global.SERVER, ":")
	return arr[0] + ":1234"
}

func (g *grpcCli) NewConn(token, taskId string) error {
	conn, err := grpc.Dial(g.getDial(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("fail to dial: %v", err)
	}
	g.conn = conn
	g.taskId = taskId
	g.token = token
	g.statusMap = make(map[string]string)
	g.Ctx = context.Background()
	err = g.newLogClient()
	if err != nil {
		return fmt.Errorf("fail to create uploadLogStreamCli cli:%s", err.Error())
	}
	err = g.newIssueClient()
	if err != nil {
		return fmt.Errorf("fail to create uploadIssueCli cli:%s", err.Error())
	}
	return err
}

func (g *grpcCli) newLogClient() error {
	logCli := glog.NewLogClient(g.conn)
	uploadStreamCli, err := logCli.UploadLogStream(g.Ctx)
	if err != nil {
		return err
	}
	g.logUploadLogStreamClient = uploadStreamCli
	return nil
}

func (g *grpcCli) newIssueClient() error {
	issueCli := issue.NewIssueClient(g.conn)
	uploadIssueCli, err := issueCli.UploadIssueStream(g.Ctx)
	if err != nil {
		return err
	}
	g.issueUploadIssueStreamClient = uploadIssueCli
	return nil
}

func (g *grpcCli) GetIssueUploadIssueStreamClient() issue.Issue_UploadIssueStreamClient {
	return g.issueUploadIssueStreamClient
}

func (g *grpcCli) GetLogUploadLogStreamClient() glog.Log_UploadLogStreamClient {
	return g.logUploadLogStreamClient
}

func (g *grpcCli) GetTaskID() string {
	return g.taskId
}

func (g *grpcCli) GetToken() string {
	return g.token
}

func (g *grpcCli) GetConn() *grpc.ClientConn {
	return g.conn
}

func (g *grpcCli) UploadLogStatus() (*glog.Res, error) {
	if g.conn == nil {
		return nil, nil
	}
	logCli := glog.NewLogClient(g.conn)
	statusInfo := &glog.StatusInfo{
		BaseInfo: &glog.BaseConnectionInfo{
			Token:  g.GetToken(),
			TaskId: g.GetTaskID(),
		},
		Stag:   global.STAG,
		Status: g.getStatus(),
		Time:   timestamppb.Now(),
	}
	res, err := logCli.UploadLogStatus(g.Ctx, statusInfo)
	if err != nil {
		return nil, fmt.Errorf("Fail to upload log status:%s", err.Error())
	}
	return res, nil
}
