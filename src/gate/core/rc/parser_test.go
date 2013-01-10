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

package rc

import (
	"strings"
	"testing"
)

func TestReadAnonymous(t *testing.T) {
	in := strings.NewReader("test = foobar")
	file, err := Read(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Anonymous.Resources) != 1 {
		t.Fatalf("bad anonymous length: %d\n", len(file.Anonymous.Resources))
	}
	if file.Anonymous.Resources["test"] != "foobar" {
		t.Fatalf("missing or wrong test key: %s", file.Anonymous.Resources["test"])
	}
}

func TestReadNamed(t *testing.T) {
	in := strings.NewReader("[what]\ntest = foobar")
	file, err := Read(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Anonymous.Resources) != 0 {
		t.Fatalf("bad anonymous length: %d\n", len(file.Anonymous.Resources))
	}
	if len(file.Sections) != 1 {
		t.Fatalf("bad sections length: %d\n", len(file.Anonymous.Resources))
	}
	section := file.Sections["what"]
	if section == nil {
		t.Fatalf("missing section what")
	}
	if len(section.Resources) != 1 {
		t.Fatalf("bad section length: %d\n", len(section.Resources))
	}
	if section.Resources["test"] != "foobar" {
		t.Fatalf("missing or wrong test key: %s", section.Resources["test"])
	}
}

func TestReadAnonymousAndNamed(t *testing.T) {
	in := strings.NewReader("titi = toto\n[s1] ignored\ntest = foobar\nfoo	=\tbar\n\n[s2]\nwhatever=    nothing\n")
	file, err := Read(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Anonymous.Resources) != 1 {
		t.Fatalf("bad anonymous length: %d\n", len(file.Anonymous.Resources))
	}
	if file.Anonymous.Resources["titi"] != "toto" {
		t.Fatalf("missing or wrong titi anonymous key: %s", file.Anonymous.Resources["titi"])
	}
	if len(file.Sections) != 2 {
		t.Fatalf("bad sections length: %d\n", len(file.Anonymous.Resources))
	}
	s1 := file.Sections["s1"]
	if s1 == nil {
		t.Fatalf("missing section s1")
	}
	if len(s1.Resources) != 2 {
		t.Fatalf("bad s1 length: %d\n", len(s1.Resources))
	}
	if s1.Resources["test"] != "foobar" {
		t.Fatalf("missing or wrong test s1 key: %s", s1.Resources["test"])
	}
	if s1.Resources["foo"] != "bar" {
		t.Fatalf("missing or wrong foo s1 key: %s", s1.Resources["foo"])
	}
	s2 := file.Sections["s2"]
	if s2 == nil {
		t.Fatalf("missing section s2")
	}
	if len(s2.Resources) != 1 {
		t.Fatalf("bad s2 length: %d\n", len(s2.Resources))
	}
	if s2.Resources["whatever"] != "nothing" {
		t.Fatalf("missing or wrong whatever s2 key: %s", s2.Resources["whatever"])
	}
}
