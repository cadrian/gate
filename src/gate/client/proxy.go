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
package client

import (
	"gate/core"
	"gate/core/errors"
	"gate/server"
)

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var _proxy server.Server

func dirname() (result string) {
	path := strings.Split(os.Args[0], "/")
	if len(path) == 1 {
		result = "."
	} else {
		result = strings.Join(path[:len(path)-1], "/")
	}
	return
}

func startServer() (err error) {
	cmd := exec.Command("at", "now")
	in, err := cmd.StdinPipe()
	if err != nil {
		return errors.Decorated(err)
	}
	err = cmd.Start()
	if err != nil {
		return errors.Decorated(err)
	}

	in.Write([]byte(fmt.Sprintf("%s/server\n", dirname())))

	err = in.Close()
	if err != nil {
		return errors.Decorated(err)
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Decorated(err)
	}

	// let the server get cozy
	time.Sleep(200 * time.Millisecond)

	return
}

func proxy(config core.Config) (result server.Server, err error) {
	result = _proxy
	if result == nil {
		var (
			host, p string
			port int64
			s server.Server
		)
		host, err = config.Eval("", "connection", "host")
		if err != nil {
			return
		}
		p, err = config.Eval("", "connection", "port")
		if err != nil {
			return
		}
		port, err = strconv.ParseInt(p, 10, 32)
		if err != nil {
			err = errors.Decorated(err)
			return
		}
		s, err = server.Proxy(host, int(port))
		if err != nil {
			err = startServer()
			if err != nil {
				return
			}
			s, err = server.Proxy(host, int(port))
			if err != nil {
				return
			}
		}
		result = s
		_proxy = result
	}
	return
}
