// This file is part of Gate.
// Copyright (C) 2012-2015 Cyril Adrian <cyril.adrian@gmail.com>
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

// Blowfish vault keys (now considered weak)

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var _ Key = &bf_key{}

type bf_key struct {
	name     string
	pass     string
	delcount int64
	addcount int64
}

func (self *bf_key) Name() string {
	return self.name
}

func (self *bf_key) Password() string {
	if self.IsDeleted() {
		return ""
	}
	return self.pass
}

func (self *bf_key) IsDeleted() bool {
	return self.delcount > self.addcount
}

func (self *bf_key) Delete() {
	self.delcount = self.addcount + 1
}

func (self *bf_key) Encoded() string {
	return fmt.Sprintf("%s:%d:%d:%s\n", self.name, self.addcount, self.delcount, self.pass)
}

func (self *bf_key) Merge(other Key) {
	okey := other.(*bf_key)

	if self.delcount < okey.delcount {
		self.delcount = okey.delcount
	}
	if self.addcount < okey.addcount {
		self.pass = okey.pass
		self.addcount = okey.addcount
	}
}

func (self *bf_key) SetPassword(pass string) {
	self.pass = pass
	self.addcount = self.addcount + 1
}

var bf_decoder = regexp.MustCompile("(?P<name>[^:]+):(?P<add>[0-9]+):(?P<del>[0-9]+):(?P<pass>.*)")

func bf_decode(v *vault, out io.ReadCloser, barrier chan error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))
	_, err := buffer.ReadFrom(out)
	if err != nil {
		barrier <- err
		return
	}
	data := string(buffer.Bytes())

	for _, line := range strings.Split(data, "\n") {
		if line != "" {
			linematch := bf_decoder.FindSubmatchIndex([]byte(line))
			name := decode_group(bf_decoder, line, "name", linematch)
			pass := decode_group(bf_decoder, line, "pass", linematch)
			delcount, err := decode_group_int(bf_decoder, line, "del", linematch)
			if err != nil {
				barrier <- err
				continue
			}
			addcount, err := decode_group_int(bf_decoder, line, "add", linematch)
			if err != nil {
				barrier <- err
				continue
			}

			k := &bf_key{
				name:     name,
				pass:     pass,
				delcount: delcount,
				addcount: addcount,
			}
			v.data[name] = k
		}
	}

	barrier <- io.EOF
}

func bf_newkey(name string, pass string) Key {
	return &bf_key{
		name:     name,
		pass:     pass,
		delcount: 0,
		addcount: 1,
	}
}
