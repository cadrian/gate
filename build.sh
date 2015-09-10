#!/bin/sh

ECHO=${1:-echo}
BUILD=${BUILD:-go}

export GOPATH=$(dirname $(readlink -f $0))
export PATH=$GOPATH/bin:"$PATH"
export TMPDIR=${TMPDIR:-$GOPATH/.tmp} # because /tmp may have noexec option

mkdir -p $TMPDIR

$ECHO Cleaning up
rm -rf bin pkg

$ECHO Fetching deps
go get code.google.com/p/go.crypto/scrypt
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
go get github.com/pebbe/zmq4
go get github.com/sbinet/liner

$ECHO
$ECHO Generating mocks
find src -name mocks.go -exec rm {} +
TESTS=$(find src/gate -name \*_test.go -exec dirname {} \; | uniq | cut -c5-)

while read pkg itf; do
    mockgen -self_package=$pkg -package=$(basename $pkg) -destination=src/$pkg/mocks.go $pkg $itf
done <<EOF
gate/server Server
gate/client/commands Commander,Command
gate/client/remote Remoter,Remote,Proxy
gate/client/ui UserInteraction
gate/core Config,XdgContext
EOF

$ECHO Launching tests
go test -i $TESTS
go test $TESTS || exit 1

$ECHO
case $BUILD in
    go)
        compile="go install"
        ;;
    gccgo)
        compile="go install -compiler gccgo -gccgoflags '$CPPFLAGS $CFLAGS $LDFLAGS -static-libgcc'"
        ;;
    *)
        echo "Unknown BUILD=$BUILD -- don't build $(basename exe)"
        ;;
esac

$ECHO Building Gate executables using "'$BUILD'"
find src -name main.go -print | while read main_go; do
    main=$(dirname $main_go)
    exe=${main#src/}
    $ECHO " - "$(basename $exe)
    eval "$compile $exe"
done

$ECHO
$ECHO Done
