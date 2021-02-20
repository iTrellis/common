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
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/iTrellis/common/event"
)

type stdWriter struct {
	options  Options
	logger   Logger
	stopChan chan bool
	logChan  chan *Event
	out      io.Writer

	subscriber event.Subscriber
}

// NewStdWriter 标准窗体的输出对象
func NewStdWriter(log Logger, opts ...Option) (Writer, error) {
	c := &stdWriter{
		logger:   log,
		out:      os.Stdout,
		stopChan: make(chan bool, 1),
	}
	c.init(opts...)

	go c.looperLog()

	var err error
	c.subscriber, err = event.NewDefSubscriber(c.Publish)
	if err != nil {
		return nil, err
	}

	_, err = log.Subscriber(c.subscriber)
	if err != nil {
		c.stopChan <- true
		return nil, err
	}

	return c, nil
}

func (p *stdWriter) init(opts ...Option) {

	for _, o := range opts {
		o(&p.options)
	}

	if p.options.buffer <= 0 {
		p.logChan = make(chan *Event, defaultChanBuffer)
	} else {
		p.logChan = make(chan *Event, p.options.buffer)
	}

	if len(p.options.separator) == 0 {
		p.options.separator = "\t"
	}
}

func (p *stdWriter) Publish(evts ...interface{}) {
	for _, evt := range evts {
		switch eType := evt.(type) {
		case Event:
			p.logChan <- &eType
		case *Event:
			p.logChan <- eType
		case Level:
			p.options.level = eType
		default:
			panic(fmt.Errorf("unsupported event type: %s", reflect.TypeOf(evt).Name()))
		}
	}
}

func (p *stdWriter) looperLog() {
	for {
		select {
		case log := <-p.logChan:
			if log.Level < p.Level() {
				continue
			}
			data := []byte(generateLogs(log, p.options.separator))
			if p.options.color {
				if color := LevelColors[log.Level]; len(color) != 0 {
					data = []byte(fmt.Sprintf(color, string(data)))
				}
			}
			_, _ = p.Write(data)
		case <-p.stopChan:
			return
		}
	}
}

func (p stdWriter) Write(bs []byte) (int, error) {
	return p.out.Write(bs)
}

func (p *stdWriter) Level() Level {
	return p.options.level
}

func (p *stdWriter) GetID() string {
	return p.subscriber.GetID()
}

func (p *stdWriter) Stop() {
	if err := p.logger.RemoveSubscriber(p.subscriber.GetID()); err != nil {
		p.logger.Criticalf("failed remove Chan Writer: %s", err.Error())
	}
	p.stopChan <- true
}
