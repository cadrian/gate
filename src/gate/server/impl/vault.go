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

// Keys vault management

import (
	"gate/core"
	"gate/core/errors"
	"gate/core/exec"
)

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"sort"
)

// Return a reader
type In func() (io.ReadCloser, error)

// Return a writer
type Out func() (io.WriteCloser, error)

// The vault interface.
type Vault interface {
	Open(master string, config core.Config) error
	IsOpen() bool
	Close(config core.Config) error
	Item(name string) (Key, error)
	List(filter string) ([]string, error)
	Merge(other Vault) error
	Save(force bool, config core.Config) error
	SetRandom(name string, recipe string) error
	SetPass(name string, pass string) error
	Unset(name string) error
	SetMaster(master string) error
}

type vault struct {
	data	map[string]Key
	dirty	bool
	in	In
	out	Out
	open	bool
	master	string
	recipes map[string]Generator
	decode	func(*vault, io.ReadCloser, chan error)
	newkey	func(string, string) Key
}

var _ Vault = &vault{}

func finalize(v *vault) {
	v.Close(nil)
}

// Create a new vault.
func NewVault(in In, out Out) (result Vault) {
	v := &vault{
		data:	 make(map[string]Key),
		in:	 in,
		out:	 out,
		recipes: make(map[string]Generator, 32),
		decode:	 bf_decode,
		newkey:	 bf_newkey,
	}
	runtime.SetFinalizer(v, finalize)
	result = v
	return
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

	barrier := make(chan error)

	prepare := func(cmd *exec.Cmd) (err error) {
		cmd.Env = append(os.Environ(), fmt.Sprintf("VAULT_MASTER=%s", master))
		cmd.Stdin = instream

		out, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Decorated(err)
		}

		go self.decode(self, out, barrier)

		return
	}

	run := func(cmd *exec.Cmd) (err error) {
		e := <-barrier
		if e != io.EOF {
			err = errors.Decorated(e)
		}

		return
	}

	err = exec.Command(prepare, run, "openssl", cipher, "-d", "-a", "-pass", "env:VAULT_MASTER")
	if err != nil {
		return
	}

	self.master = master
	self.open = true
	return
}

func (self *vault) IsOpen() bool {
	return self.open
}

func (self *vault) Close(config core.Config) (err error) {
	if config != nil {
		err = self.Save(false, config)
		if err != nil {
			return
		}
	}

	self.data = make(map[string]Key)
	self.open = false
	self.master = ""

	runtime.GC()

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
		if k.IsDeleted() {
			log.Printf("   # %s", k.Name())
		} else if re_filter.MatchString(k.Name()) {
			log.Printf("   + %s", k.Name())
			result = append(result, k.Name())
		} else {
			log.Printf("   - %s", k.Name())
		}
	}
	result = result[:len(result)]
	sort.Strings(result)
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

	pipe := make(chan io.WriteCloser, 1)

	prepare := func(cmd *exec.Cmd) (err error) {
		cmd.Env = append(os.Environ(), fmt.Sprintf("VAULT_MASTER=%s", self.master))
		cmd.Stdout = outstream
		p, err := cmd.StdinPipe()
		if err != nil {
			return errors.Decorated(err)
		}
		pipe <- p
		return
	}

	run := func(cmd *exec.Cmd) (err error) {
		p := <-pipe
		for _, k := range self.data {
			code := k.Encoded()
			n, err := p.Write([]byte(code))
			if err != nil {
				return errors.Decorated(err)
			}
			if n < len(code) {
				return errors.Newf("Incomplete write")
			}
		}
		err = p.Close()
		if err != nil {
			return errors.Decorated(err)
		}
		return
	}

	err = exec.Command(prepare, run, "openssl", cipher, "-a", "-pass", "env:VAULT_MASTER")

	return
}

func (self *vault) Save(force bool, config core.Config) (err error) {
	if self.dirty || force {
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
	if ok {
		k.SetPassword(pass)
	} else {
		k = self.newkey(name, pass)
		self.data[name] = k
	}
	self.dirty = true
	return
}

func (self *vault) SetMaster(master string) (err error) {
	if master == "" {
		err = errors.Newf("empty master not allowed")
	} else {
		self.master = master
		self.dirty = true
	}
	return
}
