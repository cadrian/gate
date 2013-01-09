FILES = $(shell find src -name \*.go -print)
TESTS = $(shell find src/gate -name \*_test.go -exec dirname {} \; | uniq | cut -c5-)

all: dep target/server target/console target/menu
	go test -i $(TESTS)
	go test $(TESTS)

dep: target/.dep_flag

target/.dep_flag: target/.flag
	go get github.com/sbinet/liner
	go get code.google.com/p/gomock/gomock
	go get code.google.com/p/gomock/mockgen
	touch target/.dep_flag

target/server: target/.flag Makefile $(FILES)
	go build -o target/server src/server.go

target/console: target/.flag Makefile $(FILES)
	go build -o target/console src/console.go

target/menu: target/.flag Makefile $(FILES)
	go build -o target/menu src/menu.go

target/.flag: Makefile
	mkdir -p target
	touch target/.flag

clean:
	rm -rf bin pkg target

.PHONY: all dep clean
.SILENT:
