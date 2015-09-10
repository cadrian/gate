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

// Client-side access to the server (implements all the RPC operations).

import (
	"gate/core"
	"gate/server"
	"gate/server/channel"
)

type proxy struct {
	channel channel.ChannelClient
}

var _ server.Server = &proxy{}

type ProxyStartFunc func() error

// Return a new proxy to the Gate server identified by the host name and port.
func Proxy(config core.Config, startFunc ProxyStartFunc) (result server.Server, err error) {
	p := &proxy{}
	p.channel = channel.HttpChannelClient(config, startFunc, p)
	err = p.channel.Connect()
	if err == nil {
		result = p
	}
	return
}

func (self *proxy) IsOpen(thenClose bool, reply *bool) error {
	return self.channel.IsOpen(thenClose, reply)
}

func (self *proxy) Get(name string, reply *string) error {
	return self.channel.Get(name, reply)
}

func (self *proxy) List(filter string, reply *[]string) error {
	return self.channel.List(filter, reply)
}

func (self *proxy) Open(master string, reply *bool) error {
	return self.channel.Open(master, reply)
}

func (self *proxy) Merge(args server.MergeArgs, reply *bool) error {
	return self.channel.Merge(args, reply)
}

func (self *proxy) Save(force bool, reply *bool) error {
	return self.channel.Save(force, reply)
}

func (self *proxy) Set(args server.SetArgs, reply *string) error {
	return self.channel.Set(args, reply)
}

func (self *proxy) Unset(key string, reply *bool) error {
	return self.channel.Unset(key, reply)
}

func (self *proxy) Stop(status int, reply *bool) (err error) {
	err = self.channel.Stop(status, reply)
	self.channel.Disconnect()
	return
}

func (self *proxy) Ping(info string, reply *string) error {
	return self.channel.Ping(info, reply)
}

func (self *proxy) SetMaster(master string, reply *bool) error {
	return self.channel.SetMaster(master, reply)
}
