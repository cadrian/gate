all:
	sh build.sh ':'

clean:
	rm -rf bin pkg

install: all
	sh install.sh

.PHONY: all clean install
