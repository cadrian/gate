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

// Configuration management

import (
	"gate/core/errors"
	"gate/core/rc"
)

import (
	"fmt"
	"os"
	"strings"
)

// An object mirroring the user configuration
type Config interface {
	// Get a configuration value; may be transformed by the evaluator (e.g. os.Getenv).
	Eval(file string, section string, key string, evaluator func(string) string) (string, error)

	// Get the list of extra configuration files, directly useable as first argument of Eval()
	// "config.rc" itself is ommitted
	ListConfigFiles() ([]string, error)

	VaultPath() (string, error)
}

type config struct {
	files map[string]*rc.File
	vault string
	main_rc string
}

// Get the user configuration
func NewConfig() (result Config, err error) {
	main_rc := "config.rc"
	if len(os.Args) > 1 {
		rc := os.Args[1]
		if strings.HasSuffix(rc, ".rc") {
			_, err = os.Stat(rc)
			if err == nil {
				main_rc = rc
			}
		}
	}

	result = &config{
		files: make(map[string]*rc.File),
		main_rc: main_rc,
	}
	return
}

func (self *config) ListConfigFiles() (result []string, err error) {
	xdg, err := Xdg()
	if err != nil {
		return
	}
	config_path, err := xdg.ConfigHome()
	if err != nil {
		return
	}
	config, err := os.Open(config_path)
	if err != nil {
		err = errors.Decorated(err)
		return
	}
	names, err := config.Readdirnames(0)
	if err != nil {
		err = errors.Decorated(err)
		return
	}
	result = make([]string, 0, len(names) - 1)
	for _, name := range names {
		if name != "config.rc" && strings.HasSuffix(name, ".rc") {
			result = append(result, name)
		}
	}
	return
}

func (self *config) findFile(file string) (result *rc.File, err error) {
	result, ok := self.files[file]
	if ok {
		return
	}

	xdg, err := Xdg()
	if err != nil {
		return
	}
	in, err := xdg.ReadConfig(file)
	if err != nil {
		return
	}

	result, err = rc.Read(in)
	if err != nil {
		in.Close()
		return
	}
	err = in.Close()
	if err != nil {
		return
	}

	self.files[file] = result
	return
}

func (self *config) rawValue(file string, section string, key string) (result string, err error) {
	f, err := self.findFile(file)
	if err != nil {
		return
	}

	var ok bool
	var sec *rc.Section
	if section == "" {
		sec = f.Anonymous
		if sec == nil {
			return "", errors.Newf("No anonymous section in file %s", section, file)
		}
	} else {
		sec, ok = f.Sections[section]
		if !ok {
			return "", errors.Newf("Unknown section [%s] in file %s", section, file)
		}
	}
	result, ok = sec.Resources[key]
	if !ok {
		return "", errors.Newf("Unknown key %s in section [%s] of file %s", key, section, file)
	}
	return
}

type eval_context struct {
	out []rune
	varname []rune
	state int
	in_string bool
}

func (self *eval_context) append(b... rune) {
	self.out = append(self.out, b...)
}

func (self *eval_context) eval(env func(string) string, pb *rune) (next bool) {
	switch self.state {
	case 0: // nominal
		if pb != nil {
			switch b := *pb; b {
			case '#':
				self.append(b)
				self.state = 1
			case '\\':
				self.state = 2
			case '$':
				self.varname = make([]rune, 0, 128)
				self.state = 10
			case '\'':
				self.append(b)
				if !self.in_string {
					self.state = 20
				}
			case '"':
				self.append(b)
				self.in_string = !self.in_string
			default:
				self.append(b)
			}
		}
		next = true
	case 1: // comment
		if pb != nil {
			self.append(*pb)
		}
		next = true
	case 2: // escaped
		if pb != nil {
			switch b := *pb; b {
			case 't':
				self.append('\t')
			case 'n':
				self.append('\n')
			case 'r':
				self.append('\r')
			case '\\':
				self.append('\\')
			case '\'':
				self.append('\'')
			case '$':
				self.append('$')
			default: // unknown escape sequence
				self.append('\\', b)
			}
		}
		self.state = 0
		next = true
	case 10: // after '$'
		if pb == nil {
			self.state = 0
			next = true
		} else {
			switch b := *pb; {
			case b == '{':
				self.state = 12
				next = true
			case (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b == '_'):
				self.varname = append(self.varname, b)
				self.state = 11
				next = true
			default:
				exp := []rune(env(string(self.varname)))
				self.append(exp...)
				self.state = 0
			}
		}
	case 11: // after '$' and a letter
		if pb == nil {
			exp := []rune(env(string(self.varname)))
			self.append(exp...)
			self.state = 0
			next = true
		} else {
			switch b := *pb; {
			case (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b == '_'):
				self.varname = append(self.varname, b)
				next = true
			default:
				exp := []rune(env(string(self.varname)))
				self.append(exp...)
				self.state = 0
			}
		}
	case 12: // after '${'
		if pb == nil {
			exp := []rune(env(string(self.varname)))
			self.append(exp...)
			self.state = 0
			next = true
		} else {
			switch b := *pb; b {
			case '}':
				exp := []rune(env(string(self.varname)))
				self.append(exp...)
				self.state = 0
				next = true
			default:
				self.varname = append(self.varname, b)
				next = true
			}
		}
	case 20: // inside simple-quoted string
		if pb != nil {
			switch b := *pb; b {
			case '\'':
				self.append(b)
				self.state = 0
			default:
				self.append(b)
			}
		}
		next = true
	}
	return
}

func eval(raw string, evaluator func(string) string) (result string) {
	ctx := &eval_context{
		out : make([]rune, 0, 2 * len(raw)),
	}
	for _, b := range raw {
		for next := false; !next; {
			next = ctx.eval(evaluator, &b)
		}
	}
	for next := false; !next; {
		next = ctx.eval(evaluator, nil)
	}

	result = string(ctx.out)
	return
}

func (self *config) Eval(file string, section string, key string, evaluator func(string) string) (result string, err error) {
	if file == "" {
		file = self.main_rc
	}
	result, err = self.rawValue(file, section, key)
	if err != nil {
		return
	}
	if evaluator != nil {
		result = eval(result, evaluator)
	}
	return
}

func (self *config) getVaultPath() (result string, err error) {
	xdg, err := Xdg()
	if err != nil {
		return
	}
	data_home, err := xdg.DataHome()
	if err != nil {
		return
	}
	result = fmt.Sprintf("%s/vault", data_home)
	return
}

func (self *config) VaultPath() (result string, err error) {
	result = self.vault
	if result == "" {
		result, err = self.getVaultPath()
		if err == nil {
			self.vault = result
		}
	}
	return
}
