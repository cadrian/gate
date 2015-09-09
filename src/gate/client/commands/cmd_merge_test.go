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
	"gate/client/remote"
	"gate/client/ui"
	"gate/core"
	"gate/server"
)

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestMergeRun1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := NewMockCommander(ctrl)
	rem := remote.NewMockRemoter(ctrl)
	srv := server.NewMockServer(ctrl)
	cfg := core.NewMockConfig(ctrl)
	mmi := ui.NewMockUserInteraction(ctrl)
	merge := &cmd_merge{cmd, rem, srv, cfg, mmi}

	rmt := remote.NewMockRemote(ctrl)
	rem.EXPECT().Remote("").Return(rmt, nil)

	xdg := core.NewMockXdgContext(ctrl)
	cfg.EXPECT().Xdg().Return(xdg, nil)
	xdg.EXPECT().RuntimeDir().Return("runtimeDir", nil)

	rmt.EXPECT().LoadVault("runtimeDir/merge_vault").Return(nil)

	pass := "remote pass"
	mmi.EXPECT().ReadPassword(gomock.Any()).Return(pass, nil)

	srv.EXPECT().Merge(server.MergeArgs{"runtimeDir/merge_vault", pass}, gomock.Any()).Do(func(_ server.MergeArgs, reply *bool) {
		*reply = true
	})

	save := NewMockCommand(ctrl)
	cmd.EXPECT().Command("save").Return(save)
	save.EXPECT().Run([]string{"merge"})

	path := "vault_path"
	cfg.EXPECT().VaultPath().Return(path, nil)
	rmt.EXPECT().SaveVault(path)

	err := merge.Run([]string{"merge"})
	if err != nil {
		t.Error(err)
	}
}

func TestMergeRun2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := NewMockCommander(ctrl)
	rem := remote.NewMockRemoter(ctrl)
	srv := server.NewMockServer(ctrl)
	cfg := core.NewMockConfig(ctrl)
	mmi := ui.NewMockUserInteraction(ctrl)
	merge := &cmd_merge{cmd, rem, srv, cfg, mmi}

	rmt := remote.NewMockRemote(ctrl)
	rem.EXPECT().Remote("foo").Return(rmt, nil)

	xdg := core.NewMockXdgContext(ctrl)
	cfg.EXPECT().Xdg().Return(xdg, nil)
	xdg.EXPECT().RuntimeDir().Return("runtimeDir", nil)

	rmt.EXPECT().LoadVault("runtimeDir/merge_vault").Return(nil)

	pass := "remote pass"
	mmi.EXPECT().ReadPassword(gomock.Any()).Return(pass, nil)

	srv.EXPECT().Merge(server.MergeArgs{"runtimeDir/merge_vault", pass}, gomock.Any()).Do(func(_ server.MergeArgs, reply *bool) {
		*reply = true
	})

	save := NewMockCommand(ctrl)
	cmd.EXPECT().Command("save").Return(save)
	save.EXPECT().Run([]string{"merge", "foo"})

	path := "vault_path"
	cfg.EXPECT().VaultPath().Return(path, nil)
	rmt.EXPECT().SaveVault(path)

	err := merge.Run([]string{"merge", "foo"})
	if err != nil {
		t.Error(err)
	}
}
