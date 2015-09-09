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
	"gate/server"
)

import (
	"fmt"
)

type cmd_merge cmd

var _ Command = &cmd_merge{}

func (self *cmd_merge) Name() string {
	return "merge"
}

func (self *cmd_merge) Run(line []string) (err error) {
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

	xdg, err := self.config.Xdg()
	if err != nil {
		return
	}

	dir, err := xdg.RuntimeDir()
	if err != nil {
		return
	}

	merge_vault := fmt.Sprintf("%s/merge_vault", dir)

	err = remote.LoadVault(merge_vault)
	if err != nil {
		return
	}

	pass, err := self.mmi.ReadPassword(`Please enter the encryption phrase
to the remote vault`)
	if err != nil {
		return
	}

	if pass != "" {
		var merged bool
		err = self.server.Merge(server.MergeArgs{merge_vault, pass}, &merged)
		if err != nil {
			return
		}
		if !merged {
			err = errors.Newf("Could not merge %s", merge_vault)
			return
		}

		cmd := self.commander.Command("save")
		err = cmd.Run(line)
		if err != nil {
			return
		}

		var vault_path string
		vault_path, err = self.config.VaultPath()
		if err != nil {
			return
		}

		err = remote.SaveVault(vault_path)
	}

	return
}

func (self *cmd_merge) Complete(line []string) (result []string, err error) {
	return
}

func (self *cmd_merge) Help(line []string) (result string, err error) {
	return
}
