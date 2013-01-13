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

import (
	"gate/client/commands"
	"gate/core"
	"gate/core/errors"
	"gate/server"
)

import (
	"fmt"
	"github.com/sbinet/liner"
	"io"
	"strings"
)

type readline struct {
	server server.Server
	state *liner.State
	lastline string
}

func (self *readline) run(line []string) (err error) {
	cmd := commands.Command(line[0])
	return cmd.Run(line)
}

func (self *readline) complete(line string) (result []string, err error) {
	words := strings.Split(line, " ")
	if len(words) == 1 {
		return commands.Commands(fmt.Sprintf("^%s", words[0]))
	}
	cmd := commands.Command(words[0])
	candidates, err := cmd.Complete(words)
	if err != nil {
		return
	}

	result = make([]string, 0, len(candidates))
	n := len(words) - 1
	for _, candidate := range candidates {
		words[n] = candidate
		result = append(result, strings.Join(words, " "))
	}

	return
}

func (self *readline) loop(config core.Config) (err error) {
	var line string

	done := false
	for !done {
		line, err = self.state.Prompt("> ")
		if err == nil && len(line) > 0 {
			err = self.run(strings.Split(line, " "))
			if err != nil {
				e, ok := err.(errors.StackError)
				if !ok {
					return
				}
				fmt.Println(strings.Split(e.Nested.Error(), "\n")[0])
			}
			if line != self.lastline {
				self.state.AppendHistory(line)
				self.lastline = line
			}
		} else {
			done = true
		}
	}

	fmt.Println()
	return
}

// Run the console.
func Console(config core.Config) (err error) {
	fmt.Printf(`
[1;32mWelcome to the Gate administration console![0m

[32mGate Copyright (C) 2012-2013 Cyril Adrian <cyril.adrian@gmail.com>
This program comes with ABSOLUTELY NO WARRANTY; for details type [33mshow w[32m.
This is free software, and you are welcome to redistribute it
under certain conditions; type [33mshow c[32m for details.[0m

Type [33mhelp[0m for details on available options.
Just hit [33m<enter>[0m to exit.

`)

	srv, err := proxy(config)
	if err != nil {
		return
	}

	commands.Init(srv, config)

	state := liner.NewLiner()
	defer state.Close()

	rl := &readline {
		server: srv,
		state: state,
	}

	complete := func (line string) (result []string) {
		result, err := rl.complete(line)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	state.SetCompleter(complete)

	e := rl.loop(config)
	if e != io.EOF {
		err = e
	}
	return
}
