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
	pub := logger.NewPublisher()

	stdLog := logger.NewStdLogger()

	stdLog = logger.WithPrefix(stdLog, "test", "aha")
	_, err := pub.Subscriber(stdLog)
	if err != nil {
		panic(err)
	}

	_, err = logger.NewFileLogger(
		logger.FileFileName("haha.log"),
		logger.FileLevel(logger.DebugLevel),
		logger.FileMoveFileType(1),
		logger.FilePublisher(pub),
	)
	if err != nil {
		panic(err)
	}

	pub = pub.WithPrefix("stack", logger.Stack)

	for index := 0; index < 10; index++ {
		logger.Debug(pub, "index", index, "msg", &msg{Name: "haha", Age: 123})

		if index == 5 {
			stdLog.SetLevel(logger.InfoLevel)
		}

		pub.Info("example_info", index, "msg", &msg{Name: "i am  info", Age: 123})
		pub.Warn("example_warn", index, "msg", &msg{Name: "i am warn", Age: 123})
		pub.Error("example error", index, "msg", &msg{Name: "i am error", Age: 123})
		pub.Critical("example_critical", index, "msg", &msg{Name: "i am critial", Age: 123})
		time.Sleep(time.Second)
	}

	pub.ClearSubscribers()
	pub.Error("msg", &msg{Name: "non print", Age: 123})
	time.Sleep(time.Second)
}
