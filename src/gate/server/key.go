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

// Vault keys

import (
	"fmt"
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

var _ Key = &key{}

type key struct {
	name string
	pass string
	delcount int64
	addcount int64
}

func (self *key) Name() string {
	return self.name
}

func (self *key) Password() string {
	if self.IsDeleted() {
		return ""
	}
	return self.pass
}

func (self *key) IsDeleted() bool {
	return self.delcount > self.addcount
}

func (self *key) Delete() {
	self.delcount = self.addcount + 1
}

func (self *key) Encoded() string {
	return fmt.Sprintf("%s:%d:%d:%s\n", self.name, self.addcount, self.delcount, self.pass)
}

func (self *key) Merge(other Key) {
	okey := other.(*key)

	if self.delcount < okey.delcount {
		self.delcount = okey.delcount
	}
	if self.addcount < okey.addcount {
		self.pass = okey.pass
		self.addcount = okey.addcount
	}
}

func (self *key) SetPassword(pass string) {
	self.pass = pass
	self.addcount = self.addcount + 1
}
