// This file is part of Gate.
// Copyright (C) 2012-2015 Cyril Adrian <cyril.adrian@gmail.com>
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
	"gate/client/remote"
	"gate/client/ui"
	"gate/core"
	"gate/core/errors"
	"gate/server"
)

import (
	"regexp"
	"sort"
)

type CompositeCommand interface {
	Command
	Commander
}

type Commander interface {
	Command(name string) Command
	Default() Command
	Commands(filter string) ([]string, error)
}

type commander struct {
	commands map[string]Command
	defcmd   string
}

var _ Commander = &commander{}

type Command interface {
	Name() string
	Run(line []string) error
	Complete(line []string) ([]string, error)
	Help(line []string) (string, error)
}

type cmd struct {
	commander Commander
	remoter   remote.Remoter
	server    server.Server
	config    core.Config
	mmi       ui.UserInteraction
}

func NewCommander(remoter remote.Remoter, srv server.Server, config core.Config, mmi ui.UserInteraction) (result Commander, err error) {
	cmd := &commander{
		commands: make(map[string]Command),
		defcmd:   "get",
	}
	result = cmd

	cmd.commands["add"] = &cmd_add{result, remoter, srv, config, mmi}
	cmd.commands["del"] = &cmd_del{result, remoter, srv, config, mmi}
	cmd.commands["help"] = &cmd_help{result, remoter, srv, config, mmi}
	cmd.commands["list"] = &cmd_list{result, remoter, srv, config, mmi}
	cmd.commands["load"] = &cmd_load{result, remoter, srv, config, mmi}
	cmd.commands["master"] = &cmd_master{result, remoter, srv, config, mmi}
	cmd.commands["merge"] = &cmd_merge{result, remoter, srv, config, mmi}
	cmd.commands["remote"] = newRemote(result, remoter, srv, config, mmi)
	cmd.commands["save"] = &cmd_save{result, remoter, srv, config, mmi}
	cmd.commands["show"] = &cmd_show{result, remoter, srv, config, mmi}
	cmd.commands["stop"] = &cmd_stop{result, remoter, srv, config, mmi}
	cmd.commands["get"] = &cmd_get{result, remoter, srv, config, mmi}

	return
}

func (self *commander) Command(name string) (result Command) {
	result, _ = self.commands[name]
	return
}

func (self *commander) Default() Command {
	return self.commands[self.defcmd]
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
