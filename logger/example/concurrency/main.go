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

package main

import (
	"sync"

	"github.com/go-kit/log/level"
	"github.com/iTrellis/common/logger"
	"go.uber.org/zap/zapcore"
)

func main() {

	stdLog, err := logger.NewLogger(

		logger.LogFileOption(
			logger.OptionFilename("./haha_zap.log"),
			logger.OptionStdPrinters([]string{"stdout"}),
			logger.OptionMoveFileType(1),
			logger.OptionMaxLength(5000000),
			// logger.OptionMaxBackups(10),
		),
		logger.CallerSkip(1),
		logger.Caller(),
		logger.EncoderConfig(&zapcore.EncoderConfig{
			TimeKey:    "ts",
			EncodeTime: zapcore.RFC3339NanoTimeEncoder,
		}),
		logger.Encoding("console"),
	)
	if err != nil {
		panic(err)
	}

	std2 := level.NewFilter(stdLog, level.AllowAll())

	wg := sync.WaitGroup{}
	f := func(index int) {
		defer wg.Done()
		for sub := 0; sub < 10000; sub++ {
			level.Info(std2).Log("msg", "test1111", "index", index, "sub", sub)
			level.Info(std2).Log("msg", "test2222", "index", index, "sub", sub)
			level.Info(std2).Log("msg", "test3333", "index", index, "sub", sub)
		}
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go f(i)
	}
	wg.Wait()
}
