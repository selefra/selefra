package logger

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/natefinch/lumberjack"
	"github.com/selefra/selefra/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Logger struct {
	logger *zap.Logger
	config *Config
	name   string
}

func (l *Logger) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.NoLevel:
		return
	case hclog.Trace:
		l.Trace(msg, args...)
	case hclog.Debug:
		l.Debug(msg, args...)
	case hclog.Info:
		l.Info(msg, args...)
	case hclog.Warn:
		l.Warn(msg, args...)
	case hclog.Error:
		l.Warn(msg, args...)
	}
}

func (l *Logger) Trace(msg string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(msg, args...))
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(msg, args...))
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, args...))
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, args...))
}

func (l *Logger) IsTrace() bool {
	return false
}

func (l *Logger) IsDebug() bool {
	return l.config.TranslationLevel() <= zapcore.DebugLevel
}

func (l *Logger) IsInfo() bool {
	return l.config.TranslationLevel() <= zapcore.InfoLevel
}

func (l *Logger) IsWarn() bool {
	return l.config.TranslationLevel() <= zapcore.WarnLevel
}

func (l *Logger) IsError() bool {
	return l.config.TranslationLevel() <= zapcore.ErrorLevel
}

func (l *Logger) ImpliedArgs() []interface{} {
	return nil
}

func (l *Logger) With(args ...interface{}) hclog.Logger {
	return l
}

func (l *Logger) Name() string {
	return l.name
}

func (l *Logger) Named(name string) hclog.Logger {
	l.name = name
	return l
}

func (l *Logger) ResetNamed(name string) hclog.Logger {
	return l
}

func (l *Logger) SetLevel(level hclog.Level) {
	return
}

func (l *Logger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(l.StandardWriter(opts), "", 0)
}

func (l *Logger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return os.Stdin
}

type Config struct {
	Source              string `yaml:"source,omitempty" json:"source,omitempty"`
	FileLogEnabled      bool   `yaml:"file_log_enabled,omitempty" json:"file_log_enabled,omitempty"`
	ConsoleLogEnabled   bool   `yaml:"enable_console_log,omitempty" json:"enable_console_log,omitempty"`
	EncodeLogsAsJson    bool   `yaml:"encode_logs_as_json,omitempty" json:"encode_logs_as_json,omitempty"`
	Directory           string `yaml:"directory,omitempty" json:"directory,omitempty"`
	Level               string `yaml:"level,omitempty" json:"level,omitempty"`
	LevelIdentUppercase bool   `yaml:"level_ident_uppercase,omitempty" json:"level_ident_uppercase,omitempty"`
	MaxAge              int    `yaml:"max_age,omitempty" json:"max_age,omitempty"`
	ShowLine            bool   `yaml:"show_line,omitempty" json:"show_line,omitempty"`
	ConsoleNoColor      bool   `yaml:"console_no_color,omitempty" json:"console_no_color,omitempty"`
	MaxSize             int    `yaml:"max_size,omitempty" json:"max_size,omitempty"`
	MaxBackups          int    `yaml:"max_backups,omitempty" json:"max_backups,omitempty"`
	TimeFormat          string `yaml:"time_format,omitempty" json:"time_format,omitempty"`
	Prefix              string `yaml:"prefix,omitempty" json:"prefix"`
}

func (c *Config) EncodeLevel() zapcore.LevelEncoder {
	switch {
	case c.LevelIdentUppercase && c.ConsoleNoColor:
		return zapcore.CapitalLevelEncoder
	case c.LevelIdentUppercase && !c.ConsoleNoColor:
		return zapcore.CapitalColorLevelEncoder
	case !c.LevelIdentUppercase && c.ConsoleNoColor:
		return zapcore.LowercaseLevelEncoder
	case !c.LevelIdentUppercase && !c.ConsoleNoColor:
		return zapcore.LowercaseColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func (c *Config) TranslationLevel() zapcore.Level {
	switch strings.ToLower(c.Level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (c *Config) GetEncoder() zapcore.Encoder {
	if c.EncodeLogsAsJson {
		return zapcore.NewJSONEncoder(c.GetEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(c.GetEncoderConfig())
}

func (c *Config) GetEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    c.EncodeLevel(),
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func (c *Config) GetLogWriter(level string) zapcore.WriteSyncer {
	filename := filepath.Join(*global.WORKSPACE, c.Directory, c.Source+".log")
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		LocalTime:  true,
		Compress:   false,
	}
	return zapcore.AddSync(lumberjackLogger)
}

func (c *Config) GetEncoderCore() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := c.TranslationLevel(); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, zapcore.NewCore(c.GetEncoder(), c.GetLogWriter(c.TranslationLevel().String()), c.GetLevelPriority(level)))
	}
	return cores
}

func (c *Config) GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool {
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool {
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool {
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool {
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool {
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool {
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool {
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool {
			return level == zap.DebugLevel
		}
	}
}

func NewLogger(c Config) (*Logger, error) {
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

	return &Logger{logger: logger, config: &c}, nil
}
