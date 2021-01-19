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

package errors

import (
	"fmt"
	"testing"

	"github.com/iTrellis/go-common/testutils"
)

func TestNew(t *testing.T) {
	sErr := New("new_test")

	testutils.Equals(t, sErr.FullError(), fmt.Sprintf("%s#%s:%s", sErr.Namespace(), sErr.ID(), sErr.Message()))
	testutils.Equals(t, sErr.Error(), sErr.Message())
	testutils.Equals(t, sErr.Namespace(), "T:S")
	testutils.Equals(t, sErr.Message(), "new_test")
	testutils.Assert(t, sErr.ID() != "", "error id should not be nil")

	fErr := Newf("format_err_test: %d, %s", 1, "test")

	testutils.Equals(t, fErr.Error(), "format_err_test: 1, test")
	testutils.Equals(t, fErr.FullError(), fmt.Sprintf("%s#%s:format_err_test: 1, test", fErr.Namespace(), fErr.ID()))
	testutils.Equals(t, fErr.Namespace(), "T:S")
	testutils.Equals(t, fErr.Message(), "format_err_test: 1, test")
	testutils.Assert(t, fErr.ID() != "", "error id should not be nil")
}
