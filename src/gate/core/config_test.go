// This file is part of Gate.
// Copyright (C) 2012-2013 Cyril Adrian <cyril.adrian@gmail.com>
//
// Gate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// Gate is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.	 See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gate.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"testing"
)

func mock_getenv(name string) string {
	switch name {
	case "FOO":
		return "bar"
	}
	return ""
}

func checkeval(expected string, toeval string, t *testing.T) {
	actual := eval(toeval, mock_getenv)
	if expected != actual {
		t.Errorf("[%s] => %s != %s", toeval, expected, actual)
	}
}

func TestExpand(t *testing.T) {
	checkeval("bar", "$FOO", t)
	checkeval("bar", "${FOO}", t)
	checkeval("$FOO", "\\$FOO", t)
	checkeval("\"bar\"", "\"$FOO\"", t)
	checkeval("\"'bar'\"", "\"'$FOO'\"", t)
	checkeval("'$FOO'", "'$FOO'", t)
	checkeval("'\"$FOO\"'", "'\"$FOO\"'", t)
	checkeval("bar'$FOO'", "$FOO'$FOO'", t)
	checkeval("bar # anything $FOO ", "$FOO # anything $FOO ", t)
}
