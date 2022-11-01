package oci

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"io"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func loadBar(doneFlag *bool) {
	go func() {
		dotLen := 0
		for *doneFlag {
			time.Sleep(1 * time.Second)
			dotLen++
			ui.PrintCustomizeFNotN(ui.InfoColor, "\rWaiting for DB to download %s", strings.Repeat(".", dotLen%6)+strings.Repeat(" ", 6-dotLen%6))
		}
	}()
}

func RunDB() error {
	const goos = runtime.GOOS
	doneFlag := true
	loadBar(&doneFlag)
	ref := global.PkgBasePath + goos + global.PkgTag
	ctx := context.Background()
	resolver := docker.NewResolver(docker.ResolverOptions{})
	tempDir, _ := utils.GetTempPath()
	_ = os.MkdirAll(tempDir, 0755)
	fileStore := content.NewFile(tempDir)
	_, err := os.Stat(tempDir + "/pgsql/bin")
	dataPath := tempDir + "/pgsql/data"
	ctlPath := tempDir + "/pgsql/bin/pg_ctl"
	initPath := tempDir + "/pgsql/bin/initdb"
	confPath := tempDir + "/pgsql/data/postgresql.conf"
	if goos == "windows" {
		ctlPath = tempDir + "/pgsql/bin/pg_ctl.exe"
		initPath = tempDir + "/pgsql/bin/initdb.exe"
	}

	if os.IsNotExist(err) {
		_, err := oras.Copy(ctx, resolver, ref, fileStore, tempDir)
		doneFlag = false
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		cmd := exec.Command(initPath, "-D", dataPath, "-U", "postgres")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf(err.Error() + ": " + stderr.String())
		}
		err = ChangePort(confPath, "15432")
		if err != nil {
			return fmt.Errorf(err.Error())
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

func ChangePort(filePath, port string) error {
	//读写方式打开文件
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("open file filed:%s", err)
	}
	//defer关闭文件
	defer file.Close()

	//读取文件内容到io中
	reader := bufio.NewReader(file)
	pos := int64(0)
	for {
		//读取每一行内容
		line, err := reader.ReadString('\n')
		if err != nil {
			//读到末尾
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				return fmt.Errorf("read file filed:%s", err.Error())
			}
		}
		//根据关键词覆盖当前行
		if strings.Contains(line, "#port = 5432") {
			bytes := []byte("port = " + port)
			file.WriteAt(bytes, pos)
		}
		//每一行读取完后记录位置
		pos += int64(len(line))
	}
	return nil
}
