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

// Client-side access to the server (implements all the RPC operations).

import (
	"gate/core/errors"
)

import (
	"fmt"
	"net/rpc"
	"time"
)

type proxy struct {
	client *rpc.Client
}

var _ Server = &proxy{}

// Return a new proxy to the Gate server identified by the host name and port.
func Proxy(host string, port int, wait bool) (result Server, err error) {
	var client *rpc.Client

	endpoint := fmt.Sprintf("%s:%d", host, port)
	client, err = rpc.DialHTTP("tcp", endpoint)
	if wait {
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

	result = &proxy{
		client: client,
	}
	return
}

func (self *proxy) IsOpen(thenClose bool, reply *bool) (err error) {
	err = self.client.Call("Gate.IsOpen", thenClose, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Get(name string, reply *string) (err error) {
	err = self.client.Call("Gate.Get", name, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) List(filter string, reply *[]string) (err error) {
	err = self.client.Call("Gate.List", filter, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Open(master string, reply *bool) (err error) {
	err = self.client.Call("Gate.Open", master, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Merge(args MergeArgs, reply *bool) (err error) {
	err = self.client.Call("Gate.Merge", args, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Save(force bool, reply *bool) (err error) {
	err = self.client.Call("Gate.Save", force, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Set(args SetArgs, reply *string) (err error) {
	err = self.client.Call("Gate.Set", args, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Unset(key string, reply *bool) (err error) {
	err = self.client.Call("Gate.Unset", key, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Stop(status int, reply *bool) (err error) {
	err = self.client.Call("Gate.Stop", status, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) SetMaster(master string, reply *bool) (err error) {
	err = self.client.Call("Gate.SetMaster", master, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Ping(info string, reply *string) (err error) {
	err = self.client.Call("Gate.Ping", info, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}
