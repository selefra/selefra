package grpcClient

import (
	"github.com/selefra/selefra/global"
	issue "github.com/selefra/selefra/pkg/grpcClient/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
)

var opts []grpc.DialOption

func getDial() string {
	var dialMap = make(map[string]string)
	dialMap["dev-api.selefra.io"] = "dev-tcp.selefra.io"
	dialMap["main-api.selefra.io"] = "main-tcp.selefra.io"
	dialMap["pre-api.selefra.io"] = "pre-tcp.selefra.io"
	if dialMap[global.SERVER] != "" {
		return dialMap[global.SERVER]
	}
	arr := strings.Split(global.SERVER, ":")
	return arr[0]
}

func InitConn() (issue.IssueClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(getDial()+":1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := issue.NewIssueClient(conn)
	return client, conn, err
}
