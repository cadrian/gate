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
	"gate/core"
	"gate/core/errors"
	"gate/client/remote"
	"gate/client/ui"
	"gate/server"
)

import (
	"fmt"
	"strings"
)

type cmd_remote struct {
	*cmd
	*commander
}

var _ CompositeCommand = &cmd_remote{}

func newRemote(host_commander Commander, remoter remote.Remoter, srv server.Server, config core.Config, mmi ui.UserInteraction) *cmd_remote {
	command := &cmd {
		host_commander,
		remoter,
		srv,
		config,
		mmi,
	}

	cmder := &commander {
		make(map[string]Command),
		"list",
	}

	cmder.commands["list"] = &cmd_remote_list{cmder, remoter, srv, config, mmi}

	return &cmd_remote{command, cmder}
}

func (self *cmd_remote) Name() string {
	return "remote"
}

func (self *cmd_remote) Run(line []string) (err error) {
	var cmd Command
	if len(line) > 1 {
		cmd = self.Command(line[1])
		if cmd == nil {
			err = errors.Newf("unknown remote command: %s", line[0])
			return
		}
	} else {
		cmd = self.Default()
	}
	err = cmd.Run(line)
	return
}

func (self *cmd_remote) Complete(line []string) (result []string, err error) {
	if len(line) == 2 {
		return self.Commands(fmt.Sprintf("^%s", line[1]))
	}
	cmd := self.Command(line[1])
	if cmd == nil {
		return
	}
	return cmd.Complete(line)
}

func (self *cmd_remote) Help(line []string) (result string, err error) {
	var (
		cmd Command
		remotes []string
		commands []string
	)
	if len(line) > 1 {
		cmd = self.Command(line[1])
	}
	if cmd != nil {
		return cmd.Help(line)
	} else {
		var (
			commands_help []string
			remotes_help string
			h string
		)

		commands, err = self.Commands("")
		if err != nil {
			return
		}

		commands_help = make([]string, 0, len(commands))
		for _, name := range commands {
			cmd := self.Command(name)
			h, err = cmd.Help(line)
			if err != nil {
				return
			}
			commands_help = append(commands_help, h)
		}

		remotes, err = self.Command("list").(*cmd_remote_list).listRemotes()
		if err != nil {
			return
		}

		switch len(remotes) {
		case 0:
			remotes_help = "There are no remotes defined."
		case 1:
			remotes_help = fmt.Sprintf("There is only one remote defined: [1m%s[0m", remotes[0])
		default:
			remotes_help = fmt.Sprintf("The defined remotes are:\n		   [1;33m|[0m [1m%s[0m", strings.Join(remotes, ", "))
		}

		result = fmt.Sprintf(`%s
		   [1;33m|[0m [33m[remote][0m note:
		   [1;33m|[0m The [33mload[0m, [33msave[0m, [33mmerge[0m, and [33mremote[0m commands require
		   [1;33m|[0m an extra argument if there is more than one available
		   [1;33m|[0m remotes.
		   [1;33m|[0m In that case, the argument is the remote to select.
		   [1;33m|[0m
		   [1;33m|[0m %s
`,
			strings.Join(commands_help, "\n"),
			remotes_help,
		)
	}

	return
}
