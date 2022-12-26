/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/selefra/selefra/cmd"
	"github.com/selefra/selefra/pkg/grpcClient"
	"github.com/selefra/selefra/ui"
	"log"
	"runtime/debug"
)

func main() {
	defer func() {
		if grpcClient.Cli.GetLogUploadLogStreamClient() != nil {
			err := grpcClient.Cli.GetIssueUploadIssueStreamClient().CloseSend()
			log.Fatalf("fail to close issue stream:%s", err.Error())
		}
		if grpcClient.Cli.GetConn() != nil {
			err := grpcClient.Cli.GetConn().Close()
			log.Fatalf("fail to close grpc conn:%s", err.Error())
		}
		if err := recover(); err != nil {
			ui.PrintErrorF("Panic: %v\n%s", err, debug.Stack())
		}
	}()
	cmd.Execute()
}
