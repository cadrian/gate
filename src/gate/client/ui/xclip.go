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

package ui

// copy to clipboard

import (
	"gate/core/errors"
	"gate/core/exec"
)

import (
	"io"
)

// Copy the data string into the X clipboard (both primary and clipboard)
func (self *interaction) Xclip(data string) (err error) {
	err = xclip(data, "primary")
	if err != nil {
		return
	}

	err = xclip(data, "clipboard")
	if err != nil {
		return
	}

	return
}

// Fetch the password from the server and xclips the corresponding password
func (self *interaction) XclipPassword(name string) (err error) {
	var pass string
	err = self.server.Get(name, &pass)
	if err != nil {
		return
	}

	err = self.Xclip(pass)

	return
}

func xclip(name string, selection string) (err error) {
	pipe := make(chan io.WriteCloser, 1)

	prepare := func(cmd *exec.Cmd) (err error) {
		p, err := cmd.StdinPipe()
		if err != nil {
			return errors.Decorated(err)
		}
		pipe <- p
		return
	}

	run := func(cmd *exec.Cmd) (err error) {
		p := <-pipe
		p.Write([]byte(name))
		err = p.Close()
		if err != nil {
			return errors.Decorated(err)
		}
		return
	}

	err = exec.Command(prepare, run, "xclip", "-selection", selection)
	return
}
