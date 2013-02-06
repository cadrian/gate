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
	"bytes"
	"fmt"
	"os"
)

func (self *interaction) ReadPassword(text string) (result string, err error) {
	command, err := self.config.Eval("", "password", "command", os.Getenv)
	if err != nil {
		return
	}
	env := func(name string) string {
		switch name {
		case "TEXT":
			return text
		}
		return ""
	}
	arguments, err := self.config.Eval("", "password", "arguments", env)
	if err != nil {
		return
	}

	buffer := &bytes.Buffer{}

	type barrierData struct {
		n   int64
		err error
	}
	barrier := make(chan barrierData)

	prepare := func(cmd *exec.Cmd) (err error) {
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

	run := func(cmd *exec.Cmd) (err error) {
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
