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

package impl

// scrypt vault keys (strong keys)

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var _ Key = &scrypt_key{}

type scrypt_key struct {
	name	 string
	salt	 string
	pass	 string
	delcount int64
	addcount int64
}

func (self *scrypt_key) Name() string {
	return self.name
}

func (self *scrypt_key) Password() string {
	if self.IsDeleted() {
		return ""
	}
	return self.pass
}

func (self *scrypt_key) IsDeleted() bool {
	return self.delcount > self.addcount
}

func (self *scrypt_key) Delete() {
	self.delcount = self.addcount + 1
}

func (self *scrypt_key) Encoded() string {
	return fmt.Sprintf("%s:%s:%d:%d:%s\n", self.name, self.salt, self.addcount, self.delcount, self.pass)
}

func (self *scrypt_key) Merge(other Key) {
	okey := other.(*scrypt_key)

	if self.delcount < okey.delcount {
		self.delcount = okey.delcount
	}
	if self.addcount < okey.addcount {
		self.pass = okey.pass
		self.addcount = okey.addcount
	}
}

func (self *scrypt_key) SetPassword(pass string) {
	self.pass = pass
	self.addcount = self.addcount + 1
}

var scrypt_decoder = regexp.MustCompile("(?P<name>[^:]+):(?P<salt>[^:]+):(?P<add>[0-9]+):(?P<del>[0-9]+):(?P<pass>.*)")

func scrypt_decode(v *vault, out io.ReadCloser, barrier chan error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))
	_, err := buffer.ReadFrom(out)
	if err != nil {
		barrier <- err
		return
	}
	data := string(buffer.Bytes())

	for _, line := range strings.Split(data, "\n") {
		if line != "" {
			linematch := scrypt_decoder.FindSubmatchIndex([]byte(line))
			name := decode_group(scrypt_decoder, line, "name", linematch)
			salt64 := decode_group(scrypt_decoder, line, "salt", linematch)
			salt, err := base64.StdEncoding.DecodeString(salt64)
			if err != nil {
				barrier <- err
				continue
			}
			pass := decode_group(scrypt_decoder, line, "pass", linematch)
			delcount, err := decode_group_int(scrypt_decoder, line, "del", linematch)
			if err != nil {
				barrier <- err
				continue
			}
			addcount, err := decode_group_int(scrypt_decoder, line, "add", linematch)
			if err != nil {
				barrier <- err
				continue
			}

			k := &scrypt_key{
				name:	  name,
				salt:	  string(salt),
				pass:	  pass,
				delcount: delcount,
				addcount: addcount,
			}
			v.data[name] = k
		}
	}

	barrier <- io.EOF
}

func scrypt_newkey(name string, pass string) Key {
	k := &scrypt_key{
		name:	  name,
		pass:	  pass,
		delcount: 0,
		addcount: 1,
	}
	k.set_salt()
	return k
}

func (self *scrypt_key) set_salt() {
	in, err := os.Open("/dev/random")
	if err != nil {
		return
	}
	defer in.Close()

	data := make([]byte, 0, 2)
	length := 0
	for length == 0 {
		n, err := in.Read(data)
		if err != nil {
			return
		}
		if n == 2 {
			length = (int(data[0])&0xff)<<8 | int(data[1])
		}
	}
	salt := make([]byte, 0, length)
	n, err := in.Read(data)
	if err != nil {
		return
	}
	self.salt = base64.StdEncoding.EncodeToString(salt[0:n])
}
