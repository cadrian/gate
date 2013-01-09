/*
 * This file is part of Gate.
 * Copyright (C) 2012-2013 Cyril Adrian <cyril.adrian@gmail.com>
 *
 * Gate is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3 of the License.
 *
 * Gate is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.	 See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Gate.  If not, see <http://www.gnu.org/licenses/>.
 */
package server

import (
	"gate/core"
	"gate/core/errors"
)

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type MergeArgs struct {
	Vault string
	Master string
}

type SetArgs struct {
	Key string
	Pass string
	Recipe string
}

type Server interface {
	Open(master string, reply *bool) error
	IsOpen(thenClose bool, reply *bool) error
	Get(key string, reply *string) error
	Set(args SetArgs, reply *string) error
	Unset(key string, reply *bool) error
	List(filter string, reply *[]string) error
	Merge(args MergeArgs, reply *bool) error
	Save(force bool, reply *bool) error
	Stop(status int, reply *bool) error
}

type server struct {
	vault Vault
	config core.Config
	listener net.Listener
}

var _ Server = &server{}

func newVault(file string) Vault {
	in := func() (result io.ReadCloser, err error) {
		return os.Open(file)
	}
	out := func() (result io.WriteCloser, err error) {
		return os.Create(file)
	}
	return NewVault(in, out)
}

func Start(config core.Config, port int) (result Server, err error) {
	xdg, err := core.Xdg()
	if err != nil {
		return
	}
	data_home, err := xdg.DataHome()
	if err != nil {
		return
	}
	vault_path := fmt.Sprintf("%s/vault", data_home)
	srv := &server{
		vault: newVault(vault_path),
		config: config,
	}
	rpc.Register(srv)
	rpc.HandleHTTP()
	srv.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		err = errors.Decorated(err)
		return
	}
	go http.Serve(srv.listener, nil)
	result = srv
	return
}

func (self *server) IsOpen(thenClose bool, reply *bool) (err error) {
	if self.vault.IsOpen() {
		*reply = true
		if thenClose {
			err = self.vault.Close()
		}
	} else {
		*reply = false
	}
	return
}

func (self *server) Get(name string, reply *string) (err error) {
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot get %s", name)
	}
	key, err := self.vault.Item(name)
	if err != nil {
		return
	}
	if key == nil || key.IsDeleted() {
		return errors.Newf("Unknown key %s", name)
	}
	*reply = key.Password()
	return
}

func (self *server) List(filter string, reply *[]string) (err error) {
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot list")
	}
	*reply, err = self.vault.List(filter)
	return
}

func (self *server) Open(master string, reply *bool) (err error) {
	if self.vault.IsOpen() {
		return errors.Newf("Vault is already open: cannot open")
	}
	err = self.vault.Open(master, self.config)
	return
}

func (self *server) Merge(args MergeArgs, reply *bool) (err error) {
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot merge")
	}
	vault := newVault(args.Vault)
	err = vault.Open(args.Master, self.config)
	if err != nil {
		return
	}
	if !vault.IsOpen() {
		return errors.Newf("Merge vault is not open: cannot merge")
	}
	err = self.vault.Merge(vault)
	if err != nil {
		return
	}
	err = vault.Close()
	if err != nil {
		return
	}
	*reply = true
	return
}

func (self *server) Save(force bool, reply *bool) (err error) {
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot save")
	}
	err = self.vault.Save(force, self.config)
	if err != nil {
		return
	}
	*reply = true
	return
}

func (self *server) Set(args SetArgs, reply *string) (err error) {
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot set")
	}
	if args.Recipe != "" {
		err = self.vault.SetRandom(args.Key, args.Recipe)
	} else {
		err = self.vault.SetPass(args.Key, args.Pass)
	}
	if err != nil {
		return
	}
	err = self.Get(args.Key, reply)
	return
}

func (self *server) Unset(key string, reply *bool) (err error) {
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot unset")
	}
	err = self.vault.Unset(key)
	*reply = err == nil
	return
}

func (self *server) Stop(status int, reply *bool) (err error) {
	if self.vault.IsOpen() {
		err = self.vault.Close()
		if err != nil {
			return
		}
	}
	err = self.listener.Close()
	if err != nil {
		return errors.Decorated(err)
	}
	*reply = true
	return
}
