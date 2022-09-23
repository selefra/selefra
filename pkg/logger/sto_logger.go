package logger

import (
	"fmt"
	"github.com/selefra/selefra/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

type StoLogger struct {
	logger *zap.Logger
	config *Config
	name   string
}

func (l *StoLogger) DebugF(msg string, args ...any) {
	l.logger.Debug(fmt.Sprintf(msg, args...))
}

func (l *StoLogger) InfoF(msg string, args ...any) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func (l *StoLogger) WarnF(msg string, args ...any) {
	l.logger.Warn(fmt.Sprintf(msg, args...))
}

func (l *StoLogger) ErrorF(msg string, args ...any) {
	l.logger.Error(fmt.Sprintf(msg, args...))
}

func (l *StoLogger) Fatal(msg string, args ...zap.Field) {
	l.logger.Fatal(fmt.Sprintf(msg, args))
}

func (l *StoLogger) FatalF(msg string, args ...any) {
	l.logger.Fatal(fmt.Sprintf(msg, args))
}

func (l *StoLogger) Debug(msg string, args ...zap.Field) {
	l.logger.Debug(fmt.Sprintf(msg, args))
}

func (l *StoLogger) Info(msg string, args ...zap.Field) {
	l.logger.Info(fmt.Sprintf(msg, args))
}

func (l *StoLogger) Warn(msg string, args ...zap.Field) {
	l.logger.Warn(fmt.Sprintf(msg, args))
}

func (l *StoLogger) Error(msg string, args ...zap.Field) {
	l.logger.Error(fmt.Sprintf(msg, args))
}

func (l *StoLogger) Name() string {
	return l.name
}

func NewStoLogger(c Config) (*StoLogger, error) {
	logDir := filepath.Join(*global.WORKSPACE, c.Directory)
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(logDir, 0755)
	}
	if err != nil {
		return nil, nil
	}
	logger := zap.New(zapcore.NewTee(c.GetEncoderCore()...))

	if c.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	return &StoLogger{logger: logger, config: &c}, nil
}
