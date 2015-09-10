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

type httpChannelServer struct {
	config	 core.Config
	server	 server.Server
	handler	 *blockingHandler
	listener net.Listener
}

type httpChannelClient struct {
	config	 core.Config
	proxy  server.Server
	startFunc server.ProxyStartFunc
	client *rpc.Client
}

var _ server.Server = &httpChannelServer{}
var _ ChannelServer = &httpChannelServer{}
var _ ChannelClient = &httpChannelClient{}

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

func HttpChannelServer(config core.Config, server server.Server) ChannelServer {
	return &httpChannelServer{
		config: config,
		server: server,
		handler: &blockingHandler{
			lock: &sync.RWMutex{},
		},
	}
}

func (self *httpChannelServer) Bind() (err error) {
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

func (self *httpChannelServer) Disconnect() {
	self.handler.lock.Lock() // will never unlock, but the server is dead anyway (this barrier ensures that all connections are served)
}

func (self *httpChannelServer) IsOpen(thenClose bool, reply *bool) error {
	return self.server.IsOpen(thenClose, reply)
}

func (self *httpChannelServer) Get(name string, reply *string) error {
	return self.server.Get(name, reply)
}

func (self *httpChannelServer) List(filter string, reply *[]string) error {
	return self.server.List(filter, reply)
}

func (self *httpChannelServer) Open(master string, reply *bool) error {
	return self.server.Open(master, reply)
}

func (self *httpChannelServer) Merge(args server.MergeArgs, reply *bool) error {
	return self.server.Merge(args, reply)
}

func (self *httpChannelServer) Save(force bool, reply *bool) error {
	return self.server.Save(force, reply)
}

func (self *httpChannelServer) Set(args server.SetArgs, reply *string) error {
	return self.server.Set(args, reply)
}

func (self *httpChannelServer) Unset(key string, reply *bool) error {
	return self.server.Unset(key, reply)
}

func (self *httpChannelServer) Stop(status int, reply *bool) (err error) {
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

func (self *httpChannelServer) Ping(info string, reply *string) error {
	return self.server.Ping(info, reply)
}

func (self *httpChannelServer) SetMaster(master string, reply *bool) error {
	return self.server.SetMaster(master, reply)
}

// ----------------------------------------------------------------

func HttpChannelClient(config core.Config, startFunc server.ProxyStartFunc, proxy server.Server) ChannelClient {
	return &httpChannelClient{
		config: config,
		proxy: proxy,
		startFunc: startFunc,
	}
}

func (self *httpChannelClient) Connect() (err error) {
	host, port := networkConfig(self.config)
	endpoint := fmt.Sprintf("%s:%d", host, port)
	client, err := rpc.DialHTTP("tcp", endpoint)
	if err != nil {
		e := self.startFunc()
		if e != nil {
			err = errors.Decorated(e)
			return
		}
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

func (self *httpChannelClient) Disconnect() {
}

func (self *httpChannelClient) IsOpen(thenClose bool, reply *bool) (err error) {
	err = self.client.Call("Gate.IsOpen", thenClose, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) Get(name string, reply *string) (err error) {
	err = self.client.Call("Gate.Get", name, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) List(filter string, reply *[]string) (err error) {
	err = self.client.Call("Gate.List", filter, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) Open(master string, reply *bool) (err error) {
	err = self.client.Call("Gate.Open", master, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) Merge(args server.MergeArgs, reply *bool) (err error) {
	err = self.client.Call("Gate.Merge", args, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) Save(force bool, reply *bool) (err error) {
	err = self.client.Call("Gate.Save", force, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) Set(args server.SetArgs, reply *string) (err error) {
	err = self.client.Call("Gate.Set", args, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) Unset(key string, reply *bool) (err error) {
	err = self.client.Call("Gate.Unset", key, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) Stop(status int, reply *bool) (err error) {
	err = self.client.Call("Gate.Stop", status, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) SetMaster(master string, reply *bool) (err error) {
	err = self.client.Call("Gate.SetMaster", master, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *httpChannelClient) Ping(info string, reply *string) (err error) {
	err = self.client.Call("Gate.Ping", info, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}
