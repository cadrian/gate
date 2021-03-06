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
	"fmt"
)

type cmd_get cmd

var _ Command = &cmd_get{}

func (self *cmd_get) Name() string {
	return "get"
}

func (self *cmd_get) Run(line []string) (err error) {
	err = self.mmi.XclipPassword(line[len(line)-1])
	return
}

func (self *cmd_get) Complete(line []string) (result []string, err error) {
	var word string
	if len(line) > 1 {
		word = line[len(line)-1]
	}
	err = self.server.List(fmt.Sprintf("^%s", word), &result)
	return
}

func (self *cmd_get) Help(line []string) (result string, err error) {

	result = `
[33mget <key>[0m	   Get a password using the given key.
		   If that key exists the password is stored in the clipboard.

`

	return
}
