package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjackv2 "gopkg.in/natefinch/lumberjack.v2"
	gormLogger "gorm.io/gorm/logger"
)

const (
	Console = "console"
	File    = "file"
)

var (
	Debug   = zap.DebugLevel
	Target  = File
	DBLevel = gormLogger.Info
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func Init(filename string) {
	w := zapcore.AddSync(&lumberjackv2.Logger{
		Filename:   filename,
		MaxSize:    1024, // megabytes
		MaxBackups: 10,
		MaxAge:     7, // days
	})

	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewEncoderConfig()),
		writeSyncer,
		Debug,
	)

	Logger = zap.New(core, zap.AddCaller())
	Sugar = Logger.Sugar()
}
