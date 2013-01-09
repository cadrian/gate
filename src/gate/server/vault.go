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
	"gate/core"
	"gate/core/errors"
)

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

type In func() (io.ReadCloser, error)
type Out func() (io.WriteCloser, error)

type Vault interface {
	Open(master string, config core.Config) error
	IsOpen() bool
	Close() error
	Item(name string) (Key, error)
	List(filter string) ([]string, error)
	Merge(other Vault) error
	Save(force bool, config core.Config) error
	SetRandom(name string, recipe string) error
	SetPass(name string, pass string) error
	Unset(name string) error
}

type vault struct {
	data map[string]*key
	dirty bool
	in In
	out Out
	open bool
	master string
	recipes map[string]Generator
}

var _ Vault = &vault{}

func NewVault(in In, out Out) (result Vault) {
	result = &vault{
		data: make(map[string]*key),
		in: in,
		out: out,
		recipes: make(map[string]Generator, 32),
	}
	return
}

var decoder = regexp.MustCompile("^(?P<name>[^:]+):(?P<add>[0-9]+):(?P<del>[0-9]+):(?P<pass>.*)$")

func decode_group(data string, name string, match []int) string {
	return string(decoder.ExpandString(make([]byte, 0, 1024), "$" + name, data, match))
}

func decode_group_int(data string, name string, match []int) (result int64, err error) {
	s := string(decoder.ExpandString(make([]byte, 0, 1024), "$"+name, data, match))
	result, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.Decorated(err)
	}
	return
}

func (self *vault) decode(out io.ReadCloser, barrier chan error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))
	_, err := buffer.ReadFrom(out)
	if err != nil {
		barrier <- err
		return
	}
	data := string(buffer.Bytes())

	for _, linematch := range decoder.FindAllStringIndex(data, -1) {
		name := decode_group(data, "name", linematch)
		pass := decode_group(data, "pass", linematch)
		delcount, err := decode_group_int(data, "del", linematch)
		if err != nil {
			barrier <- err
			continue
		}
		addcount, err := decode_group_int(data, "add", linematch)
		if err != nil {
			barrier <- err
			continue
		}

		k := &key{
			name: name,
			pass: pass,
			delcount: delcount,
			addcount: addcount,
		}
		self.data[name] = k
	}

	barrier <- io.EOF
}

func (self *vault) Open(master string, config core.Config) (err error) {
	instream, err := self.in()
	if err != nil {
		return errors.Decorated(err)
	}
	defer instream.Close()

	cipher, err := config.Eval("", "vault", "openssl.cipher", os.Getenv)
	if err != nil {
		return
	}

	cmd := exec.Command("openssl", cipher, "-d", "-a", "-pass", "env:VAULT_MASTER")
	cmd.Env = append(os.Environ(), fmt.Sprintf("VAULT_MASTER=%s", master))
	cmd.Stdin = instream

	out, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Decorated(err)
	}

	barrier := make(chan error)
	go self.decode(out, barrier)

	err = cmd.Start()
	if err != nil {
		return errors.Decorated(err)
	}

	err = <-barrier
	if err != io.EOF {
		return errors.Decorated(err)
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Decorated(err)
	}

	self.master = master
	self.open = true
	return
}

func (self *vault) IsOpen() bool {
	return self.open
}

func (self *vault) Close() (err error) {
	self.data = make(map[string]*key)
	self.open = false
	self.master = ""
	return
}

func (self *vault) Item(name string) (result Key, err error) {
	result, ok := self.data[name]
	if !ok {
		err = errors.Newf("Unknown key: %s", name)
	}
	return
}

func (self *vault) List(filter string) (result []string, err error) {
	re_filter, err := regexp.Compile(filter)
	if err != nil {
		err = errors.Decorated(err)
		return
	}

	result = make([]string, 0, len(self.data))
	for _, k := range self.data {
		if !k.IsDeleted() && re_filter.MatchString(k.Name()) {
			result = append(result, k.Name())
		}
	}
	result = result[:len(result)]
	return
}

func (self *vault) Merge(o Vault) (err error) {
	other := o.(*vault)
	for keyname, key := range self.data {
		other_key, ok := other.data[keyname]
		if ok {
			key.Merge(other_key)
		}
	}
	for keyname, key := range other.data {
		_, ok := self.data[keyname]
		if !ok {
			self.data[keyname] = key
		}
	}
	self.dirty = true
	return
}

func (self *vault) save(config core.Config) (err error) {
	outstream, err := self.out()
	if err != nil {
		return errors.Decorated(err)
	}
	defer outstream.Close()

	cipher, err := config.Eval("", "vault", "openssl.cipher", os.Getenv)
	if err != nil {
		return err
	}

	cmd := exec.Command("openssl", cipher, "-a", "-pass", "env:VAULT_MASTER")
	cmd.Env = append(os.Environ(), fmt.Sprintf("VAULT_MASTER=%s", self.master))
	cmd.Stdout = outstream
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return errors.Decorated(err)
	}

	for _, k := range self.data {
		code := k.Encoded()
		n, err := pipe.Write([]byte(code))
		if err != nil {
			return errors.Decorated(err)
		}
		if n < len(code) {
			return errors.Newf("Incomplete write")
		}
	}

	return
}

func (self *vault) Save(force bool, config core.Config) (err error) {
	if force || self.dirty {
		err = self.save(config)
		if err != nil {
			return
		}
		self.dirty = false
	}
	return
}

func (self *vault) generatePass(recipe string) (result string, err error) {
	gen, ok := self.recipes[recipe]
	if !ok {
		gen, err = NewGenerator(recipe)
		if err != nil {
			return
		}
		self.recipes[recipe] = gen
	}
	return gen.New()
}

func (self *vault) SetRandom(name string, recipe string) (err error) {
	pass, err := self.generatePass(recipe)
	if err != nil {
		return
	}
	return self.SetPass(name, pass)
}

func (self *vault) Unset(name string) (err error) {
	delete(self.data, name)
	return
}

func (self *vault) SetPass(name string, pass string) (err error) {
	k, ok := self.data[name]
	if !ok {
		k = &key{}
		self.data[name] = k
	}
	k.SetPassword(pass)
	self.dirty = true
	return
}
