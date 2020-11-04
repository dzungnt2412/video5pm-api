package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *HPLog

func init() {
	InitLoggerDefaultDev()
}

// InitLoggerDefault -- init logger default
func InitLoggerDefault() {
	// init production encoder conf
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	// init production conf
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoderCfg
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	// build logger
	logger, _ := cfg.Build()

	sugarLog := logger.Sugar()
	cfgParams := make(map[string]interface{})
	Log = &HPLog{cfgParams, cfg.Level, logger, sugarLog}
}

// InitLoggerDefaultDev -- init logger dev
func InitLoggerDefaultDev() {
	// init production encoder conf
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	// init production conf
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig = encoderCfg
	cfg.OutputPaths = []string{"stdout"}
	// build logger
	logger, _ := cfg.Build()

	sugarLog := logger.Sugar()
	cfgParams := make(map[string]interface{})
	Log = &HPLog{cfgParams, cfg.Level, logger, sugarLog}
}

// HPLog is a utility struct for logging data in an extremely high performance system.
// We can use both Logger and SugarLog for logging. For more information,
// just visit https://godoc.org/go.uber.org/zap
type HPLog struct {
	// configuration
	config   map[string]interface{}
	logLevel zap.AtomicLevel
	// Logger for logging
	Logger *zap.Logger
	// Sugar for logging
	*zap.SugaredLogger
}

// Close will flush log to file
func (l *HPLog) Close() {
	_ = l.Logger.Sync()
}
