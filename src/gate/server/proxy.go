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
	"gate/core/errors"
)

import (
	"fmt"
	"net/rpc"
)

type proxy struct {
	client *rpc.Client
}

var _ Server = &proxy{}

func Proxy(serverAddress string, port int) (result Server, err error) {
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", serverAddress, port))
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
	err = self.client.Call("Server.IsOpen", thenClose, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Get(name string, reply *string) (err error) {
	err = self.client.Call("Server.Get", name, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) List(filter string, reply *[]string) (err error) {
	err = self.client.Call("Server.List", filter, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Open(master string, reply *bool) (err error) {
	err = self.client.Call("Server.Open", master, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}

func (self *proxy) Merge(args MergeArgs, reply *bool) (err error) {
	err = self.client.Call("Server.Merge", args, reply)
	if err != nil {
		err = errors.Decorated(err)
	}
	return
}
