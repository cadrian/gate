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
)

import (
	"fmt"
	"github.com/sbinet/liner"
	"strings"
)

func run(line []string) (err error) {
	cmd := commands.Command(line[0])
	return cmd.Run(line)
}

func loop(config core.Config) (err error) {
	srv, err := proxy(config)

	commands.Init(srv)

	state := liner.NewLiner()
	defer state.Close()

	var line string

	done := false
	for !done {
		line, err = state.Prompt("> ")
		if err == nil && len(line) > 0 {
			err = run(strings.Split(line, " "))
			if err != nil {
				return
			}
			state.AppendHistory(line)
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
[1;32mWelcome to the pwdmgr administration console![0m

[32mpwdmgr Copyright (C) 2012 Cyril Adrian <cyril.adrian@gmail.com>
This program comes with ABSOLUTELY NO WARRANTY; for details type [33mshow w[32m.
This is free software, and you are welcome to redistribute it
under certain conditions; type [33mshow c[32m for details.[0m

Type [33mhelp[0m for details on available options.
Just hit [33m<enter>[0m to exit.

`)

	err = loop(config)
	return
}
