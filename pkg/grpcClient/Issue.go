package grpcClient

import (
	issue "github.com/selefra/selefra/pkg/grpcClient/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var opts []grpc.DialOption

func InitConn() (issue.IssueClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial("dev-tcp.selefra.io:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := issue.NewIssueClient(conn)
	return client, conn, err
}
