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
	"gate/core/errors"
)

import (
	"fmt"
)

type cmd_save cmd

var _ Command = &cmd_save{}

func (self *cmd_save) Name() string {
	return "save"
}

func (self *cmd_save) Run(line []string) (err error) {
	var remoteName string
	if len(line) > 1 {
		remoteName = line[1]
	} else {
		remoteName = ""
	}

	remote, err := self.remoter.Remote(remoteName)
	if err != nil {
		return
	}

	vault_path, err := self.config.VaultPath()
	if err != nil {
		return
	}

	var saved bool
	err = self.server.Save(true, &saved)
	if err != nil {
		return
	}
	if !saved {
		err = errors.Newf("Could not save vault")
		return
	}

	err = remote.SaveVault(vault_path)

	return
}

func (self *cmd_save) Complete(line []string) (result []string, err error) {
	return
}

func (self *cmd_save) Help(line []string) (result string, err error) {
	var remote_note string
	if len(line) == 1 {
		remote_note = "the [1mremote note[0m above"
	} else {
		remote_note = "note perusing [1mhelp remote[0m"
	}

	result = fmt.Sprintf(`
[33msave [remote][0m      Save the password vault upto the server.
		   [33m[remote][0m: see %s
`,
		remote_note,
	)

	return
}
