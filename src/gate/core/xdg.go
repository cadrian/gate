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

package core

import (
	"gate/core/errors"
)

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type XdgContext interface {
	ReadData(file string) (io.ReadCloser, error)
	ReadConfig(file string) (io.ReadCloser, error)
	CacheHome() (string, error)
	RuntimeDir() (string, error)
	DataHome() (string, error)
	ConfigHome() (string, error)
}

type xdgContext struct {
	data_dirs []string
	config_dirs []string
	cache_home string
	runtime_dir string
	data_home string
	config_home string
}

var xdg *xdgContext

func getenv(name string, ext func(string) string, def func() string) (result string) {
	result = os.Getenv(name)
	if result == "" && def != nil {
		result = def()
	}
	if result != "" && ext != nil {
		result = ext(result)
	}
	return
}

func getdirs(env string, home string, dflt string) []string {
	def := func () string {
		return dflt
	}
	dirs := strings.Split(getenv(env, nil, def), ":")
	result := make([]string, 0, len(dirs) + 1)
	result = append(result, home)
	result = append(result, dirs...)
	return result
}

func Xdg() (result XdgContext, err error) {
	if xdg == nil {
		var (
			data string
			config string
		)
		xdg = &xdgContext{
		}
		data, err = xdg.dataHome()
		if err != nil {
			return
		}
		xdg.data_dirs = getdirs("XDG_DATA_DIRS", data, "/usr/local/share/:/usr/share/")
		config, err = xdg.configHome()
		if err != nil {
			return
		}
		xdg.config_dirs = getdirs("XDG_CONFIG_DIRS", config, "/usr/local/etc:/etc/xdg")
	}
	result = xdg
	return
}

func checkdir(dirname string) (result string, err error) {
	result = dirname
	info, err := os.Stat(dirname)
	if err == nil {
		if !info.IsDir() {
			err = errors.Newf("%s is not a directory", dirname)
		}
	} else {
		err = os.MkdirAll(dirname, os.ModeDir | 0700)
		if err != nil {
			err = errors.Decorated(err)
		}
	}
	return
}

func read(file string, dirs []string) (result io.ReadCloser, err error) {
	for _, dir := range dirs {
		path := fmt.Sprintf("%s/gate/%s", dir, file)
		result, err = os.Open(path)
		if err == nil {
			return
		}
	}
	err = errors.Newf("Could not find file %s (looked in %s/gate)", file, strings.Join(dirs, "/gate, "))
	return
}

func (self *xdgContext) ReadData(file string) (result io.ReadCloser, err error) {
	return read(file, self.data_dirs)
}

func (self *xdgContext) ReadConfig(file string) (result io.ReadCloser, err error) {
	return read(file, self.config_dirs)
}

func (self *xdgContext) CacheHome() (string, error) {
	if self.cache_home == "" {
		def := func () string {
			return fmt.Sprintf("%s/.cache/gate", getenv("HOME", nil, nil))
		}
		self.cache_home = getenv("XDG_CACHE_HOME", nil, def)
	}
	return checkdir(self.cache_home)
}

func (self *xdgContext) RuntimeDir() (string, error) {
	if self.runtime_dir == "" {
		tmpext := func (tmp string) string {
			return fmt.Sprintf("%s/gate", tmp)
		}
		tmpdef := func () string {
			return fmt.Sprintf("/tmp/gate-%s", getenv("USER", nil, nil))
		}
		xdgdef := func () string {
			return getenv("TMPDIR", tmpext, tmpdef)
		}
		self.runtime_dir = getenv("XDG_RUNTIME_DIR", nil, xdgdef)
	}
	return checkdir(self.runtime_dir)
}

func (self *xdgContext) dataHome() (string, error) {
	if self.data_home == "" {
		def := func () string {
			return fmt.Sprintf("%s/.local/share", getenv("HOME", nil, nil))
		}
		self.data_home = getenv("XDG_DATA_HOME", nil, def)
	}
	return checkdir(self.data_home)
}

func (self *xdgContext) DataHome() (result string, err error) {
	result, err = self.dataHome()
	if err != nil {
		return
	}
	return fmt.Sprintf("%s/gate", result), nil
}

func (self *xdgContext) configHome() (string, error) {
	if self.config_home == "" {
		def := func () string {
			return fmt.Sprintf("%s/.config", getenv("HOME", nil, nil))
		}
		self.config_home = getenv("XDG_CONFIG_HOME", nil, def)
	}
	return checkdir(self.config_home)
}

func (self *xdgContext) ConfigHome() (result string, err error) {
	result, err = self.configHome()
	if err != nil {
		return
	}
	return fmt.Sprintf("%s/gate", result), nil
}
