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

// Some execution facilities, because Gate relies a lot on process spawning
package exec

import (
	"gate/core/errors"
)

import (
	"io"
	"os/exec"
)

type Cmd exec.Cmd

// Spawn a given command.
//  - prepare() is called after the command creation but before actually starting it
//  - run() is called while the command is running and before waiting for its completion
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

// A pipe to the command's stdin
func (self *Cmd) StdinPipe() (io.WriteCloser, error) {
	return ((*exec.Cmd)(self)).StdinPipe()
}

// A pipe to the command's stdout
func (self *Cmd) StdoutPipe() (io.ReadCloser, error) {
	return ((*exec.Cmd)(self)).StdoutPipe()
}

// A pipe to the command's stderr
func (self *Cmd) StderrPipe() (io.ReadCloser, error) {
	return ((*exec.Cmd)(self)).StderrPipe()
}
