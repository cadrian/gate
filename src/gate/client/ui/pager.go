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

package ui

import (
	"gate/core/errors"
	"gate/core/exec"
)

import (
	"io"
	"os"
)

func (self *interaction) Pager(text string) (err error) {
	pipe := make(chan io.WriteCloser, 1)

	prepare := func(cmd *exec.Cmd) (err error) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		p, err := cmd.StdinPipe()
		if err != nil {
			return errors.Decorated(err)
		}
		pipe <- p
		return
	}

	run := func(cmd *exec.Cmd) (err error) {
		p := <-pipe
		p.Write([]byte(text))
		err = p.Close()
		if err != nil {
			return errors.Decorated(err)
		}
		return
	}

	err = exec.Command(prepare, run, "less", "-R")
	return
}
