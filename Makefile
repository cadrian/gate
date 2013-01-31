BUILD=go #gccgo
ECHO=:

all:
	BUILD=$(BUILD) sh build.sh $(ECHO)

clean:
	rm -rf bin pkg

install: all
	sh install.sh

.PHONY: all clean install
