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

package commands

import (
	"fmt"
	"strings"
)

type cmd_help cmd

var _ Command = &cmd_help{}

func (self *cmd_help) Name() string {
	return "help"
}

func (self *cmd_help) Run(line []string) (err error) {

	var (
		header, extra string
		help          []string
	)

	if len(line) == 1 {
		header = "[1;32mKnown commands[0m"
		commands, err := self.commander.Commands(".*")
		if err != nil {
			return err
		}

		help = make([]string, 0, len(commands))
		for _, c := range commands {
			cmd := self.commander.Command(c)
			hlp, e := cmd.Help(line)
			if e != nil {
				return e
			}
			help = append(help, hlp)
			extra = `
Any other input is understood as a [33mget[0m command of the entry as key.
`
		}
	} else {
		c := line[1]
		header = fmt.Sprintf("[1;32mHelp for command %s[0m", c)
		cmd := self.commander.Command(c)
		hlp, err := cmd.Help(line)
		if err != nil {
			return err
		}

		help = []string{hlp}
	}

	self.mmi.Pager(fmt.Sprintf(`
%s
%s%s
[1m--------[0m
[32mGate Copyright (C) 2012-2013 Cyril Adrian <cyril.adrian@gmail.com>[0m
[32mThis program comes with ABSOLUTELY NO WARRANTY; for details type [33mshow w[32m.[0m
[32mThis is free software, and you are welcome to redistribute it[0m
[32munder certain conditions; type [33mshow c[32m for details.[0m
`,
		header,
		strings.Join(help, ""),
		extra,
	))

	return
}

func (self *cmd_help) Complete(line []string) (result []string, err error) {
	return
}

func (self *cmd_help) Help(line []string) (result string, err error) {

	result = `
[33mhelp[0m		   Show this screen
`

	return
}
