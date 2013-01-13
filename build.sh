#!/bin/bash

export GOPATH=$(dirname $(readlink -f $0))
export PATH=$GOPATH/bin:"$PATH"

echo Fetching deps
go get github.com/sbinet/liner
go get code.google.com/p/gomock/gomock
go get code.google.com/p/gomock/mockgen

echo Testing
TESTS=$(find src/gate -name \*_test.go -exec dirname {} \; | uniq | cut -c5-)

mkdir -p src/gate/mock_server
mockgen gate/server Server > src/gate/mock_server/server.go

go test -i $TESTS
go test $TESTS || exit 1

find src -name main.go -print | while read main_go; do
    main=$(dirname $main_go)
    exe=${main#src/}
    echo Building $(basename $exe)
    go install $exe
done

echo Done
