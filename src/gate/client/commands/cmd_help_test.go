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
	"gate/client/remote"
	"gate/client/ui"
	"gate/core"
	"gate/server"
)

import (
	"code.google.com/p/gomock/gomock"
	"testing"
)

func TestHelpRun1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := NewMockCommander(ctrl)
	rem := remote.NewMockRemoter(ctrl)
	srv := server.NewMockServer(ctrl)
	cfg := core.NewMockConfig(ctrl)
	mmi := ui.NewMockUserInteraction(ctrl)
	help := &cmd_help{cmd, rem, srv, cfg, mmi}

	cmd_foo := NewMockCommand(ctrl)
	cmd_bar := NewMockCommand(ctrl)

	cmd.EXPECT().Commands(".*").Return([]string{"foo", "bar"}, nil)
	cmd.EXPECT().Command("foo").Return(cmd_foo)
	cmd.EXPECT().Command("bar").Return(cmd_bar)

	cmd_foo.EXPECT().Help([]string{"help"}).Return("\nhelp for foo\n", nil)
	cmd_bar.EXPECT().Help([]string{"help"}).Return("\nhelp for bar\n", nil)

	mmi.EXPECT().Pager(`
[1;32mKnown commands[0m

help for foo

help for bar

Any other input is understood as a password request using the given key.
If that key exists the password is stored in the clipboard.

[1m--------[0m
[32mGate Copyright (C) 2012-2013 Cyril Adrian <cyril.adrian@gmail.com>[0m
[32mThis program comes with ABSOLUTELY NO WARRANTY; for details type [33mshow w[32m.[0m
[32mThis is free software, and you are welcome to redistribute it[0m
[32munder certain conditions; type [33mshow c[32m for details.[0m
`)

	err := help.Run([]string{"help"})
	if err != nil {
		t.Error(err)
	}
}

func TestHelpRun2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := NewMockCommander(ctrl)
	rem := remote.NewMockRemoter(ctrl)
	srv := server.NewMockServer(ctrl)
	cfg := core.NewMockConfig(ctrl)
	mmi := ui.NewMockUserInteraction(ctrl)
	help := &cmd_help{cmd, rem, srv, cfg, mmi}

	cmd_foo := NewMockCommand(ctrl)

	cmd.EXPECT().Command("foo").Return(cmd_foo)

	cmd_foo.EXPECT().Help([]string{"help", "foo"}).Return("\nhelp for foo\n", nil)

	mmi.EXPECT().Pager(`
[1;32mHelp for command foo[0m

help for foo

[1m--------[0m
[32mGate Copyright (C) 2012-2013 Cyril Adrian <cyril.adrian@gmail.com>[0m
[32mThis program comes with ABSOLUTELY NO WARRANTY; for details type [33mshow w[32m.[0m
[32mThis is free software, and you are welcome to redistribute it[0m
[32munder certain conditions; type [33mshow c[32m for details.[0m
`)

	err := help.Run([]string{"help", "foo"})
	if err != nil {
		t.Error(err)
	}
}
