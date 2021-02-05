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
	log := logger.NewLogger()

	w, err := logger.ChanWriter(log, logger.ChanWiterLevel(logger.DebugLevel))
	if err != nil {
		panic(err)
	}
	fw, err := logger.FileWriter(log,
		logger.FileWiterFileName("haha.log"),
		logger.FileWiterLevel(logger.DebugLevel),
		logger.FileWiterRoutings(2),
		logger.FileWiterMoveFileType(1),
	)
	if err != nil {
		panic(err)
	}
	// _, _ = log.Subscriber(w)
	// _, _ = log.Subscriber(fw)

	log = log.WithPrefix(logger.Stack)

	for index := 0; index < 10; index++ {
		logger.Debug(log, "example_debug", index, &msg{Name: "haha", Age: 123})

		log.Info("example\tinfo", index, &msg{Name: "i am  info", Age: 123})
		log.Warn("example_warn", index, &msg{Name: "i am warn", Age: 123})
		log.Error("example error", index, &msg{Name: "i am error", Age: 123})
		log.Critical("example_critical", index, &msg{Name: "i am critial", Age: 123})
		time.Sleep(time.Second)
	}

	w.Stop()
	fw.Stop()
	log.ClearSubscribers()
	log.Error("example error", &msg{Name: "non print", Age: 123})
	time.Sleep(time.Second)
}
