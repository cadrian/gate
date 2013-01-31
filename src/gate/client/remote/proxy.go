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

// Remote proxy

import (
	"gate/core/errors"
	"gate/core/exec"
)

import (
	"fmt"
	"io"
	"os"
)

type proxy struct {
	properties
}

var _ Proxy = &proxy{}

var ProxyAllowedKeys map[string]bool = map[string]bool{
	"host": true,
	"port": false,
	"protocol": false,
	"user": false,
	"pass": false,
}

func newProxy() Proxy {
	return &proxy {
		properties {
			allowed: ProxyAllowedKeys,
			properties: make(map[string]string),
		},
	}
}

func (self *proxy) Install(cmd *exec.Cmd) (err error) {
	var (
		url string
	)
	host := self.getProperty("host")
	user := self.getProperty("user")
	if user == "" {
		url = host
	} else {
		pass := self.getProperty("pass")
		if pass == "" {
			url = fmt.Sprintf("%s@%s", user, host)
		} else {
			url = fmt.Sprintf("%s:%s@%s", user, escape_pass_url(pass), host)
		}
	}
	port := self.getProperty("port")
	if port != "" {
		url = fmt.Sprintf("%s:%s", url, port)
	}
	protocol := self.getProperty("protocol")
	if protocol != "" {
		url = fmt.Sprintf("%s://%s", protocol, url)
	}
	cmd.Env = append(os.Environ(), fmt.Sprintf("ALL_PROXY=%s", url))
	return
}

func (self *proxy) SetProperty(key, value string) error {
	return self.setProperty(key, value)
}

func (self *proxy) ResetProperty(key string) error {
	return self.resetProperty(key)
}

func (self *proxy) StoreProperties(out io.Writer) (err error) {
	if self.countProperties() > 0 {
		_, err = out.Write([]byte("[proxy]\n"))
		if err != nil {
			return errors.Decorated(err)
		}

		err = self.storeProperties(out)
	}

	return
}
