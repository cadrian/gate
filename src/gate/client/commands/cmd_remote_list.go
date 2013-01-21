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
	"strings"
)

type cmd_remote_list cmd

var _ Command = &cmd_remote_list{}

func (self *cmd_remote_list) Name() string {
	return "list"
}

func (self *cmd_remote_list) listRemotes() (result []string, err error) {
	files, err := self.config.ListConfigFiles()
	if err != nil {
		return
	}
	result = make([]string, 0, len(files))
	for _, file := range files {
		result = append(result, file[0:len(file)-3])
	}
	return
}

func (self *cmd_remote_list) Run(line []string) (err error) {
	remotes, err := self.listRemotes()
	if err != nil {
		return
	}
	self.mmi.Pager(strings.Join(remotes, "\n"))
	return
}

func (self *cmd_remote_list) Complete(line []string) (result []string, err error) {
	return
}

func (self *cmd_remote_list) Help(line []string) (result string, err error) {
	result = `
[33mremote list[0m	   Lists the known remotes.
`
	return
}
