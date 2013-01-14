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

type Remoter interface {
	Remote(name string) (Remote, error)
}

type remoter struct {
	remotes map[string]Remote
}

type Remote interface {
	Name() string

	Load(file string) error
	Save(file string) error

	Proxy() (Proxy, error)

	SetProperty(key, value string) error
	ResetProperty(key string) error
}

type remote struct {
	remoter remoter
	name string
	properties map[string]string
	proxy Proxy
}

type Proxy interface {
	IsInstalled() bool
	Install()
	Remove()

	SetProperty(key, value string) error
	ResetProperty(key string) error

	url() string
}

type proxy struct {
	properties map[string]string
}
