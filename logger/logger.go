package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

func NewZapLogger(logFileName, setMode, applicationVersion, applicationName string) *zap.SugaredLogger {
	// encode config mode checking and setup
	var encoderCfg zapcore.EncoderConfig
	if strings.Contains(setMode, "stg") ||
		strings.Contains(setMode, "prd") {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// zap logging config setup using previous encoder config
	var config zap.Config
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.Development = false
	config.DisableCaller = false
	config.DisableStacktrace = false
	config.Sampling = nil
	config.EncoderConfig = encoderCfg
	config.Encoding = "json"
	config.OutputPaths = []string{"stderr"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.InitialFields = map[string]interface{}{
		"service_pid":     os.Getpid(),
		"service_name":    applicationName,
		"service_version": applicationVersion,
	}

	// Set up file output using Lumberjack
	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename: logFileName,
	})

	// Create file core
	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileSyncer, config.Level)
	core := zapcore.NewTee(
		zap.Must(config.Build()).Core(),
		fileCore.With([]zap.Field{
			zap.Int("service_pid", os.Getpid()),
			zap.String("service_name", applicationName),
			zap.String("service_version", applicationVersion),
		}),
	)

	// setup sugared log
	sugared := zap.New(core).Sugar()
	defer func(sugar *zap.SugaredLogger) {
		_ = sugar.Sync()
	}(sugared)
	return sugared
}
