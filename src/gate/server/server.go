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

type Server interface {
	Open(master string, reply *bool) error
	IsOpen(thenClose bool, reply *bool) error
	Get(key string, reply *string) error
	List(filter string, reply *[]string) error
}

type server struct {
	vault Vault
	config core.Config
}

var _ Server = &server{}

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
	in := func() (result io.ReadCloser, err error) {
		return os.Open(vault_path)
	}

	result = &server{
		vault: NewVault(in),
		config: config,
	}
	rpc.Register(result)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		err = errors.Decorated(err)
		return
	}
	go http.Serve(listener, nil)
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
		return errors.Newf("Vault is not open: cannot retrieve %s", name)
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
