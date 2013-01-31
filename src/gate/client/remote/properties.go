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

// Core properties management

import (
	"gate/core/errors"
)

import (
	"fmt"
	"io"
)

type properties struct {
	allowed map[string]bool
	properties map[string]string
}

func (self *properties) setProperty(key, value string) (err error) {
	for k, mandatory := range self.allowed {
		if k == key {
			if value == "" && mandatory {
				err = errors.Newf("cannot reset '%s': mandatory property", key)
			} else {
				self.properties[key] = value
			}
			return
		}
	}
	return errors.Newf("key not allowed: %s", key)
}

func (self *properties) resetProperty(key string) (err error) {
	if self.allowed[key] {
		err = errors.Newf("cannot reset '%s': mandatory property", key)
	} else {
		delete(self.properties, key)
	}
	return
}

func (self *properties) countProperties() int {
	return len(self.properties)
}

func (self *properties) getProperty(name string) (result string) {
	result, _ = self.properties[name]
	return
}

func (self *properties) storeProperties(out io.Writer) (err error) {
	if err != nil {
		return errors.Decorated(err)
	}

	for property, value := range self.properties {
		if value != "" {
			_, err = out.Write([]byte(fmt.Sprintf("%s = %s\n", property, value)))
			if err != nil {
				return errors.Decorated(err)
			}
		}
	}

	return
}
