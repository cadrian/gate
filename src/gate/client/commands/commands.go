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
	"gate/client/ui"
	"gate/core"
	"gate/core/errors"
	"gate/server"
)

import (
	"regexp"
	"sort"
)

type Commander interface {
	Command(name string) Cmd
	Commands(filter string) ([]string, error)
}

type commander struct {
	commands map[string]Cmd
}

var _ Commander = &commander{}

type Cmd interface {
	Name() string
	Run(line []string) error
	Complete(line []string) ([]string, error)
	Help(line []string) (string, error)
}

type cmd struct {
	commander Commander
	server server.Server
	config core.Config
	mmi ui.UserInteraction
}

func NewCommander(srv server.Server, config core.Config, mmi ui.UserInteraction) (result Commander, err error) {
	cmd := &commander{
		commands: make(map[string]Cmd),
	}

	cmd.commands["add"] = &cmd_add{result, srv, config, mmi}
	cmd.commands["help"] = &cmd_help{result, srv, config, mmi}
	cmd.commands["list"] = &cmd_list{result, srv, config, mmi}
	cmd.commands["load"] = &cmd_load{result, srv, config, mmi}
	cmd.commands["master"] = &cmd_master{result, srv, config, mmi}
	cmd.commands["merge"] = &cmd_merge{result, srv, config, mmi}
	cmd.commands["rem"] = &cmd_rem{result, srv, config, mmi}
	cmd.commands["remote"] = &cmd_remote{result, srv, config, mmi}
	cmd.commands["save"] = &cmd_save{result, srv, config, mmi}
	cmd.commands["show"] = &cmd_show{result, srv, config, mmi}
	cmd.commands["stop"] = &cmd_stop{result, srv, config, mmi}
	cmd.commands["get"] = &cmd_get{result, srv, config, mmi}

	result = cmd

	return
}

func (self *commander) Command(name string) (result Cmd) {
	result, ok := self.commands[name]
	if !ok {
		result = self.commands["get"]
	}
	return
}

func (self *commander) Commands(filter string) (result []string, err error) {
	re_filter, err := regexp.Compile(filter)
	if err != nil {
		err = errors.Decorated(err)
		return
	}
	result = make([]string, 0, len(self.commands))
	for _, cmd := range self.commands {
		name := cmd.Name()
		if re_filter.MatchString(name) {
			result = append(result, name)
		}
	}
	result = result[:len(result)]
	sort.Strings(result)
	return
}
