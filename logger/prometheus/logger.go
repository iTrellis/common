/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

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
			logger.OptionConcurrencyWrite(),
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
