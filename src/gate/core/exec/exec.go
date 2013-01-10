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

package exec

import (
	"gate/core/errors"
)

import (
	"io"
	"os/exec"
)

type Cmd exec.Cmd

func Command(prepare func(cmd *Cmd) error, run func(cmd *Cmd) error, command string, arguments... string) (err error) {
	cmd := exec.Command(command, arguments...)

	if prepare != nil {
		err = prepare((*Cmd)(cmd))
		if err != nil {
			return
		}
	}

	err = cmd.Start()
	if err != nil {
		return errors.Decorated(err)
	}

	if run != nil {
		err = run((*Cmd)(cmd))
		if err != nil {
			return
		}
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Decorated(err)
	}

	return
}

func (self *Cmd) StdinPipe() (io.WriteCloser, error) {
	return ((*exec.Cmd)(self)).StdinPipe()
}

func (self *Cmd) StdoutPipe() (io.ReadCloser, error) {
	return ((*exec.Cmd)(self)).StdoutPipe()
}

func (self *Cmd) StderrPipe() (io.ReadCloser, error) {
	return ((*exec.Cmd)(self)).StderrPipe()
}
