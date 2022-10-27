package oci

import (
	"bytes"
	"context"
	"fmt"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
	"os"
	"os/exec"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func RunDB() error {
	ui.PrintInfoLn("Running DB ...")
	ref := "localhost:5001/postgre:latest"
	ctx := context.Background()
	resolver := docker.NewResolver(docker.ResolverOptions{})
	tempDir, _ := utils.GetTempPath()
	_ = os.MkdirAll(tempDir, 0755)
	fileStore := content.NewFile(tempDir)
	_, err := os.Stat(tempDir + "/pgsql")
	dataPath := tempDir + "/pgsql/data"
	ctlPath := tempDir + "/pgsql/bin/pg_ctl"
	initPath := tempDir + "/pgsql/bin/initdb"
	if os.IsNotExist(err) {
		_, err := oras.Copy(ctx, resolver, ref, fileStore, tempDir)
		check(err)
		cmd := exec.Command(initPath, "-D", dataPath, "-U", "postgres")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf(err.Error() + ": " + stderr.String())
		}
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ctlPath, "-D", dataPath, "stop")
	_ = cmd.Run()

	cmd = exec.Command(ctlPath, "-D", dataPath, "-l", tempDir+"/pgsql/logfile", "start")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(err.Error() + ": " + stderr.String())
	}
	ui.PrintErrorLn("Running DB Success")
	return nil
}
