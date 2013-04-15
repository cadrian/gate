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

package server

import (
	"gate/core/errors"
)

import (
	"fmt"
	"regexp"
	"strconv"
)

// A vault key
type Key interface {
	Name() string
	Password() string
	IsDeleted() bool
	Delete()
	Encoded() string
	Merge(other Key)
	SetPassword(pass string)
}

func decode_group(dec *regexp.Regexp, data string, name string, match []int) (result string) {
	result = string(dec.ExpandString(make([]byte, 0, 1024), fmt.Sprintf("${%s}", name), data, match))
	return
}

func decode_group_int(dec *regexp.Regexp, data string, name string, match []int) (result int64, err error) {
	s := decode_group(dec, data, name, match)
	result, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.Decorated(err)
	}
	return
}
