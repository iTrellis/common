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

// Debug 调试
func Debug(l LogFarm, fields ...interface{}) {
	l.Debug(fields...)
}

// Debugf 调试
func Debugf(l LogFarm, msg string, fields ...interface{}) {
	l.Debugf(msg, fields...)
}

// Info 信息
func Info(l LogFarm, fields ...interface{}) {
	l.Info(fields...)
}

// Infof 信息
func Infof(l LogFarm, msg string, fields ...interface{}) {
	l.Infof(msg, fields...)
}

// Error 错误
func Error(l LogFarm, fields ...interface{}) {
	l.Error(fields...)
}

// Errorf 错误
func Errorf(l LogFarm, msg string, fields ...interface{}) {
	l.Errorf(msg, fields...)
}

// Warn 警告
func Warn(l LogFarm, fields ...interface{}) {
	l.Warn(fields...)
}

// Warnf 警告
func Warnf(l LogFarm, msg string, fields ...interface{}) {
	l.Warnf(msg, fields...)
}

// Critical 异常
func Critical(l LogFarm, fields ...interface{}) {
	l.Critical(fields...)
}

// Criticalf 异常
func Criticalf(l LogFarm, msg string, fields ...interface{}) {
	l.Criticalf(msg, fields...)
}

// Panic 异常
func Panic(l LogFarm, fields ...interface{}) {
	l.Panic(fields...)
}

// Panicf 异常
func Panicf(l LogFarm, msg string, fields ...interface{}) {
	l.Panicf(msg, fields...)
}
