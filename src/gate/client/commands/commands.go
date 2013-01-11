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
	"gate/server"
)

type Cmd interface {
	Name() string
	Run(line []string) error
	Complete(line []string, word string) ([]string, error)
	Help(line []string) (string, error)
}

var (
	commands_map map[string]Cmd
)

func Init(srv server.Server) (err error) {
	commands_map = make(map[string]Cmd)
	commands_map["add"] = &cmd_add{srv}
	commands_map["help"] = &cmd_help{srv}
	commands_map["list"] = &cmd_list{srv}
	commands_map["load"] = &cmd_load{srv}
	commands_map["master"] = &cmd_master{srv}
	commands_map["merge"] = &cmd_merge{srv}
	commands_map["rem"] = &cmd_rem{srv}
	commands_map["remote"] = &cmd_remote{srv}
	commands_map["save"] = &cmd_save{srv}
	commands_map["show"] = &cmd_show{srv}
	commands_map["stop"] = &cmd_stop{srv}
	commands_map["get"] = &cmd_get{srv}
	return
}

func Command(name string) (result Cmd) {
	result, ok := commands_map[name]
	if !ok {
		result = commands_map["get"]
	}
	return
}
