package core

import (
	"testing"
)

func mock_getenv(name string) string {
	switch name {
	case "FOO":
		return "bar"
	}
	return ""
}

func checkeval(expected string, toeval string, t *testing.T) {
	actual := eval(toeval, mock_getenv)
	if expected != actual {
		t.Errorf("[%s] => %s != %s", toeval, expected, actual)
	}
}

func TestExpand(t *testing.T) {
	checkeval("bar", "$FOO", t)
	checkeval("bar", "${FOO}", t)
	checkeval("$FOO", "\\$FOO", t)
	checkeval("\"bar\"", "\"$FOO\"", t)
	checkeval("\"'bar'\"", "\"'$FOO'\"", t)
	checkeval("'$FOO'", "'$FOO'", t)
	checkeval("'\"$FOO\"'", "'\"$FOO\"'", t)
	checkeval("bar'$FOO'", "$FOO'$FOO'", t)
	checkeval("bar # anything $FOO ", "$FOO # anything $FOO ", t)
}
