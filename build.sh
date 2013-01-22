#!/bin/sh

ECHO=${1:-echo}

export GOPATH=$(dirname $(readlink -f $0))
export PATH=$GOPATH/bin:"$PATH"
export TMPDIR=${TMPDIR:-$GOPATH/.tmp} # because /tmp may have noexec option

mkdir -p $TMPDIR

$ECHO Fetching deps
go get github.com/sbinet/liner
go get code.google.com/p/gomock/gomock
go get code.google.com/p/gomock/mockgen
go get code.google.com/p/go.crypto/scrypt

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
gate/core Config
EOF

$ECHO Launching tests
go test -i $TESTS
go test $TESTS || exit 1

$ECHO
find src -name main.go -print | while read main_go; do
    main=$(dirname $main_go)
    exe=${main#src/}
    $ECHO Building $(basename $exe)
    go install $exe
done

$ECHO
$ECHO Done
