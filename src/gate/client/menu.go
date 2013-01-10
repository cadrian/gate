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
	"gate/core"
	"gate/core/errors"
	"gate/core/exec"
)

import (
	"fmt"
	"io"
	"os"
)

func displayMenu(config core.Config, list []string) (err error) {
	command, err := config.Eval("", "menu", "command", os.Getenv)
	if err != nil {
		return
	}
	arguments, err := config.Eval("", "menu", "arguments", nil)
	if err != nil {
		return
	}

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

		for _, entry := range list {
			p.Write([]byte(entry + "\n"))
		}

		err = p.Close()
		if err != nil {
			return errors.Decorated(err)
		}
		return
	}

	err = exec.Command(prepare, run, "bash", "-c", fmt.Sprintf("%s %s", command, arguments))

	return
}

func Menu(config core.Config) (err error) {
	server, err := proxy(config)
	if err != nil {
		return
	}
	var list []string
	err = server.List(".*", &list)
	if err != nil {
		return errors.Decorated(err)
	}
	if len(list) > 0 {
		err = displayMenu(config, list)
	}
	return
}
