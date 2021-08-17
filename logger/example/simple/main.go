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
	"time"

	"github.com/iTrellis/common/logger"
)

type msg struct {
	Name string
	Age  int
}

func main() {

	fileLog, err := logger.NewFileLogger(logger.OptionFilename("haha.log"))
	if err != nil {
		panic(err)
	}

	stdLog, err := logger.NewLogger(

		logger.LogFileOption(
			logger.OptionFilename("./haha_zap.log"),
			logger.OptionStdPrinters([]string{"stdout"}),
			logger.OptionMoveFileType(1),
			logger.OptionMaxLength(500000),
			logger.OptionMaxBackups(10),
		),
		logger.CallerSkip(1),
		logger.Caller(),
		logger.StackTrace(),
		logger.Encoding("json"),
	)
	if err != nil {
		panic(err)
	}

	fileLog.Write([]byte("key=value"))

	for index := 0; index < 10000; index++ {
		stdLog.Debug("example_debug", "index", index, "msg", &msg{Name: "haha", Age: 123})

		stdLog.Info("example_info", "index", index, "msg", &msg{Name: "i am info", Age: 123})
		stdLog.Warn("example_warn", "index", index, "msg", &msg{Name: "i am warn", Age: 123})
		stdLog.Error("example_error", "index", index, "msg", &msg{Name: "i am error", Age: 123})
		stdLog.Panic("example_panic", "index", index, "msg", &msg{Name: "i am panic", Age: 123})
		time.Sleep(time.Millisecond * 100)
	}

	time.Sleep(time.Second)
	stdLog.Fatal("example_fatal", "msg", &msg{Name: "i am fatal", Age: 123})
}
