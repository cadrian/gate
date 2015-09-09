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

package server

// A passwords generator

import (
	"gate/core/errors"
)

import (
	"io"
	"os"
)

// Password generator
type Generator interface {
	// Generate a new password and return it.
	New() (result string, err error)
}

type generator struct {
	recipe []generator_mix
	length int
}

var _ Generator = &generator{}

type generator_mix struct {
	quantity   int
	ingredient string
}

func (self *generator) New() (result string, err error) {
	in, err := os.Open("/dev/random")
	if err != nil {
		return "", errors.Decorated(err)
	}
	defer in.Close()
	return self.generated(in)
}

func (self *generator) generated(in io.Reader) (result string, err error) {
	for _, mix := range self.recipe {
		result, err = mix.extend(in, result)
		if err != nil {
			return
		}
	}
	return
}

func (self generator_mix) extend(in io.Reader, pass string) (result string, err error) {
	result = pass
	for i := 0; i < self.quantity; i++ {
		result, err = self.extend_pass(in, result)
		if err != nil {
			return
		}
	}
	return
}

func (self generator_mix) extend_pass(in io.Reader, pass string) (result string, err error) {
	data := make([]byte, 0, 3)
	n, err := in.Read(data)
	if err != nil {
		return "", errors.Decorated(err)
	}
	if n < 3 {
		return "", errors.New("not enough data")
	}
	b1 := data[0]
	b2 := data[1]
	b := byte(int((b1&0x7f)<<8+b2) % len(self.ingredient))
	i := int(data[2]) % (len(pass) + 1)
	result = string(append(append(append(make([]byte, 0, len(pass)+1), pass[:i]...), b), pass[i:]...))
	return
}

type parse_generator_context struct {
	recipe          []generator_mix
	total_quantity  int
	last_quantity   int
	last_ingredient string
	index           int
	source          string
}

// Return a generator using the given source.
func NewGenerator(source string) (result Generator, err error) {
	context := &parse_generator_context{
		recipe:         make([]generator_mix, 0, 128),
		total_quantity: 0,
		index:          0,
		source:         source,
	}
	err = context.parse_recipe()
	if err == nil {
		result = &generator{
			recipe: context.recipe,
			length: context.total_quantity,
		}
	}
	return
}

func (self *parse_generator_context) parse_recipe() (err error) {
	for err == nil && self.index < len(self.source) {
		err = self.parse_mix()
	}
	return
}

func (self *parse_generator_context) parse_mix() (err error) {
	err = self.parse_quantity()
	if err != nil {
		return
	}
	err = self.parse_ingredient()
	if err != nil {
		return
	}
	if self.last_quantity == 0 {
		self.last_quantity = 1
	}
	self.recipe = append(
		self.recipe,
		generator_mix{
			quantity:   self.last_quantity,
			ingredient: self.last_ingredient,
		},
	)
	return
}

func (self *parse_generator_context) parse_quantity() (err error) {
	self.last_quantity = 0
	for self.index < len(self.source) {
		switch b := self.source[self.index]; {
		case b >= '0' && b <= '9':
			self.last_quantity = self.last_quantity*10 + int(b-'0')
			self.index++
		default:
			break
		}
	}
	return
}

const (
	letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	figures = "0123456789"
	symbols = "(-_)~#{[|^@]}+=<>,?./!§"
)

func (self *parse_generator_context) parse_ingredient() (err error) {
	self.last_ingredient = ""
	for self.index < len(self.source) {
		switch b := self.source[self.index]; b {
		case 'a':
			self.last_ingredient = self.last_ingredient + letters
			self.index++
		case 'n':
			self.last_ingredient = self.last_ingredient + figures
			self.index++
		case 's':
			self.last_ingredient = self.last_ingredient + symbols
			self.index++
		case '+':
			self.index++
			break
		default:
			err = errors.Newf("expected one of 'a', 'n', 's', '+'; not '%v'", rune(b))
			break
		}
	}
	if err == nil && self.last_ingredient == "" {
		if self.last_quantity == 0 {
			err = errors.Newf("expected ingredient or quantity at %d", self.index)
		} else {
			err = errors.Newf("expected ingredient at %d", self.index)
		}
	}
	return
}
