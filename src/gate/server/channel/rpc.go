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

package channel

// Go's RPC channel implementation

import (
	"gate/core"
	"gate/core/errors"
	"gate/server"
)

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type blockingHandler struct {
	lock *sync.RWMutex
}

func (self *blockingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	http.DefaultServeMux.ServeHTTP(w, r)
}

type rpcChannelServer struct {
	config	 core.Config
	server	 server.Server
	handler	 *blockingHandler
	listener net.Listener
}

type rpcChannelClient struct {
	proxy  server.Server
	host   string
	port   int
	wait   bool
	client *rpc.Client
}

var _ server.Server = &rpcChannelServer{}
var _ ChannelServer = &rpcChannelServer{}
var _ ChannelClient = &rpcChannelClient{}

func networkConfig(config core.Config) (host, port string) {
	var e error
	host, e = config.Eval("", "connection", "host", os.Getenv)
	if e != nil {
		host = "127.0.0.1"
	}
	port, e = config.Eval("", "connection", "port", os.Getenv)
	if e != nil {
		port = "8532"
	}
	return
}

// ----------------------------------------------------------------

func RpcChannelServer(config core.Config, server server.Server) ChannelServer {
	return &rpcChannelServer{
		config: config,
		server: server,
		handler: &blockingHandler{
			lock: &sync.RWMutex{},
		},
	}
}

func (self *rpcChannelServer) Bind() (err error) {
	rpc.RegisterName("Gate", self)
	rpc.HandleHTTP()

	host, port := networkConfig(self.config)
	endpoint := fmt.Sprintf("%s:%s", host, port)
	self.listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		err = errors.Decorated(err)
		return
	}

	go http.Serve(self.listener, self.handler)

	return
}

func (self *rpcChannelServer) Disconnect() {
	self.handler.lock.Lock() // will never unlock, but the server is dead anyway (this barrier ensures that all connections are served)
}

func (self *rpcChannelServer) IsOpen(thenClose bool, reply *bool) error {
	return self.server.IsOpen(thenClose, reply)
}

func (self *rpcChannelServer) Get(name string, reply *string) error {
	return self.server.Get(name, reply)
}

func (self *rpcChannelServer) List(filter string, reply *[]string) error {
	return self.server.List(filter, reply)
}

func (self *rpcChannelServer) Open(master string, reply *bool) error {
	return self.server.Open(master, reply)
}

func (self *rpcChannelServer) Merge(args server.MergeArgs, reply *bool) error {
	return self.server.Merge(args, reply)
}

func (self *rpcChannelServer) Save(force bool, reply *bool) error {
	return self.server.Save(force, reply)
}

func (self *rpcChannelServer) Set(args server.SetArgs, reply *string) error {
	return self.server.Set(args, reply)
}

func (self *rpcChannelServer) Unset(key string, reply *bool) error {
	return self.server.Unset(key, reply)
}

func (self *rpcChannelServer) Stop(status int, reply *bool) (err error) {
	err = self.server.Stop(status, reply)
	if err != nil {
		return
	}
	err = self.listener.Close()
	if err != nil {
		return errors.Decorated(err)
	}
	return
}

func (self *rpcChannelServer) Ping(info string, reply *string) error {
	return self.server.Ping(info, reply)
}

func (self *rpcChannelServer) SetMaster(master string, reply *bool) error {
	return self.server.SetMaster(master, reply)
}

// ----------------------------------------------------------------

func RpcChannelClient(host string, port int, wait bool, proxy server.Server) ChannelClient {
	return &rpcChannelClient{
		proxy: proxy,
		host:  host,
		port:  port,
		wait:  wait,
	}
}

func (self *rpcChannelClient) Connect() (err error) {
	endpoint := fmt.Sprintf("%s:%d", self.host, self.port)
	client, err := rpc.DialHTTP("tcp", endpoint)
	if self.wait {
		for delay := 100 * time.Millisecond; err != nil && delay <= 3*time.Second; delay *= 2 {
			// if the server just started, maybe it needs time to settle
			time.Sleep(delay)
			client, err = rpc.DialHTTP("tcp", endpoint)
		}
	}
	if err != nil {
		err = errors.Decorated(err)
		return
	}

	self.client = client
	return
}

func (self *rpcChannelClient) Disconnect() {
}

func (self *rpcChannelClient) IsOpen(thenClose bool, reply *bool) (err error) {
	err = self.client.Call("Gate.IsOpen", thenClose, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) Get(name string, reply *string) (err error) {
	err = self.client.Call("Gate.Get", name, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) List(filter string, reply *[]string) (err error) {
	err = self.client.Call("Gate.List", filter, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) Open(master string, reply *bool) (err error) {
	err = self.client.Call("Gate.Open", master, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) Merge(args server.MergeArgs, reply *bool) (err error) {
	err = self.client.Call("Gate.Merge", args, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) Save(force bool, reply *bool) (err error) {
	err = self.client.Call("Gate.Save", force, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) Set(args server.SetArgs, reply *string) (err error) {
	err = self.client.Call("Gate.Set", args, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) Unset(key string, reply *bool) (err error) {
	err = self.client.Call("Gate.Unset", key, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) Stop(status int, reply *bool) (err error) {
	err = self.client.Call("Gate.Stop", status, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) SetMaster(master string, reply *bool) (err error) {
	err = self.client.Call("Gate.SetMaster", master, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *rpcChannelClient) Ping(info string, reply *string) (err error) {
	err = self.client.Call("Gate.Ping", info, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}
