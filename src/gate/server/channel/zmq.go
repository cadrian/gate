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

// 0mq channel implementation - not finished!!!

import (
	"gate/core"
	"gate/core/errors"
	"gate/server"
)

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"os"
	"sync"
	"time"
)

type zmqChannelServer struct {
	config core.Config
	server server.Server
}

type zmqChannelClient struct {
	config core.Config
	proxy  server.Server
}

// ----------------------------------------------------------------

func ZmqChannelServer(config core.Config, server server.Server) ChannelServer {
	return &zmqChannelServer{
		config: config,
		server: server,
	}
}

func (self *zmqChannelServer) Bind() (err error) {
	host, port := networkConfig(self.config)
	return
}

func (self *zmqChannelServer) Disconnect() {
}

func (self *zmqChannelServer) IsOpen(thenClose bool, reply *bool) error {
	return self.server.IsOpen(thenClose, reply)
}

func (self *zmqChannelServer) Get(name string, reply *string) error {
	return self.server.Get(name, reply)
}

func (self *zmqChannelServer) List(filter string, reply *[]string) error {
	return self.server.List(filter, reply)
}

func (self *zmqChannelServer) Open(master string, reply *bool) error {
	return self.server.Open(master, reply)
}

func (self *zmqChannelServer) Merge(args server.MergeArgs, reply *bool) error {
	return self.server.Merge(args, reply)
}

func (self *zmqChannelServer) Save(force bool, reply *bool) error {
	return self.server.Save(force, reply)
}

func (self *zmqChannelServer) Set(args server.SetArgs, reply *string) error {
	return self.server.Set(args, reply)
}

func (self *zmqChannelServer) Unset(key string, reply *bool) error {
	return self.server.Unset(key, reply)
}

func (self *zmqChannelServer) Stop(status int, reply *bool) error {
	return self.server.Stop(status, reply)
}

func (self *zmqChannelServer) Ping(info string, reply *string) error {
	return self.server.Ping(info, reply)
}

func (self *zmqChannelServer) SetMaster(master string, reply *bool) error {
	return self.server.SetMaster(master, reply)
}

// ----------------------------------------------------------------

func ZmqChannelClient(config core.Config, startFunc server.ProxyStartFunc, proxy server.Server) ChannelClient {
	return &zmqChannelClient{
		config: config,
		proxy: proxy,
		startFunc: startFunc,
	}
}

func (self *zmqChannelClient) Connect() (err error) {
	host, port := networkConfig(self.config)
	return
}

func (self *zmqChannelClient) Disconnect() {
}

func (self *zmqChannelClient) IsOpen(thenClose bool, reply *bool) (err error) {
	return
}

func (self *zmqChannelClient) Get(name string, reply *string) (err error) {
	return
}

func (self *zmqChannelClient) List(filter string, reply *[]string) (err error) {
	return
}

func (self *zmqChannelClient) Open(master string, reply *bool) (err error) {
	return
}

func (self *zmqChannelClient) Merge(args server.MergeArgs, reply *bool) (err error) {
	return
}

func (self *zmqChannelClient) Save(force bool, reply *bool) (err error) {
	return
}

func (self *zmqChannelClient) Set(args server.SetArgs, reply *string) (err error) {
	return
}

func (self *zmqChannelClient) Unset(key string, reply *bool) (err error) {
	return
}

func (self *zmqChannelClient) Stop(status int, reply *bool) (err error) {
	return
}

func (self *zmqChannelClient) SetMaster(master string, reply *bool) (err error) {
	return
}

func (self *zmqChannelClient) Ping(info string, reply *string) (err error) {
	return
}
