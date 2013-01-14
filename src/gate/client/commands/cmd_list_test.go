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
	"gate/client/ui"
	"gate/core"
	"gate/server"
)

import (
	"code.google.com/p/gomock/gomock"
	"testing"
)

func TestListRun1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := server.NewMockServer(ctrl)
	cfg := core.NewMockConfig(ctrl)
	mmi := ui.NewMockUserInteraction(ctrl)
	list := &cmd_list{srv, cfg, mmi}

	srv.EXPECT().List(".*", gomock.Any()).Do(func (filter string, reply *[]string) {
		*reply = []string{"key1", "key2"}
	})
	mmi.EXPECT().Pager("key1\nkey2\n")

	err := list.Run([]string{"list"})
	if err != nil {
		t.Error(err)
	}
}
