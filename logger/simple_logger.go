package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapSugarLogger func(msg string, keysAndValues ...interface{})

func (l zapSugarLogger) Log(kv ...interface{}) error {
	l("", kv...)
	return nil
}

// NewZapSugarLogger returns a Go kit log.Logger that sends
// log events to a zap.Logger.
func NewZapSugarLogger(logger *zap.Logger, level zapcore.Level, opts ...zap.Option) SimpleLogger {
	sugarLogger := logger.WithOptions(opts...).Sugar()
	var sugar zapSugarLogger
	switch level {
	case zapcore.DebugLevel:
		sugar = sugarLogger.Debugw
	case zapcore.InfoLevel:
		sugar = sugarLogger.Infow
	case zapcore.WarnLevel:
		sugar = sugarLogger.Warnw
	case zapcore.ErrorLevel:
		sugar = sugarLogger.Errorw
	case zapcore.DPanicLevel:
		sugar = sugarLogger.DPanicw
	case zapcore.PanicLevel:
		sugar = sugarLogger.Panicw
	case zapcore.FatalLevel:
		sugar = sugarLogger.Fatalw
	default:
		sugar = sugarLogger.Infow
	}
	return sugar
}
