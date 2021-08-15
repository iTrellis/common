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

package logger

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewWithZapLogger(l *zap.Logger) Logger {
	if l == nil {
		return &noop{}
	}
	return &ZapLogger{logger: l}
}

type ZapLogger struct {
	options *LogConfig
	logger  *zap.Logger
}

var _ Logger = (*ZapLogger)(nil)

func NewLogger(opts ...Option) (*ZapLogger, error) {
	zl := &ZapLogger{
		options: &LogConfig{},
	}
	for _, o := range opts {
		o(zl.options)
	}

	if zl.options.EncoderConfig == nil {
		zl.options.EncoderConfig = &zapcore.EncoderConfig{
			TimeKey:          "ts",
			LevelKey:         "level",
			NameKey:          "log",
			CallerKey:        "caller",
			MessageKey:       "msg",
			StacktraceKey:    "stacktrace",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      zapcore.LowercaseLevelEncoder,
			EncodeTime:       zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration:   zapcore.SecondsDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: zl.options.FileOptions.Separator,
		}
	}

	level := zap.NewAtomicLevelAt(zl.options.Level.ToZapLevel())

	var encoder zapcore.Encoder
	switch zl.options.Encoding {
	case "", "console":
		encoder = zapcore.NewConsoleEncoder(*zl.options.EncoderConfig)
	case "json":
		encoder = zapcore.NewJSONEncoder(*zl.options.EncoderConfig)
	default:
		return nil, errors.New("unknown encoding")
	}

	var ws []zapcore.WriteSyncer
	for _, op := range zl.options.FileOptions.StdPrinters {
		switch op {
		case "stderr":
			ws = append(ws, zapcore.AddSync(os.Stderr))
		case "stdout":
			ws = append(ws, zapcore.AddSync(os.Stdout))
		default:
			return nil, errors.New("unknown std printers")
		}
	}

	if zl.options.FileOptions.Filename != "" {
		if err := zl.options.FileOptions.Check(); err != nil {
			return nil, err
		}

		w, err := NewFileLoggerWithOptions(zl.options.FileOptions)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(ws...), level)

	var options []zap.Option
	if zl.options.CallerSkip != 0 {
		options = append(options, zap.AddCallerSkip(zl.options.CallerSkip))
	}

	if zl.options.StackTrace {
		options = append(options, zap.AddStacktrace(level))
	}

	if zl.options.Caller {
		options = append(options, zap.AddCaller())
	}

	zl.logger = zap.New(core, options...)
	return zl, nil
}

func (p *ZapLogger) GetZapLogger() *zap.Logger {
	return p.logger
}

// With (fields ...Field)
func (p *ZapLogger) With(kvs ...interface{}) Logger {
	newZL := &ZapLogger{
		options: p.options,
	}

	lenFields := len(kvs)
	var fields []zap.Field
	for i := 0; i < lenFields; i += 2 {
		k := kvs[i]
		var v interface{} = errors.New("MISSING VALUE")
		if i+1 < lenFields {
			v = kvs[i+1]
		}
		fields = append(fields, zap.Any(toString(k), v))
	}
	newZL.logger = p.logger.With(fields...)

	return newZL
}

// Debug(msg string, fields ...Field)
func (p *ZapLogger) Log(kvs ...interface{}) error {
	p.Info("", kvs...)
	return nil
}

// Debug(msg string, fields ...Field)
func (p *ZapLogger) Debug(msg string, kvs ...interface{}) {
	fields := p.genKVs(kvs...)
	p.logger.Debug(msg, fields...)
}

// Debugf(msg string, fields ...Field)
func (p *ZapLogger) Debugf(msg string, kvs ...interface{}) {
	p.Debug(fmt.Sprintf(msg, kvs...))
}

// Info(msg string, fields ...Field)
func (p *ZapLogger) Info(msg string, kvs ...interface{}) {
	fields := p.genKVs(kvs...)
	p.logger.Info(msg, fields...)
}

// Infof(msg string, fields ...Field)
func (p *ZapLogger) Infof(msg string, kvs ...interface{}) {
	p.Info(fmt.Sprintf(msg, kvs...))
}

// Warn(msg string, fields ...Field)
func (p *ZapLogger) Warn(msg string, kvs ...interface{}) {
	fields := p.genKVs(kvs...)
	p.logger.Warn(msg, fields...)
}

// Warnf(msg string, fields ...Field)
func (p *ZapLogger) Warnf(msg string, kvs ...interface{}) {
	p.Warn(fmt.Sprintf(msg, kvs...))
}

// Error(msg string, fields ...Field)
func (p *ZapLogger) Error(msg string, kvs ...interface{}) {
	fields := p.genKVs(kvs...)
	p.logger.Error(msg, fields...)
}

// Errorf(msg string, fields ...Field)
func (p *ZapLogger) Errorf(msg string, kvs ...interface{}) {
	p.Error(fmt.Sprintf(msg, kvs...))
}

// Panic(msg string, fields ...Field)
func (p *ZapLogger) Panic(msg string, kvs ...interface{}) {
	fields := p.genKVs(kvs...)
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	p.logger.Panic(msg, fields...)
}

// Panicf(msg string, fields ...Field)
func (p *ZapLogger) Panicf(msg string, kvs ...interface{}) {
	p.Panic(fmt.Sprintf(msg, kvs...))
}

// Fatal(msg string, fields ...Field)
func (p *ZapLogger) Fatal(msg string, kvs ...interface{}) {
	fields := p.genKVs(kvs...)
	p.logger.Fatal(msg, fields...)
}

// Fatalf(msg string, fields ...Field)
func (p *ZapLogger) Fatalf(msg string, kvs ...interface{}) {
	p.Fatal(fmt.Sprintf(msg, kvs...))
}

func (p *ZapLogger) genKVs(kvs ...interface{}) []zap.Field {

	lenFields := len(kvs)
	n := 4 + (lenFields+1)/2*2

	logs := make([]zap.Field, 0, n)

	for i := 0; i < lenFields; i += 2 {
		k := kvs[i]
		var v interface{} = errors.New("MISSING VALUE")
		if i+1 < lenFields {
			v = kvs[i+1]
		}
		logs = append(logs, zap.Any(toString(k), v))
	}

	return logs
}
