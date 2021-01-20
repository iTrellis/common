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
	"fmt"

	"github.com/iTrellis/common/files"
)

const (
	oldpath = "test.txt"
	newpath = "testing.txt"
)

func main() {

	f := files.New()

	n, e := f.Write(oldpath, "testing 1\n")
	if e != nil {
		fmt.Println("failed rewrite:", e)
		return
	}
	fmt.Println("write bytes:", n)

	fi, e := f.FileInfo(oldpath)
	if e != nil {
		fmt.Println("failed get file information:", e)
		return
	}

	fmt.Println(fi)

	b, n, e := f.Read(oldpath)
	if e != nil {
		fmt.Println("failed read:", e)
		return
	}
	fmt.Println("read bytes:", n)
	fmt.Println("read content:")
	fmt.Println(string(b))

	if e = f.Rename(oldpath, newpath); e != nil {
		fmt.Println("failed rename:", e)
		return
	}

	n, e = f.WriteAppendBytes(newpath, []byte("testing 2\n"))
	if e != nil {
		fmt.Println("failed write append:", e)
		return
	}
	fmt.Println("write append bytes:", n)

	if e = f.Rename(newpath, oldpath); e != nil {
		fmt.Println("failed rename:", e)
		return
	}

	b, n, e = f.Read(oldpath)
	if e != nil {
		fmt.Println("failed read again:", e)
		return
	}
	fmt.Println("read bytes:", n)
	fmt.Println("content:")
	fmt.Println(string(b))
}
