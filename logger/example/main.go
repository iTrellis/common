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

	stdLog := logger.NewStdLogger()

	stdCallsLog := stdLog.WithPrefix("std_callers_log", logger.RuntimeCallers(9))
	stdLog = stdLog.WithPrefix("std_caller_log", logger.RuntimeCaller(8))

	pub := logger.NewPublisher()
	_, err := pub.Subscriber(stdLog)
	if err != nil {
		panic(err)
	}

	_, err = pub.Subscriber(stdCallsLog)
	if err != nil {
		panic(err)
	}

	fLog, err := logger.NewFileLogger(
		logger.FileFileName("haha.log"),
		logger.FileLevel(logger.DebugLevel),
		logger.FileMoveFileType(1),
	)
	if err != nil {
		panic(err)
	}
	fLog = fLog.WithPrefix("writer", "test")
	_, err = pub.Subscriber(fLog)
	if err != nil {
		panic(err)
	}

	pub = pub.WithPrefix("test", "aha", "pub_caller", logger.RuntimeCallers(5))

	stdLog.Info("key", "value")

	for index := 0; index < 10; index++ {
		stdLog.Debug("index", index, "msg", &msg{Name: "haha", Age: 123})
		stdLog.Info("index", index, "msg", &msg{Name: "haha", Age: 234})

		if index == 5 {
			stdLog.SetLevel(logger.WarnLevel)
		} else if index == 7 {
			pub.SetLevel(logger.ErrorLevel)
		}

		pub.Log(logger.InfoLevel, "example_info", index, "msg", &msg{Name: "i am  info", Age: 123})
		pub.Log(logger.WarnLevel, "example_warn", index, "msg", &msg{Name: "i am warn", Age: 123})
		pub.Log(logger.ErrorLevel, "example error", index, "msg", &msg{Name: "i am error", Age: 123})
		pub.Log(logger.CriticalLevel, "example_critical", index, "msg", &msg{Name: "i am critial", Age: 123})
		time.Sleep(time.Millisecond)
	}

	pub.ClearSubscribers()
	pub.Log(logger.ErrorLevel, "msg", &msg{Name: "non print", Age: 123})
	time.Sleep(time.Second)
}
