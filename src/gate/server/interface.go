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

// Head package for the server definition
package server

// Arguments to the "merge" operation.
type MergeArgs struct {
	Vault  string
	Master string
}

// Arguments to the "set" operation.
type SetArgs struct {
	Key    string
	Pass   string
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
	SetMaster(master string, reply *bool) error
}
