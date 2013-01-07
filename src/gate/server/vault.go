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

type Vault interface {
	Open(master string, config core.Config) (err error)
	Item(name string) (result Key, err error)
}

type vault struct {
	data map[string]*key
	in func() (io.Reader, error)
}

var _ Vault = &vault{}

func NewVault(in func() (io.Reader, error)) (result Vault) {
	result = &vault{
		data: make(map[string]*key),
		in: in,
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

func (self *vault) decode(out io.Reader, errs chan error) {
	buffer := &bytes.Buffer{}
	for done := false; !done; {
		_, err := buffer.ReadFrom(out)
		done = err == io.EOF
	}
	data := string(buffer.Bytes())
	for _, linematch := range decoder.FindAllStringIndex(data, -1) {
		name := decode_group(data, "name", linematch)
		pass := decode_group(data, "pass", linematch)
		delcount, err := decode_group_int(data, "del", linematch)
		if err != nil {
			errs <- err
			continue
		}
		addcount, err := decode_group_int(data, "add", linematch)
		if err != nil {
			errs <- err
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

	errs <- io.EOF
}

func (self *vault) Open(master string, config core.Config) (err error) {
	instream, err := self.in()
	if err != nil {
		return errors.Decorated(err)
	}

	cipher, err := config.Eval("", "vault", "openssl.cipher")
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
	err = cmd.Start()
	if err != nil {
		return errors.Decorated(err)
	}

	errs := make(chan error)
	go self.decode(out, errs)

	for err == nil {
		err = <-errs
	}
	if err != io.EOF {
		return errors.Decorated(err)
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Decorated(err)
	}

	return
}

func (self *vault) Item(name string) (result Key, err error) {
	result, ok := self.data[name]
	if !ok {
		err = errors.Newf("Unknown key: %s", name)
	}
	return
}
