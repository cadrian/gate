#!/bin/bash

export GOPATH=$(dirname $(readlink -f $0))
export PATH=$GOPATH/bin:"$PATH"

echo Fetching deps
go install github.com/sbinet/liner
go install code.google.com/p/gomock/gomock
go install code.google.com/p/gomock/mockgen

echo Testing
TESTS=$(find src/gate -name \*_test.go -exec dirname {} \; | uniq | cut -c5-)
go test -i $TESTS
go test $TESTS || exit 1

find src -name main.go -print | while read main_go; do
    main=$(dirname $main_go)
    exe=${main#src/}
    echo Building $(basename $exe)
    go install $exe
done

echo Done
