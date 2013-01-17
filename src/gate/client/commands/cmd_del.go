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

type cmd_del cmd

var _ Command = &cmd_del{}

func (self *cmd_del) Name() string {
	return "del"
}

func (self *cmd_del) Run(line []string) (err error) {
	return
}

func (self *cmd_del) Complete(line []string) (result []string, err error) {
	return
}

func (self *cmd_del) Help(line []string) (result string, err error) {
	return
}
