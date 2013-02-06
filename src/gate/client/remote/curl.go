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

// Curl remote

import (
	"gate/core"
	"gate/core/errors"
	"gate/core/exec"
	"gate/server"
)

import (
	"fmt"
	"io"
	"os"
)

type curl remote

var _ Remote = &curl{}

var CurlAllowedKeys map[string]bool = map[string]bool{
	"url":         true,
	"user":        false,
	"passkey":     false,
	"put_request": false,
	"get_request": false,
}

func newCurl(name string, srv server.Server, config core.Config, remoter Remoter) (Remote, error) {
	result := &curl{
		properties{
			allowed:    CurlAllowedKeys,
			properties: make(map[string]string),
		},
		srv,
		remoter,
		name,
		nil,
	}
	file := name + ".rc"
	for key, mandatory := range CurlAllowedKeys {
		value, err := config.Eval(file, "remote", key, os.Getenv)
		if err != nil && mandatory {
			return nil, err
		}
		if value != "" {
			result.properties.setProperty(key, value)
		}
	}
	return result, nil
}

func (self *curl) Name() string {
	return self.name
}

func (self *curl) arguments(option, file, request string) (result []string, err error) {
	url := self.getProperty("url")
	if url == "" {
		err = errors.Newf("missing remote vault url")
		return
	}
	result = []string{"-#", option, file, url}
	user := self.getProperty("user")
	passkey := self.getProperty("passkey")
	if user != "" {
		if passkey == "" {
			result = append(result, "-u", user)
		} else {
			var pass string
			err = self.server.Get(passkey, &pass)
			if err != nil {
				return
			}
			result = append(result, "-u", fmt.Sprintf("%s:%s", user, escape_pass_url(pass)))
		}
	}
	if request != "" {
		result = append(result, "--request", request)
	}
	return
}

func (self *curl) doCurl(option, file, request string) (err error) {
	args, err := self.arguments(option, file, request)
	if err != nil {
		return
	}

	prepare := func(cmd *exec.Cmd) (err error) {
		if self.proxy != nil {
			self.proxy.Install(cmd)
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return
	}

	err = exec.Command(prepare, nil, "curl", args...)

	return
}

func (self *curl) LoadVault(file string) error {
	return self.doCurl("-o", file, self.getProperty("get_request"))
}

func (self *curl) SaveVault(file string) error {
	return self.doCurl("-T", file, self.getProperty("put_request"))
}

func (self *curl) Proxy() Proxy {
	return self.proxy
}

func (self *curl) SetProperty(key, value string) error {
	return self.setProperty(key, value)
}

func (self *curl) ResetProperty(key string) error {
	return self.resetProperty(key)
}

func (self *curl) StoreProperties(out io.Writer) (err error) {
	_, err = out.Write([]byte("[remote]\nmethod = curl\n"))
	if err != nil {
		return errors.Decorated(err)
	}

	err = self.storeProperties(out)
	if err != nil {
		return
	}

	if self.proxy != nil {
		err = self.proxy.StoreProperties(out)
	}

	return
}
