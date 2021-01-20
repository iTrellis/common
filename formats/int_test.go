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
	"encoding/json"
	"testing"

	"github.com/iTrellis/common/formats"
	"github.com/iTrellis/common/testutils"
)

func Test_ToInt64(t *testing.T) {
	i32 := int32(1)
	var expt int64 = 1
	f, err := formats.ToInt64(i32)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)
	i64 := int64(2)
	expt = 2
	f, err = formats.ToInt64(i64)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)

	var i int = 10
	expt = 10
	f, err = formats.ToInt64(i)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)

	var s string = "10"
	f, err = formats.ToInt64(s)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)

	var js json.Number = "20"
	expt = 20
	f, err = formats.ToInt64(js)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)

	failed := map[string]string{}
	_, err = formats.ToInt64(failed)
	testutils.NotOk(t, err)
}

func Test_ToInt(t *testing.T) {
	i32 := int32(1)
	var expt int = 1
	f, err := formats.ToInt(i32)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)
	i64 := int64(2)
	expt = 2
	f, err = formats.ToInt(i64)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)

	var i int = 10
	expt = 10
	f, err = formats.ToInt(i)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)

	var s string = "10"
	f, err = formats.ToInt(s)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)

	var js json.Number = "20"
	expt = 20
	f, err = formats.ToInt(js)
	testutils.Equals(t, expt, f)
	testutils.Ok(t, err)

	failed := map[string]string{}
	_, err = formats.ToInt(failed)
	testutils.NotOk(t, err)
}
