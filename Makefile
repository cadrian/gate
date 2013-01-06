FILES = $(shell find src -name \*.go -print)

all: dep target/server target/console target/menu

dep: target/.dep_flag

target/.dep_flag: target/.flag
	go get github.com/sbinet/liner
	touch target/.dep_flag

target/server: target/.flag Makefile $(FILES)
	go build -o target/server src/server.go

target/console: target/.flag Makefile $(FILES)
	go build -o target/console src/console.go

target/menu: target/.flag Makefile $(FILES)
	go build -o target/menu src/menu.go

target/.flag:
	mkdir -p target
	touch target/.flag

clean:
	rm -rf target

.PHONY: all dep clean
