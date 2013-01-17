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
	"gate/client/ui"
	"gate/server"
)

import (
	"fmt"
)

type cmd_remote struct {
	*cmd
	*commander
}

var _ CompositeCommand = &cmd_remote{}

func newRemote(host_commander Commander, srv server.Server, config core.Config, mmi ui.UserInteraction) *cmd_remote {
	command := &cmd {
		host_commander,
		srv,
		config,
		mmi,
	}

	cmder := &commander {
		make(map[string]Command),
		"list",
	}

	cmder.commands["list"] = &cmd_remote_list{cmder, srv, config, mmi}

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
	return
}
