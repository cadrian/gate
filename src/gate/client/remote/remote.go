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

package remote

import (
	"gate/core"
	"gate/core/errors"
	"gate/core/exec"
)

import (
	"io"
)

type Remoter interface {
	Remote(name string) (Remote, error)
}

type remoter struct {
	config core.Config
	remotes map[string]Remote
}

var _ Remoter = &remoter{}

type Remote interface {
	Name() string

	LoadVault(file string) error
	SaveVault(file string) error

	Proxy() Proxy

	SetProperty(key, value string) error
	ResetProperty(key string) error
	StoreProperties(io.Writer) error
}

type remote struct {
	properties
	remoter Remoter
	name string
	proxy Proxy
}

type Proxy interface {
	Install(cmd *exec.Cmd) error

	SetProperty(key, value string) error
	ResetProperty(key string) error
	StoreProperties(io.Writer) error
}

func NewRemoter(config core.Config) Remoter {
	return &remoter {
		config: config,
		remotes: make(map[string]Remote, 32),
	}
}

func (self *remoter) Remote(name string) (result Remote, err error) {
	result, ok := self.remotes[name]
	if !ok {
		// TODO try and load it
		err = errors.Newf("Unknown remote: %s", name)
	}
	return
}
