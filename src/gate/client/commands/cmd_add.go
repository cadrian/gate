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
	"gate/core/errors"
	"gate/server"
)

import (
	"fmt"
	"os"
)

type cmd_add cmd

var _ Cmd = &cmd_add{}

func (self *cmd_add) Name() string {
	return "add"
}

func (self *cmd_add) generateArgs(key string, recipe string) (result server.SetArgs, err error) {
	if recipe == "" {
		recipe, err = self.config.Eval("", "console", "default_recipe", os.Getenv)
		if err != nil {
			return
		}
	}
	result = server.SetArgs{
		Key: key,
		Recipe: recipe,
	}
	return
}

func (self *cmd_add) promptArgs(key string) (result server.SetArgs, err error) {
	pass, err := self.mmi.ReadPassword(fmt.Sprintf("Please enter the new password for %s", key))
	if err != nil {
		return
	}
	result = server.SetArgs{
		Key: key,
		Pass: pass,
	}
	return
}

func (self *cmd_add) Run(line []string) (err error) {
	var arg server.SetArgs
	switch len(line) {
	case 2:
		arg, err = self.generateArgs(line[1], "")
	case 3:
		action := line[2]
		switch action {
		case "generate":
			arg, err = self.generateArgs(line[1], "")
		case "prompt":
			arg, err = self.promptArgs(line[1])
		default:
			err = errors.Newf("Unrecognized argument: '%s'", action)
		}
	case 4:
		recipe := line[3]
		action := line[2]
		switch action {
		case "generate":
			arg, err = self.generateArgs(line[1], recipe)
		default:
			err = errors.Newf("Unrecognized argument: '%s'", action)
		}
	default:
		err = errors.New("Invalid arguments")
	}
	if err != nil {
		return
	}

	var pass string
	err = self.server.Set(arg, &pass)
	if err != nil {
		return
	}

	err = self.mmi.Xclip(pass)
	return
}

func (self *cmd_add) Complete(line []string) (result []string, err error) {
	return
}

func (self *cmd_add) Help(line []string) (result string, err error) {

	result = `
[33madd <key> [how][0m    Add a new password. Needs at least a key.
		   If [33m[how][0m is "generate" then the password is
		   randomly generated ([1mdefault[0m).
		   If [33m[how][0m is "generate" with an extra argument then
		   the extra argument represents a "recipe" used to generate
		   the password (*).
		   If [33m[how][0m is "prompt" then the password is asked.
		   If the password already exists it is changed.
		   In all cases the password is stored in the clipboard.

		   (*) A recipe is a series of "ingredients" separated by a '+'.
		   Each "ingredient" is an optional quantity (default 1)
		   followed by a series of 'a' (alphanumeric), 'n' (numeric),
		   or 's' (symbol).
		   The password is generated using the recipe to randomly select
		   characters, and mixing them.
`

	return
}
