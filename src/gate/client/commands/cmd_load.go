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

type cmd_load cmd

var _ Command = &cmd_load{}

func (self *cmd_load) Name() string {
	return "load"
}

func (self *cmd_load) Run(line []string) (err error) {
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

	err = remote.LoadVault(vault_path)

	return
}

func (self *cmd_load) Complete(line []string) (result []string, err error) {
	return
}

func (self *cmd_load) Help(line []string) (result string, err error) {

	result = `
[33mload [remote][0m      [1mReplace[0m the local vault with the server's version.
		   Note: in that case you will be asked for the new vault
		   password (the previous vault is closed).
		   [33m[remote][0m: see note below
`

	return
}
