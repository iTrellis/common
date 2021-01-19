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

package formats_test

import (
	"testing"
	"time"

	"github.com/iTrellis/go-common/formats"
	"github.com/iTrellis/go-common/testutils"
)

func Test_ParseStringByteSize(t *testing.T) {

	s := "10b"
	var expI int64 = 10
	bInt := formats.ParseStringByteSize(s)
	testutils.Equals(t, expI, bInt.Int64())
	s = "10kb"
	expI *= 1000
	bInt = formats.ParseStringByteSize(s)
	testutils.Equals(t, expI, bInt.Int64())
	s = "10mb"
	expI *= 1000
	bInt = formats.ParseStringByteSize(s)
	testutils.Equals(t, expI, bInt.Int64())
	s = "10gb"
	expI *= 1000
	bInt = formats.ParseStringByteSize(s)
	testutils.Equals(t, expI, bInt.Int64())
	s = "10tb"
	expI *= 1000
	bInt = formats.ParseStringByteSize(s)
	testutils.Equals(t, expI, bInt.Int64())
	s = "10pb"
	expI *= 1000
	bInt = formats.ParseStringByteSize(s)
	testutils.Equals(t, expI, bInt.Int64())
	s = "10eb"
	expI *= 1000
	bInt = formats.ParseStringByteSize(s)
	testutils.Equals(t, expI, bInt.Int64())
	s = "10zb"
	expI *= 1000
	bInt = formats.ParseStringByteSize(s)
	testutils.Equals(t, expI, bInt.Int64())
	s = "10yb"
	expI *= 1000
	bInt = formats.ParseStringByteSize(s)
	testutils.Equals(t, expI, bInt.Int64())

	si := "10k"
	var expIi int64 = 10240
	biInt := formats.ParseStringByteSize(si)
	testutils.Equals(t, expIi, biInt.Int64())
	si = "10m"
	expIi *= 1024
	biInt = formats.ParseStringByteSize(si)
	testutils.Equals(t, expIi, biInt.Int64())
	si = "10g"
	expIi *= 1024
	biInt = formats.ParseStringByteSize(si)
	testutils.Equals(t, expIi, biInt.Int64())
	si = "10t"
	expIi *= 1024
	biInt = formats.ParseStringByteSize(si)
	testutils.Equals(t, expIi, biInt.Int64())
	si = "10p"
	expIi *= 1024
	biInt = formats.ParseStringByteSize(si)
	testutils.Equals(t, expIi, biInt.Int64())
	si = "10e"
	expIi *= 1024
	biInt = formats.ParseStringByteSize(si)
	testutils.Equals(t, expIi, biInt.Int64())
	si = "10z"
	expIi *= 1024
	biInt = formats.ParseStringByteSize(si)
	testutils.Equals(t, expIi, biInt.Int64())
	si = "10y"
	expIi *= 1024
	biInt = formats.ParseStringByteSize(si)
	testutils.Equals(t, expIi, biInt.Int64())
}

func Test_ParseStringTime(t *testing.T) {

	s := "1y"
	dt := formats.ParseStringTime(s)
	var expt time.Duration
	testutils.Equals(t, expt, dt)

	dt = formats.ParseStringTime(s, 10)
	expt = 10
	testutils.Equals(t, expt, dt)

	s = "1s"
	dt = formats.ParseStringTime(s)
	expt = time.Second
	testutils.Equals(t, expt, dt)
	s = "1m"
	dt = formats.ParseStringTime(s)
	expt = time.Minute
	testutils.Equals(t, expt, dt)
	s = "1h"
	dt = formats.ParseStringTime(s)
	expt = time.Hour
	testutils.Equals(t, expt, dt)
	s = "1d"
	dt = formats.ParseStringTime(s)
	expt = time.Hour * 24
	testutils.Equals(t, expt, dt)
	s = "1ns"
	dt = formats.ParseStringTime(s)
	expt = time.Nanosecond
	testutils.Equals(t, expt, dt)
	s = "1us"
	dt = formats.ParseStringTime(s)
	expt = time.Microsecond
	testutils.Equals(t, expt, dt)
	s = "1ms"
	dt = formats.ParseStringTime(s)
	expt = time.Millisecond
	testutils.Equals(t, expt, dt)
}
