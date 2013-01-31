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

package server

// Server-side (i.e. actual) server object

import (
	"gate/core"
	"gate/core/errors"
)

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

// Arguments to the "merge" operation.
type MergeArgs struct {
	Vault string
	Master string
}

// Arguments to the "set" operation.
type SetArgs struct {
	Key string
	Pass string
	Recipe string
}

// The server interface implemented both by the actual (server-side)
// object and the proxy.
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
	Ping(info string, reply *string) error
}

// A server-side server and extra (non-exported) methods
type ServerLocal interface {
	Server() Server
	Wait() (int, error)
}

type blockingHandler struct {
	lock *sync.RWMutex
}

func (self *blockingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	http.DefaultServeMux.ServeHTTP(w, r)
}

type server struct {
	vault Vault
	config core.Config
	listener net.Listener
	handler *blockingHandler
	running bool
	status chan int
}

type serverLocal struct {
	server *server
}

var _ Server = &server{}
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

	srv := &server{
		vault: newVault(vault_path),
		config: config,
		status: make(chan int),
	}
	rpc.RegisterName("Gate", srv)
	rpc.HandleHTTP()

	host, e := config.Eval("", "connection", "host", os.Getenv)
	if e != nil {
		host = "127.0.0.1"
	}
	port, e := config.Eval("", "connection", "port", os.Getenv)
	if e != nil {
		port = "8532"
	}

	endpoint := fmt.Sprintf("%s:%s", host, port)
	srv.listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		err = errors.Decorated(err)
		return
	}
	srv.running = true

	srv.handler = &blockingHandler{
		lock: &sync.RWMutex{},
	}

	go http.Serve(srv.listener, srv.handler)
	result = &serverLocal {
		server: srv,
	}

	log.Printf("Started.")
	return
}

func (self *server) IsOpen(thenClose bool, reply *bool) (err error) {
	log.Printf("IsOpen(thenClose=%t)", thenClose)
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

func (self *server) List(filter string, reply *[]string) (err error) {
	log.Printf("List(filter='%s')", filter)
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot list")
	}
	*reply, err = self.vault.List(filter)
	return
}

func (self *server) Open(master string, reply *bool) (err error) {
	log.Printf("Open(master='***')")
	if self.vault.IsOpen() {
		return errors.Newf("Vault is already open: cannot open")
	}
	err = self.vault.Open(master, self.config)
	*reply = err == nil
	return
}

func (self *server) Merge(args MergeArgs, reply *bool) (err error) {
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

func (self *server) Set(args SetArgs, reply *string) (err error) {
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

func (self *server) Unset(key string, reply *bool) (err error) {
	log.Printf("Unset(key='%s')", key)
	if !self.vault.IsOpen() {
		return errors.Newf("Vault is not open: cannot unset")
	}
	err = self.vault.Unset(key)
	*reply = err == nil
	return
}

func (self *server) Stop(status int, reply *bool) (err error) {
	log.Printf("Stop(status=%d)", status)
	if self.vault.IsOpen() {
		err = self.vault.Save(false, self.config)
		if err != nil {
			return
		}
		err = self.vault.Close()
		if err != nil {
			return
		}
	}
	err = self.listener.Close()
	if err != nil {
		return errors.Decorated(err)
	}
	self.running = false
	self.status <- status
	*reply = true
	return
}

func (self *server) Ping(info string, reply *string) (err error) {
	log.Printf("Ping(info='%s')", info)
	*reply = info
	return
}

func (self *serverLocal) Wait() (result int, err error) {
	if self.server.running {
		result = <-self.server.status
		self.server.handler.lock.Lock() // will never unlock, but the server is dead anyway (this barrier ensures that all connections are served)
	} else {
		err = errors.New("server not running")
	}
	return
}

func (self *serverLocal) Server() Server {
	return self.server
}
