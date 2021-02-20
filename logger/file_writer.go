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
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/iTrellis/common/event"
	"github.com/iTrellis/common/files"
)

type fileWriter struct {
	logger Logger

	options fileWriterOptions

	logChan chan *Event

	subscriber event.Subscriber

	stopChan chan bool

	writeFileTime time.Time
	lastMoveFlag  int
	ticker        *time.Ticker
}

// MoveFileType move file type
type MoveFileType int

// MoveFileTypes
const (
	MoveFileTypeNothing   MoveFileType = iota // 不移动
	MoveFileTypePerMinite                     // 按分钟移动
	MoveFileTypeHourly                        // 按小时移动
	MoveFileTypeDaily                         // 按天移动
)

// FileWriter 标准窗体的输出对象
func FileWriter(log Logger, opts ...OptionFileWriter) (Writer, error) {
	fw := &fileWriter{
		logger:   log,
		ticker:   time.NewTicker(time.Second * 30),
		stopChan: make(chan bool, 1),
	}

	err := fw.init(opts...)
	if err != nil {
		return nil, err
	}

	fw.subscriber, err = event.NewDefSubscriber(fw.Publish)
	if err != nil {
		fw.Stop()
		return nil, err
	}

	_, err = log.Subscriber(fw.subscriber)
	if err != nil {
		return nil, err
	}
	return fw, err
}

var fileExecutor = files.New()

func (p *fileWriter) init(opts ...OptionFileWriter) error {

	for _, o := range opts {
		o(&p.options)
	}

	if len(p.options.fileName) == 0 {
		return errors.New("file name not exist")
	}

	if p.options.chanBuffer == 0 {
		p.logChan = make(chan *Event, defaultChanBuffer)
	} else {
		p.logChan = make(chan *Event, p.options.chanBuffer)
	}

	if len(p.options.separator) == 0 {
		p.options.separator = "\t"
	}

	fi, err := fileExecutor.FileInfo(p.options.fileName)
	if err == nil {
		// 说明文件存在
		p.writeFileTime = fi.ModTime()
	} else {
		// 没有文件创建文件
		_, err = fileExecutor.WriteAppend(p.options.fileName, "")
		if err != nil {
			return err
		}
	}

	switch p.options.moveFileType {
	case MoveFileTypePerMinite:
		p.lastMoveFlag = p.writeFileTime.Minute()
	case MoveFileTypeHourly:
		p.lastMoveFlag = p.writeFileTime.Hour()
	case MoveFileTypeDaily:
		p.lastMoveFlag = p.writeFileTime.Day()
	}

	go p.looperLog()

	return nil
}

func (p *fileWriter) Publish(evts ...interface{}) {
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

func (p *fileWriter) Write(bs []byte) (int, error) {
	return fileExecutor.WriteAppendBytes(p.options.fileName, bs)
}

func (p *fileWriter) looperLog() {
	for {
		select {
		case log := <-p.logChan:
			if log.Level >= p.options.level {
				_, _ = p.innerLog(log)
			}
		case t := <-p.ticker.C:
			flag := 0
			switch p.options.moveFileType {
			case MoveFileTypePerMinite:
				flag = t.Minute()
			case MoveFileTypeHourly:
				flag = t.Hour()
			case MoveFileTypeDaily:
				flag = t.Day()
			}
			if p.lastMoveFlag == flag {
				continue
			}
			_ = p.judgeMoveFile()
		case <-p.stopChan:
			p.ticker.Stop()
			return
		}
	}
}

func (p *fileWriter) judgeMoveFile() error {

	timeNow, flag := time.Now(), 0
	switch p.options.moveFileType {
	case MoveFileTypePerMinite:
		flag = timeNow.Minute()
	case MoveFileTypeHourly:
		flag = timeNow.Hour()
	case MoveFileTypeDaily:
		flag = time.Now().Day()
	default:
		return nil
	}

	if flag == p.lastMoveFlag {
		return nil
	}
	p.lastMoveFlag = flag
	p.writeFileTime = time.Now()
	return p.moveFile()
}

func (p *fileWriter) moveFile() error {
	var timeStr string
	switch p.options.moveFileType {
	case MoveFileTypePerMinite:
		timeStr = time.Now().Format("200601021504-05.999999999")
	case MoveFileTypeHourly:
		timeStr = time.Now().Format("2006010215-0405.999999999")
	case MoveFileTypeDaily:
		timeStr = time.Now().Format("20060102-150405.999999999")
	}

	err := fileExecutor.Rename(p.options.fileName, fmt.Sprintf("%s_%s", p.options.fileName, timeStr))
	if err != nil {
		return err
	}

	if err = p.removeOldFiles(); err != nil {
		return err
	}

	_, err = fileExecutor.Write(p.options.fileName, "")

	return err
}

func (p *fileWriter) removeOldFiles() error {
	if 0 == p.options.maxBackupFile {
		return nil
	}

	path := filepath.Dir(p.options.fileName)

	// 获取日志文件列表
	dirLis, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// 根据文件名过滤日志文件
	fileSort := FileSort{}
	fileNameSplit := strings.Split(p.options.fileName, "/")
	filePrefix := fmt.Sprintf("%s_", fileNameSplit[len(fileNameSplit)-1])
	for _, f := range dirLis {
		if strings.Contains(f.Name(), filePrefix) {
			fileSort = append(fileSort, f)
		}
	}

	if len(fileSort) <= int(p.options.maxBackupFile) {
		return nil
	}

	// 根据文件修改日期排序，保留最近的N个文件
	sort.Sort(fileSort)
	for _, f := range fileSort[p.options.maxBackupFile:] {
		err := os.Remove(path + "/" + f.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *fileWriter) innerLog(evt *Event) (n int, err error) {

	if err = p.judgeMoveFile(); err != nil {
		return
	}

	n, err = p.Write([]byte(generateLogs(evt, p.options.separator)))

	if p.options.maxLength == 0 {
		return
	}

	fi, e := fileExecutor.FileInfo(p.options.fileName)
	if e != nil {
		return 0, e
	}

	if p.options.maxLength > fi.Size() {
		return
	}

	err = p.moveFile()

	return
}

func (p *fileWriter) Level() Level {
	return p.options.level
}

func (p *fileWriter) GetID() string {
	return p.subscriber.GetID()
}

func (p *fileWriter) Stop() {
	if err := p.logger.RemoveSubscriber(p.subscriber.GetID()); err != nil {
		p.logger.Criticalf("failed remove Chan Writer: %s", err.Error())
	}

	p.stopChan <- true

	close(p.logChan)
}

// FileSort 文件排序
type FileSort []os.FileInfo

func (fs FileSort) Len() int {
	return len(fs)
}

func (fs FileSort) Less(i, j int) bool {
	return fs[i].Name() > fs[j].Name()
}

func (fs FileSort) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}
