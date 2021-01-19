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

	"github.com/iTrellis/go-common/formats"
	"github.com/iTrellis/go-common/testutils"
)

func Test_QuoteToASCIIWithoutBackslashed(t *testing.T) {
	type Test struct {
		A int    `json:""`
		Q string `json:""`
	}
	test := Test{A: 1, Q: "中国"}
	bs, _ := json.Marshal(test)

	revertToString := formats.QuoteBytesToASCIIWithoutBackslashed(bs)
	nTest := &Test{}
	err := json.Unmarshal([]byte(revertToString), nTest)
	testutils.Ok(t, err)
	testutils.Equals(t, 1, nTest.A)
	testutils.Equals(t, "中国", nTest.Q)

	testutils.Equals(t, "\\a\\b\\f\\n\\r\\t\\v' @", formats.QuoteToASCIIWithoutBackslashed("\a\b\f\n\r\t\v' @"))
}
