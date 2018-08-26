package simplelogger

import (
	"runtime/debug"

	"strings"

	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// Level is type of log levels
type Level = zapcore.Level

var (
	// DebugLevel level
	DebugLevel = Level(zap.DebugLevel)
	// InfoLevel level
	InfoLevel = Level(zap.InfoLevel)
	// WarnLevel level
	WarnLevel = Level(zap.WarnLevel)
	// ErrorLevel level
	ErrorLevel = Level(zap.ErrorLevel)
	// PanicLevel level
	PanicLevel = Level(zap.PanicLevel)
	// FatalLevel level
	FatalLevel = Level(zap.FatalLevel)
)

var (
	cfg          zap.Config
	logger       *zap.Logger
	sugar        *zap.SugaredLogger
	currentLevel Level
)

func init() {
	var err error
	cfgJson := []byte(`{
		"level": "debug",
	"outputPaths": ["stderr"],
	"errorOutputPaths": ["stderr"],
	"encoding": "console",
		"encoderConfig": {
		"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
	}
}`)
	currentLevel = DebugLevel

	if err = json.Unmarshal(cfgJson, &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	rebuildLoggerFromCfg()
}

// SetLevel sets the log level
func SetLevel(lv Level) {
	currentLevel = lv
	cfg.Level.SetLevel(lv)
}

// GetLevel get the current log level
func GetLevel() Level {
	return currentLevel
}

// TraceError prints the stack and error
func TraceError(format string, args ...interface{}) {
	Error(string(debug.Stack()))
	Errorf(format, args...)
}

// SetOutput sets the output writer
func SetOutput(outputs []string) {
	cfg.OutputPaths = outputs
	rebuildLoggerFromCfg()
}

// ParseLevel converts string to Levels
func ParseLevel(s string) Level {
	if strings.ToLower(s) == "debug" {
		return DebugLevel
	} else if strings.ToLower(s) == "info" {
		return InfoLevel
	} else if strings.ToLower(s) == "warn" || strings.ToLower(s) == "warning" {
		return WarnLevel
	} else if strings.ToLower(s) == "error" {
		return ErrorLevel
	} else if strings.ToLower(s) == "panic" {
		return PanicLevel
	} else if strings.ToLower(s) == "fatal" {
		return FatalLevel
	}
	Errorf("ParseLevel: unknown level: %s", s)
	return DebugLevel
}

func rebuildLoggerFromCfg() {
	if newLogger, err := cfg.Build(); err == nil {
		if logger != nil {
			logger.Sync()
		}
		logger = newLogger
		setSugar(logger.Sugar())
	} else {
		panic(err)
	}
}

func Debugf(format string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05.000")}, args...)
	sugar.Debugf("%s - "+format, args...)
}

func Infof(format string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05.000")}, args...)
	sugar.Infof("%s - "+format, args...)
}

func Warnf(format string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05.000")}, args...)
	sugar.Warnf("%s - "+format, args...)
}

func Errorf(format string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05.000")}, args...)
	sugar.Errorf("%s - "+format, args...)
}

func Panicf(format string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05.000")}, args...)
	sugar.Panicf("%s - "+format, args...)
}

func Fatalf(format string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05.000")}, args...)
	sugar.Fatalf("%s - "+format, args...)
}

func Error(args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05.000"), " - "}, args...)
	sugar.Error(args...)
}

func Panic(args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05.000"), " - "}, args...)
	sugar.Panic(args...)
}

func Fatal(args ...interface{}) {
	args = append([]interface{}{time.Now().Format("2006-01-02 15:04:05.000"), " - "}, args...)
	sugar.Fatal(args...)
}

func setSugar(sugar_ *zap.SugaredLogger) {
	sugar = sugar_
}
