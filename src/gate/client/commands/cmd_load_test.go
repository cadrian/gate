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

func TestLoadRun1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := NewMockCommander(ctrl)
	rem := remote.NewMockRemoter(ctrl)
	srv := server.NewMockServer(ctrl)
	cfg := core.NewMockConfig(ctrl)
	mmi := ui.NewMockUserInteraction(ctrl)
	load := &cmd_load{cmd, rem, srv, cfg, mmi}

	rmt := remote.NewMockRemote(ctrl)
	rem.EXPECT().Remote("").Return(rmt, nil)

	path := "vault_path"
	cfg.EXPECT().VaultPath().Return(path, nil)

	rmt.EXPECT().LoadVault(path)

	err := load.Run([]string{"load"})
	if err != nil {
		t.Error(err)
	}
}

func TestLoadRun2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := NewMockCommander(ctrl)
	rem := remote.NewMockRemoter(ctrl)
	srv := server.NewMockServer(ctrl)
	cfg := core.NewMockConfig(ctrl)
	mmi := ui.NewMockUserInteraction(ctrl)
	load := &cmd_load{cmd, rem, srv, cfg, mmi}

	rmt := remote.NewMockRemote(ctrl)
	rem.EXPECT().Remote("foo").Return(rmt, nil)

	path := "vault_path"
	cfg.EXPECT().VaultPath().Return(path, nil)

	rmt.EXPECT().LoadVault(path)

	err := load.Run([]string{"load", "foo"})
	if err != nil {
		t.Error(err)
	}
}
