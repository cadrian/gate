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
	"gate/client/ui"
	"gate/core"
	"gate/core/errors"
	"gate/core/exec"
	"gate/server"
)

import (
	"fmt"
	"io"
	"os"
	osexec "os/exec"
	"strconv"
	"strings"
	"time"
)

var _proxy server.Server

func dirname() (result string) {
	path := strings.Split(os.Args[0], "/")
	if len(path) > 1 {
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
		dir := dirname()
		var exe string
		if dir == "" {
			exe, err = osexec.LookPath("server")
			if err != nil {
				return errors.Decorated(err)
			}
		} else {
			exe = fmt.Sprintf("%s/server", dir)
		}
		var rc string
		if len(os.Args) > 1 {
			rc = os.Args[1]
		}

		writeln := func(pattern string, arg ...interface{}) {
			line := fmt.Sprintf(pattern, arg...) + "\n"
			//fmt.Printf("|%s", line)
			p.Write([]byte(line))
		}

		writeln("%s \"%s\" > /tmp/server-%s.log 2>&1", exe, rc, time.Now().Format("20060102150405"))

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

func readNewMaster(mmi ui.UserInteraction, reason string) (result string, err error) {
	var pass1, pass2, text string

	text = fmt.Sprintf("%s,\nplease enter an encryption phrase.", reason)
	for result == "" {
		pass1, err = mmi.ReadPassword(text)
		if err != nil {
			return
		}
		pass2, err = mmi.ReadPassword("Please enter the same encryption phrase again.")
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

	mmi, err := ui.Ui(srv, config)
	if err != nil {
		return
	}

	var master string
	if vault_info == nil {
		master, err = readNewMaster(mmi, "This is a new vault")
		if err != nil {
			return
		}
	} else {
		master, err = mmi.ReadPassword("Please enter your encryption phrase\nto open the password vault.")
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
