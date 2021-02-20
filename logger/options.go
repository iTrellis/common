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

type Option func(*Options)

type Options struct {
	level     Level
	buffer    int
	separator string
	color     bool
}

// WriteLevel 设置等级
func WriteLevel(lvl Level) Option {
	return func(c *Options) {
		c.level = lvl
	}
}

// Buffer 设置Chan的大小
func Buffer(buffer int) Option {
	return func(c *Options) {
		c.buffer = buffer
	}
}

// Separator 设置打印分隔符
func Separator(separator string) Option {
	return func(c *Options) {
		c.separator = separator
	}
}

// Color 设置打印分隔符
func Color() Option {
	return func(c *Options) {
		c.color = true
	}
}

type fileWriterOptions struct {
	level Level

	separator  string
	fileName   string
	maxLength  int64
	chanBuffer int

	moveFileType MoveFileType
	// 最大保留日志个数，如果为0则全部保留
	maxBackupFile int
}

// OptionFileWriter 操作配置函数
type OptionFileWriter func(*fileWriterOptions)

// FileWriterLevel 设置等级
func FileWriterLevel(lvl Level) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.level = lvl
	}
}

// FileWriterBuffer 设置Chan的大小
func FileWriterBuffer(buffer int) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.chanBuffer = buffer
	}
}

// FileWriterSeparator 设置打印分隔符
func FileWriterSeparator(separator string) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.separator = separator
	}
}

// FileWriterFileName 设置文件名
func FileWriterFileName(name string) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.fileName = name
	}
}

// FileWriterMaxLength 设置最大文件大小
func FileWriterMaxLength(length int64) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.maxLength = length
	}
}

// FileWriterMaxBackupFile 文件最大数量
func FileWriterMaxBackupFile(num int) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.maxBackupFile = num
	}
}

// FileWriterMoveFileType 设置移动文件的类型
func FileWriterMoveFileType(typ MoveFileType) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.moveFileType = typ
	}
}
