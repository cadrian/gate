/*
 * This file is part of Gate.
 * Copyright (C) 2012-2013 Cyril Adrian <cyril.adrian@gmail.com>
 *
 * Gate is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3 of the License.
 *
 * Gate is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.	 See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Gate.  If not, see <http://www.gnu.org/licenses/>.
 */
package client

import (
	"gate/core/errors"
	"gate/core/exec"
	"gate/server"
)

import (
	"bytes"
	"io"
)

func xclip(srv server.Server, out io.Reader, barrier chan error) {
	buffer := &bytes.Buffer{}
	n, err := buffer.ReadFrom(out)
	if err != nil {
		barrier <- err
		return
	}
	name := string(buffer.Bytes()[:n-1])

	var pass string
	err = srv.Get(name, &pass)
	if err != nil {
		barrier <- err
		return
	}

	p := _xclip(pass, barrier, "primary")
	c := _xclip(pass, barrier, "clipboard")
	if p && c {
		barrier <- io.EOF
	}
}

func _xclip(name string, barrier chan error, selection string) bool {
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
		p.Write([]byte(name+"\n"))
		err = p.Close()
		if err != nil {
			return errors.Decorated(err)
		}
		return
	}

	err := exec.Command(prepare, run, "xclip", "-selection", selection)
	if err != nil {
		barrier <- err
		return false
	}

	return true
}
