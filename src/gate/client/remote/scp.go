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

package remote

// Scp remote

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
	"strings"
)

type scp remote

var _ Remote = &scp{}

var ScpAllowedKeys map[string]bool = map[string]bool{
	"file":    true,
	"host":    true,
	"user":    false,
	"options": false,
}

func newScp(name string, srv server.Server, config core.Config, remoter Remoter) (Remote, error) {
	result := &scp{
		properties{
			allowed:    ScpAllowedKeys,
			properties: make(map[string]string),
		},
		srv,
		remoter,
		name,
		nil,
	}
	file := name + ".rc"
	for key, mandatory := range ScpAllowedKeys {
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

func (self *scp) Name() string {
	return self.name
}

func (self *scp) arguments() (result []string, err error) {
	remote_file := self.getProperty("file")
	url := remote_file
	if remote_file == "" {
		err = errors.Newf("missing remote vault file")
		return
	}
	user := self.getProperty("user")
	host := self.getProperty("host")
	if host != "" {
		if user == "" {
			url = fmt.Sprintf("%s:%s", host, remote_file)
		} else {
			url = fmt.Sprintf("%s@%s:%s", user, host, remote_file)
		}
	} else if user != "" {
		err = errors.Newf("user without host")
		return
	}
	options := self.getProperty("options")
	result = append(strings.Split(options, " "), url)
	return
}

func (self *scp) LoadVault(file string) (err error) {
	args, err := self.arguments()
	if err != nil {
		return
	}

	prepare := func(cmd *exec.Cmd) (err error) {
		if self.proxy != nil {
			self.proxy.Install(cmd)
		}
		cmd.Env = append(os.Environ(), "SSH_ASKPASS=true")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return
	}

	err = exec.Command(prepare, nil, "scp", append(args, file)...)

	return
}

func (self *scp) SaveVault(file string) (err error) {
	args, err := self.arguments()
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

	err = exec.Command(prepare, nil, "scp", append([]string{file}, args...)...)

	return
}

func (self *scp) Proxy() Proxy {
	return self.proxy
}

func (self *scp) SetProperty(key, value string) error {
	return self.setProperty(key, value)
}

func (self *scp) ResetProperty(key string) error {
	return self.resetProperty(key)
}

func (self *scp) StoreProperties(out io.Writer) (err error) {
	_, err = out.Write([]byte("[remote]\nmethod = scp\n"))
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
