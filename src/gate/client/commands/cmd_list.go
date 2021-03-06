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
	"strings"
)

type cmd_list cmd

var _ Command = &cmd_list{}

func (self *cmd_list) Name() string {
	return "list"
}

func (self *cmd_list) Run(line []string) (err error) {
	var reply []string
	var filter string
	if len(line) > 1 {
		filter = line[1]
	} else {
		filter = ".*"
	}
	err = self.server.List(filter, &reply)
	if err != nil {
		return
	}

	err = self.mmi.Pager(strings.Join(append(reply, ""), "\n"))
	return
}

func (self *cmd_list) Complete(line []string) (result []string, err error) {
	return
}

func (self *cmd_list) Help(line []string) (result string, err error) {

	result = `
[33mlist[0m		   List the known passwords (show only the keys).
`

	return
}
