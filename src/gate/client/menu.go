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

// Functions specific to the "menu" command

import (
	"gate/core"
	"gate/core/errors"
	"gate/core/exec"
	"gate/client/ui"
	"gate/server"
)

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func clipboard(srv server.Server, out io.Reader, barrier chan error) {
	buffer := &bytes.Buffer{}
	n, err := buffer.ReadFrom(out)
	if err != nil {
		barrier <- err
		return
	}
	name := string(buffer.Bytes()[:n-1])

	err = ui.XclipPassword(srv, name)

	if err == nil {
		err = io.EOF
	}

	barrier <- err
	return
}

func displayMenu(config core.Config, srv server.Server, list []string) (err error) {
	command, err := config.Eval("", "menu", "command", os.Getenv)
	if err != nil {
		return
	}
	arguments, err := config.Eval("", "menu", "arguments", nil)
	if err != nil {
		return
	}

	barrier := make(chan error)
	pipe := make(chan io.WriteCloser, 1)

	prepare := func (cmd *exec.Cmd) (err error) {
		p, err := cmd.StdinPipe()
		if err != nil {
			return errors.Decorated(err)
		}

		out, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Decorated(err)
		}

		go clipboard(srv, out, barrier)

		pipe <- p
		return
	}

	run := func (cmd *exec.Cmd) (err error) {
		p := <-pipe

		for _, entry := range list {
			p.Write([]byte(entry + "\n"))
		}

		err = p.Close()
		if err != nil {
			return errors.Decorated(err)
		}

		e := <-barrier
		if e != io.EOF {
			err = errors.Decorated(e)
		}
		return
	}

	err = exec.Command(prepare, run, "bash", "-c", fmt.Sprintf("%s %s", command, arguments))

	return
}

// Get the list of passwords from the server, displays a list and puts
// the corresponding password in xclip
func Menu(config core.Config) (err error) {
	srv, err := proxy(config)
	if err != nil {
		return
	}
	var list []string
	err = srv.List(".*", &list)
	if err != nil {
		return errors.Decorated(err)
	}
	if len(list) > 0 {
		err = displayMenu(config, srv, list)
	}
	return
}
