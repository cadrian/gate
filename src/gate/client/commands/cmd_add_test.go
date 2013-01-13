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
	"gate/mocks"
	"gate/server"
)

import (
	"code.google.com/p/gomock/gomock"
	"testing"
)

func TestAddRun2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mocks.NewMockServer(ctrl)
	cfg := mocks.NewMockConfig(ctrl)

	mmi := mocks.NewMockUserInteraction(ctrl)
	add := &cmd_add{srv, cfg, mmi}

	cfg.EXPECT().Eval("", "console", "default_recipe", gomock.Any()).Return("recipe", nil)
	args := server.SetArgs{Key:"foo", Recipe:"recipe"}

	srv.EXPECT().Set(args, gomock.Any()).Do(func (_ server.SetArgs, pass *string) {
		*pass = "password"
	})

	mmi.EXPECT().Xclip("password")

	err := add.Run([]string{"add", "foo"})
	if err != nil {
		t.Error(err)
	}
}

func TestAddRun3Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mocks.NewMockServer(ctrl)
	cfg := mocks.NewMockConfig(ctrl)

	mmi := mocks.NewMockUserInteraction(ctrl)
	add := &cmd_add{srv, cfg, mmi}

	cfg.EXPECT().Eval("", "console", "default_recipe", gomock.Any()).Return("recipe", nil)
	args := server.SetArgs{Key:"foo", Recipe:"recipe"}

	srv.EXPECT().Set(args, gomock.Any()).Do(func (_ server.SetArgs, pass *string) {
		*pass = "password"
	})

	mmi.EXPECT().Xclip("password")

	err := add.Run([]string{"add", "foo", "generate"})
	if err != nil {
		t.Error(err)
	}
}

func TestAddRun3Prompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mocks.NewMockServer(ctrl)
	cfg := mocks.NewMockConfig(ctrl)

	mmi := mocks.NewMockUserInteraction(ctrl)
	add := &cmd_add{srv, cfg, mmi}

	mmi.EXPECT().ReadPassword("Please enter the new password for foo").Return("passwd", nil)
	args := server.SetArgs{Key:"foo", Pass:"passwd"}

	srv.EXPECT().Set(args, gomock.Any()).Do(func (_ server.SetArgs, pass *string) {
		*pass = "password"
	})

	mmi.EXPECT().Xclip("password")

	err := add.Run([]string{"add", "foo", "prompt"})
	if err != nil {
		t.Error(err)
	}
}

func TestAddRun4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mocks.NewMockServer(ctrl)
	cfg := mocks.NewMockConfig(ctrl)

	mmi := mocks.NewMockUserInteraction(ctrl)
	add := &cmd_add{srv, cfg, mmi}

	args := server.SetArgs{Key:"foo", Recipe:"recipe"}

	srv.EXPECT().Set(args, gomock.Any()).Do(func (_ server.SetArgs, pass *string) {
		*pass = "password"
	})

	mmi.EXPECT().Xclip("password")

	err := add.Run([]string{"add", "foo", "generate", "recipe"})
	if err != nil {
		t.Error(err)
	}
}
