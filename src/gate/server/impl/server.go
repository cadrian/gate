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

package impl

// Server-side (i.e. actual) server object

import (
	"gate/core"
	"gate/core/errors"
	"gate/server"
	"gate/server/channel"
)

import (
	"io"
	"log"
	"os"
)

// A server-side server and extra (non-exported) methods
type ServerLocal interface {
	Server() server.Server
	Wait() (int, error)
}

type serverImpl struct {
	vault	Vault
	config	core.Config
	channel channel.ChannelServer
	running bool
	status	chan int
}

type serverLocal struct {
	server *serverImpl
}

var _ server.Server = &serverImpl{}
var _ ServerLocal = &serverLocal{}

func newVault(file string) Vault {
	in := func() (result io.ReadCloser, err error) {
		return os.Open(file)
	}
	out := func() (result io.WriteCloser, err error) {
		return os.Create(file)
	}
	return NewVault(in, out)
}

// Start a server on localhost, listening on the given port
func Start(config core.Config) (result ServerLocal, err error) {
	log.Printf("Starting...")

	vault_path, err := config.VaultPath()
	if err != nil {
		return
	}

	srv := &serverImpl{
		vault:	 newVault(vault_path),
		config:	 config,
		status:	 make(chan int),
		running: true,
	}
	srv.channel = channel.RpcChannelServer(config, srv)

	err = srv.channel.Bind()
	if err != nil {
		return
	}

	result = &serverLocal{
		server: srv,
	}

	log.Printf("Started.")
	return
}

func (self *serverImpl) IsOpen(thenClose bool, reply *bool) (err error) {
	log.Printf("IsOpen(thenClose=%t)", thenClose)
	if self.vault.IsOpen() {
		*reply = true
		if thenClose {
			err = self.vault.Close(self.config)
		}
	} else {
		*reply = false
	}
	return
}

func (self *serverImpl) Get(name string, reply *string) (err error) {
	log.Printf("Get(name='%s')", name)
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

func (self *serverImpl) List(filter string, reply *[]string) (err error) {
	log.Printf("List(filter='%s')", filter)
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot list")
	}
	*reply, err = self.vault.List(filter)
	return
}

func (self *serverImpl) Open(master string, reply *bool) (err error) {
	log.Printf("Open(master='***')")
	if self.vault.IsOpen() {
		return errors.Newf("Vault is already open: cannot open")
	}
	err = self.vault.Open(master, self.config)
	*reply = err == nil
	return
}

func (self *serverImpl) Merge(args server.MergeArgs, reply *bool) (err error) {
	log.Printf("Merge(vault='%s', master='***')", args.Vault)
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
		vault.Close(self.config)
		return
	}
	err = vault.Close(self.config)
	if err != nil {
		return
	}
	*reply = true
	return
}

func (self *serverImpl) Save(force bool, reply *bool) (err error) {
	log.Printf("Save(force=%t)", force)
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

func (self *serverImpl) Set(args server.SetArgs, reply *string) (err error) {
	log.Printf("Set(key='%s', ...)", args.Key)
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

func (self *serverImpl) Unset(key string, reply *bool) (err error) {
	log.Printf("Unset(key='%s')", key)
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot unset")
	}
	err = self.vault.Unset(key)
	*reply = err == nil
	return
}

func (self *serverImpl) Stop(status int, reply *bool) (err error) {
	log.Printf("Stop(status=%d)", status)
	if self.vault.IsOpen() {
		err = self.vault.Close(self.config)
		if err != nil {
			return
		}
	}
	self.running = false
	self.status <- status
	*reply = true
	return
}

func (self *serverImpl) Ping(info string, reply *string) (err error) {
	log.Printf("Ping(info='%s')", info)
	*reply = info
	return
}

func (self *serverImpl) SetMaster(master string, reply *bool) (err error) {
	log.Printf("SetMaster(master='***')")
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot set master")
	}
	err = self.vault.SetMaster(master)
	if err == nil {
		*reply = true
	}
	return
}

func (self *serverLocal) Wait() (result int, err error) {
	if self.server.running {
		result = <-self.server.status
		self.server.channel.Disconnect()
	} else {
		err = errors.New("server not running")
	}
	return
}

func (self *serverLocal) Server() server.Server {
	return self.server
}
