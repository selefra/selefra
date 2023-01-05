package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/selefra/selefra/pkg/grpcClient"
	"github.com/selefra/selefra/pkg/grpcClient/proto/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	hclog "github.com/hashicorp/go-hclog"

	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/logger"
)

var levelMap = map[string]int{
	"trace":   0,
	"debug":   1,
	"info":    2,
	"warning": 3,
	"error":   4,
	"fatal":   5,
}

var levelColor = []*color.Color{
	InfoColor,
	InfoColor,
	InfoColor,
	WarningColor,
	ErrorColor,
	ErrorColor,
}

var step int32 = 0

var defaultLogger, _ = logger.NewLogger(logger.Config{
	FileLogEnabled:    true,
	ConsoleLogEnabled: false,
	EncodeLogsAsJson:  true,
	ConsoleNoColor:    true,
	Source:            "client",
	Directory:         "logs",
	Level:             "info",
})

func StoLogger() (*logger.StoLogger, error) {
	return logger.NewStoLogger(logger.Config{
		FileLogEnabled:    true,
		ConsoleLogEnabled: false,
		EncodeLogsAsJson:  true,
		ConsoleNoColor:    true,
		Source:            "client",
		Directory:         "logs",
		Level:             "info",
	})
}

var wsLogger *os.File

func init() {
	flag := strings.ToLower(os.Getenv("SELEFRA_CLOUD_FLAG"))
	if flag == "true" || flag == "enable" {
		_, err := os.Stat("ws.log")
		if err != nil {
			if !os.IsNotExist(err) {
				panic("Unknown error," + err.Error())
			}
			wsLogger, err = os.Create("ws.log")
		} else {
			wsLogger, err = os.OpenFile("ws.log", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		}
		if err != nil {
			panic("ws log file open error," + err.Error())
		}
	}
}

const (
	prefixManaged   = "managed"
	prefixUnmanaged = "unmanaged"
	defaultAlias    = "default"
)

var (
	ErrorColor   = color.New(color.FgRed, color.Bold)
	WarningColor = color.New(color.FgYellow, color.Bold)
	InfoColor    = color.New(color.FgWhite, color.Bold)
	SuccessColor = color.New(color.FgGreen, color.Bold)
)

type LogJOSN struct {
	Cmd   string    `json:"cmd"`
	Stag  string    `json:"stag"`
	Msg   string    `json:"msg"`
	Time  time.Time `json:"time"`
	Level string    `json:"level"`
}

func getLevel(c *color.Color) string {
	var level string
	switch c {
	case ErrorColor:
		level = "error"
	case WarningColor:
		level = "warn"
	case InfoColor:
		level = "info"
	case SuccessColor:
		level = "success"
	default:
	}
	return level
}

func createLog(msg string, c *color.Color) string {
	l := LogJOSN{
		Cmd:   global.CMD,
		Stag:  global.STAG,
		Msg:   msg,
		Time:  time.Now(),
		Level: getLevel(c),
	}
	b, err := json.Marshal(l)
	if err != nil {
		return ""
	}
	sb := string(b)
	if wsLogger != nil {
		_, _ = wsLogger.WriteString(sb + "\n")
	}
	return sb
}

func PrintErrorF(format string, a ...interface{}) {
	_, f, l, ok := runtime.Caller(1)
	if ok {
		if defaultLogger != nil {
			defaultLogger.Log(hclog.Error, "%s %s:%d", fmt.Sprintf(format, a...), f, l)
		}
	}
	PrintCustomizeF(ErrorColor, format, a...)
}

func PrintWarningF(format string, a ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Log(hclog.Warn, format, a...)
	}
	PrintCustomizeF(WarningColor, format, a...)
}

func PrintSuccessF(format string, a ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Log(hclog.Info, format, a...)
	}
	PrintCustomizeF(SuccessColor, format, a...)

}

func PrintInfoF(format string, a ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Log(hclog.Info, format, a...)
	}
	PrintCustomizeF(InfoColor, format, a...)
}

func PrintErrorLn(a ...interface{}) {
	_, f, l, ok := runtime.Caller(1)
	if ok {
		if defaultLogger != nil {
			defaultLogger.Log(hclog.Error, "%s %s:%d", fmt.Sprintln(a...), f, l)
		}
	}
	PrintCustomizeLn(ErrorColor, a...)
}

func PrintWarningLn(a ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Log(hclog.Warn, fmt.Sprintln(a...))
	}
	PrintCustomizeLn(WarningColor, a...)

}

func PrintSuccessLn(a ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Log(hclog.Info, fmt.Sprintln(a...))
	}
	PrintCustomizeLn(SuccessColor, a...)
}

func PrintInfoLn(a ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Log(hclog.Info, fmt.Sprintln(a...))
	}
	PrintCustomizeLn(InfoColor, a...)
}

func sendMsg(c *color.Color, logCli log.Log_UploadLogStreamClient, msg string) error {
	step++
	createLog(msg, c)
	if c == ErrorColor {
		grpcClient.Cli.SetStatus("error")
	}
	err := logCli.Send(&log.ConnectMsg{
		ActionName: "",
		Data: &log.LogJOSN{
			Cmd:   global.CMD,
			Stag:  global.STAG,
			Msg:   msg,
			Time:  timestamppb.Now(),
			Level: getLevel(c),
		},
		Index: step,
		Msg:   "",
		BaseInfo: &log.BaseConnectionInfo{
			Token:  grpcClient.Cli.GetToken(),
			TaskId: grpcClient.Cli.GetTaskID(),
		},
	})
	return err
}

func PrintCustomizeF(c *color.Color, format string, a ...interface{}) {
	logCli := grpcClient.Cli.GetLogUploadLogStreamClient()
	if logCli != nil {
		msg := fmt.Sprintf(format, a...)
		err := sendMsg(c, logCli, msg)
		if err != nil {
			createLog("grpc logStream error:"+err.Error(), ErrorColor)
		}
	}
	_, _ = c.Printf(format+"\n", a...)
}

func PrintCustomizeFNotN(c *color.Color, format string, a ...interface{}) {
	logCli := grpcClient.Cli.GetLogUploadLogStreamClient()
	if logCli != nil {
		msg := fmt.Sprintf(format, a...)
		err := sendMsg(c, logCli, msg)
		if err != nil {
			createLog("grpc logStream error:"+err.Error(), ErrorColor)
		}
	}
	_, _ = c.Printf(format, a...)
}

func PrintCustomizeLn(c *color.Color, a ...interface{}) {
	logCli := grpcClient.Cli.GetLogUploadLogStreamClient()
	if logCli != nil {
		str := fmt.Sprint(a...)
		err := sendMsg(c, logCli, str)
		if err != nil {
			createLog("grpc logStream error:"+err.Error(), ErrorColor)
		}
	}
	_, _ = c.Println(a...)
}

func PrintCustomizeLnNotShow(a string) {
	logCli := grpcClient.Cli.GetLogUploadLogStreamClient()
	if logCli != nil {
		err := sendMsg(InfoColor, logCli, a)
		if err != nil {
			createLog("grpc logStream error:"+err.Error(), ErrorColor)
		}
	}
}

func SaveLogToDiagnostic(diagnostics []*schema.Diagnostic) {
	_ = PrintDiagnostic(diagnostics)
}

func PrintDiagnostic(diagnostics []*schema.Diagnostic) error {
	var err error
	for i := range diagnostics {
		if int(diagnostics[i].Level()) >= levelMap[global.LOGLEVEL] {
			defaultLogger.Log(hclog.Level(levelMap[global.LOGLEVEL]+1), diagnostics[i].Content())
			PrintCustomizeLn(levelColor[int(diagnostics[i].Level())], diagnostics[i].Content())
			if diagnostics[i].Level() == schema.DiagnosisLevelError {
				err = errors.New(diagnostics[i].Content())
			}
		}
	}
	return err
}
