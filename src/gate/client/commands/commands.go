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

type Cmd interface {
	Name() string
	Run(line []string) error
	Complete(line []string) ([]string, error)
	Help(line []string) (string, error)
}

type cmd struct {
	server server.Server
	config core.Config
	mmi ui.UserInteraction
}

var (
	commands_map map[string]Cmd
)

func Init(srv server.Server, config core.Config, mmi ui.UserInteraction) (err error) {
	commands_map = make(map[string]Cmd)
	commands_map["add"] = &cmd_add{srv, config, mmi}
	commands_map["help"] = &cmd_help{srv, config, mmi}
	commands_map["list"] = &cmd_list{srv, config, mmi}
	commands_map["load"] = &cmd_load{srv, config, mmi}
	commands_map["master"] = &cmd_master{srv, config, mmi}
	commands_map["merge"] = &cmd_merge{srv, config, mmi}
	commands_map["rem"] = &cmd_rem{srv, config, mmi}
	commands_map["remote"] = &cmd_remote{srv, config, mmi}
	commands_map["save"] = &cmd_save{srv, config, mmi}
	commands_map["show"] = &cmd_show{srv, config, mmi}
	commands_map["stop"] = &cmd_stop{srv, config, mmi}
	commands_map["get"] = &cmd_get{srv, config, mmi}
	return
}

func Command(name string) (result Cmd) {
	result, ok := commands_map[name]
	if !ok {
		result = commands_map["get"]
	}
	return
}

func Commands(filter string) (result []string, err error) {
	re_filter, err := regexp.Compile(filter)
	if err != nil {
		err = errors.Decorated(err)
		return
	}
	result = make([]string, 0, len(commands_map))
	for _, cmd := range commands_map {
		name := cmd.Name()
		if re_filter.MatchString(name) {
			result = append(result, name)
		}
	}
	result = result[:len(result)]
	sort.Strings(result)
	return
}
