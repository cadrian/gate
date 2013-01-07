package server

import (
	"fmt"
)

type Key interface {
	Name() string
	Password() string
	IsDeleted() bool
	Delete()
	Encoded() string
	Merge(other Key)
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
	return self.pass
}

func (self *key) IsDeleted() bool {
	return self.delcount > self.addcount
}

func (self *key) Delete() {
	self.delcount = self.addcount + 1
}

func (self *key) Encoded() string {
	return fmt.Sprintf("%s:%d:%d:%s", self.name, self.addcount, self.delcount, self.pass)
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
