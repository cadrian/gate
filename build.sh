#!/bin/bash

export GOPATH=$(dirname $(readlink -f $0))
export PATH=$GOPATH/bin:"$PATH"

echo Fetching deps
go get github.com/sbinet/liner
go get code.google.com/p/gomock/gomock
go get code.google.com/p/gomock/mockgen

echo Testing
find src -name mocks.go -exec rm {} +
TESTS=$(find src/gate -name \*_test.go -exec dirname {} \; | uniq | cut -c5-)

while read pkg itf; do
    mockgen -self_package=$pkg -package=$(basename $pkg) -destination=src/$pkg/mocks.go $pkg $itf
done <<EOF
gate/server Server
gate/client/commands Commander,Cmd
gate/client/ui UserInteraction
gate/core Config
EOF

go test -i $TESTS
go test $TESTS || exit 1

find src -name main.go -print | while read main_go; do
    main=$(dirname $main_go)
    exe=${main#src/}
    echo Building $(basename $exe)
    go install $exe
done

echo Done
