package prometheus

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/iTrellis/common/logger"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level        string
	FileName     string
	MoveFileType int
	MaxLength    int64
	MaxBackups   int
}

func New(config *Config) log.Logger {
	stdLog, err := logger.NewLogger(
		logger.LogFileOption(
			logger.OptionFilename(config.FileName),
			logger.OptionMoveFileType(logger.MoveFileType(config.MoveFileType)),
			logger.OptionMaxLength(config.MaxLength),
			logger.OptionMaxBackups(config.MaxBackups),
		),
		logger.EncoderConfig(&zapcore.EncoderConfig{}),
	)
	if err != nil {
		panic(err)
	}

	timestampFormat := log.TimestampFormat(
		func() time.Time { return time.Now() },
		"2006-01-02T15:04:05.000Z07:00",
	)
	logger := level.NewFilter(stdLog, getLevel(config.Level))
	logger = log.With(logger, "ts", timestampFormat, "caller", log.DefaultCaller)

	return logger
}

// Set get the value of the allowed level.
func getLevel(s string) level.Option {
	switch s {
	case "debug":
		return level.AllowDebug()
	case "info":
		return level.AllowInfo()
	case "warn":
		return level.AllowWarn()
	case "error":
		return level.AllowError()
	default:
		return level.AllowInfo()
	}
}
