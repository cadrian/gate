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

type cmd_master cmd

var _ Command = &cmd_master{}

func (self *cmd_master) Name() string {
	return "master"
}

func (self *cmd_master) Run(line []string) (err error) {
	pass1, err := self.mmi.ReadPassword("Please enter the new master password")
	if err != nil {
		return
	}
	if pass1 == "" {
		err = errors.Newf("Cancelled")
		return
	}

	pass2, err := self.mmi.ReadPassword("Please enter the new master password (again)")
	if err != nil {
		return
	}

	if pass1 != pass2 {
		err = errors.Newf("Passwords don't match")
		return
	}

	var changed bool
	err = self.server.SetMaster(pass1, &changed)
	if err != nil {
		return
	}

	if !changed {
		err = errors.Newf("Could not change master")
	}

	return
}

func (self *cmd_master) Complete(line []string) (result []string, err error) {
	return
}

func (self *cmd_master) Help(line []string) (result string, err error) {
	return
}
