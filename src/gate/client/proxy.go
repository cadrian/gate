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

package client

// Access to the server


import (
	"gate/core"
	"gate/core/errors"
	"gate/core/exec"
	"gate/server"
)

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var _proxy server.Server

func dirname() (result string) {
	path := strings.Split(os.Args[0], "/")
	if len(path) == 1 {
		result = "."
	} else {
		result = strings.Join(path[:len(path)-1], "/")
	}
	return
}

func startServer() (err error) {
	pipe := make(chan io.WriteCloser, 1)

	prepare := func (cmd *exec.Cmd) (err error) {
		p, err := cmd.StdinPipe()
		if err != nil {
			return errors.Decorated(err)
		}
		pipe <- p
		return
	}

	run := func (cmd *exec.Cmd) (err error) {
		p := <-pipe
		p.Write([]byte(fmt.Sprintf("%s/server > /tmp/server.log 2>&1\n", dirname())))

		err = p.Close()
		if err != nil {
			return errors.Decorated(err)
		}
		return
	}

	err = exec.Command(prepare, run, "at", "now")
	if err != nil {
		return
	}

	return
}

func readPassword(config core.Config, text string) (result string, err error) {
	command, err := config.Eval("", "password", "command", os.Getenv)
	if err != nil {
		return
	}
	env := func (name string) string {
		switch name {
		case "TEXT":
			return text
		}
		return ""
	}
	arguments, err := config.Eval("", "password", "arguments", env)
	if err != nil {
		return
	}

	buffer := &bytes.Buffer{}

	type barrierData struct {
		n int64
		err error
	}
	barrier := make(chan barrierData)

	prepare := func (cmd *exec.Cmd) (err error) {
		out, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Decorated(err)
		}

		go func() {
			n, err := buffer.ReadFrom(out)
			barrier <- barrierData{n, err}
		}()
		return
	}

	resulter := make(chan string, 1)

	run := func (cmd *exec.Cmd) (err error) {
		data := <-barrier
		if data.err != nil {
			return errors.Decorated(err)
		}

		// the last character is a \n -- ignore it
		resulter <- string(buffer.Bytes()[:data.n-1])

		return
	}

	err = exec.Command(prepare, run, "bash", "-c", fmt.Sprintf("%s %s", command, arguments))
	if err != nil {
		return
	}

	result = <-resulter
	return
}

func readNewMaster(config core.Config, reason string) (result string, err error) {
	var pass1, pass2, text string

	text = fmt.Sprintf("%s,\nplease enter an encryption phrase.", reason)
	for result == "" {
		pass1, err = readPassword(config, text)
		if err != nil {
			return
		}
		pass2, err = readPassword(config, "Please enter the same encryption phrase again.")
		if err != nil {
			return
		}
		if pass1 == pass2 {
			result = pass1
		} else {
			text = fmt.Sprintf("Your phrases did not match.\n%s,\nplease enter an encryption phrase.", reason)
		}
	}

	return
}

func openVault(srv server.Server, config core.Config) (err error) {
	xdg, err := core.Xdg()
	if err != nil {
		return
	}

	data_home, err := xdg.DataHome()
	if err != nil {
		return
	}
	vault_path := fmt.Sprintf("%s/vault", data_home)
	vault_info, err := os.Stat(vault_path)
	if err != nil {
		return errors.Decorated(err)
	}

	var master string
	if vault_info == nil {
		master, err = readNewMaster(config, "This is a new vault")
		if err != nil {
			return
		}
	} else {
		master, err = readPassword(config, "Please enter your encryption phrase\nto open the password vault.")
		if err != nil {
			return
		}
	}

	var isopen bool
	err = srv.Open(master, &isopen)
	if err != nil {
		return
	}
	if !isopen {
		return errors.New("Could not open vault")
	}

	return
}

func proxy(config core.Config) (result server.Server, err error) {
	result = _proxy
	if result == nil {
		var (
			host, p string
			port int64
			s server.Server
		)
		host, err = config.Eval("", "connection", "host", os.Getenv)
		if err != nil {
			return
		}
		p, err = config.Eval("", "connection", "port", os.Getenv)
		if err != nil {
			return
		}
		port, err = strconv.ParseInt(p, 10, 32)
		if err != nil {
			err = errors.Decorated(err)
			return
		}
		s, err = server.Proxy(host, int(port))
		if err != nil {
			err = startServer()
			if err != nil {
				return
			}
			s, err = server.Proxy(host, int(port))
			if err != nil {
				return
			}
		}

		var isopen bool
		err = s.IsOpen(false, &isopen)
		if err != nil {
			err = errors.Decorated(err)
			return
		}

		if !isopen {
			err = openVault(s, config)
			if err != nil {
				return
			}
		}

		result = s
		_proxy = result
	}
	return
}
