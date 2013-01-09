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
package rc

import (
	"gate/core/errors"
)

import (
	"bytes"
	"fmt"
	"io"
)

type FileContent struct {
	data []rune
	index int
}

func (self *FileContent) IsValid() bool {
	return self.index >= 0 && self.index < len(self.data)
}

func (self *FileContent) Current() (result rune, err error) {
	if !self.IsValid() {
		return 0, errors.Newf("invalid current character at index %d", self.index)
	}
	result = self.data[self.index]
	return
}

func (self *FileContent) Next() error {
	if self.index > len(self.data) {
		return errors.Newf("cannot go next, index out of range: %d > %d", self.index, len(self.data))
	}
	self.index++
	return nil
}

func (self *FileContent) Back() error {
	if self.index < 0 {
		return errors.Newf("cannot go back, index out of range: %d < 0", self.index)
	}
	self.index--
	return nil
}

func (self *FileContent) SkipUntil(stop func(rune, int) bool) (result string, err error) {
	buffer := &bytes.Buffer{}
	var k rune
	i := 0
	for done := !self.IsValid(); !done; {
		k, err = self.Current()
		if err != nil {
			return
		}
		if stop(k, i) {
			done = true
		} else {
			buffer.WriteRune(k)
			self.Next()
			i++
			done = !self.IsValid()
		}
	}
	result = buffer.String()
	return
}

func (self *FileContent) SkipWord() (result string, err error) {
	return self.SkipUntil(func(k rune, index int) bool {
		switch {
		case index > 0 && k >= '0' && k <= '9',
			k >= 'A' && k <= 'Z',
			k >= 'a' && k <= 'z',
			k == '_',
			k == '.':
			return false
		}
		return true
	})
}

func (self *FileContent) SkipSymbol(symbol string) (result string, err error) {
	for _, c := range symbol {
		k, err := self.Current()
		if err == nil {
			if k != c {
				return "", errors.Newf("Expected symbol: '%s'", symbol)
			}
		} else {
			return "", err
		}
		self.Next()
	}
	result = symbol
	return
}

func (self *FileContent) SkipBlanks() (result string, err error) {
	return self.SkipUntil(func(k rune, index int) bool {
		switch k {
		case ' ', '\t', '\n', '\r':
			return false
		}
		return true
	})
}

func (self *FileContent) SkipToEndOfLine() (result string, err error) {
	return self.SkipUntil(func(k rune, index int) bool {
		switch k {
		case '\n', '\r':
			return true
		}
		return false
	})
}

func (self *FileContent) Debug() (result string) {
	return fmt.Sprintf("File content: len=%d, index=%d, is valid: %t", len(self.data), self.index, self.IsValid())
}

func readFile(in io.Reader) (result *FileContent) {
	data := &bytes.Buffer{}
	data.ReadFrom(in)
	result = &FileContent {
		data: []rune(string(data.Bytes())),
		index: 0,
	}
	return
}
