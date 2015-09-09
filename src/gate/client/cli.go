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

package client

import (
	"gate/client/commands"
	"gate/client/remote"
	"gate/client/ui"
	"gate/core"
)

import (
	"os"
)

// Run the console.
func CommandLine(config core.Config) (err error) {

	srv, err := proxy(config)
	if err != nil {
		return
	}

	remoter := remote.NewRemoter(srv, config)
	if err != nil {
		return
	}

	mmi, err := ui.Ui(srv, config)
	if err != nil {
		return
	}

	commander, err := commands.NewCommander(remoter, srv, config, mmi)
	if err != nil {
		return
	}

	cmd := commander.Command(os.Args[2])
	if cmd == nil {
		cmd = commander.Default()
	}

	err = cmd.Run(os.Args[2:])

	return
}
