#!/bin/bash

rm -f $ROOT_DIR/coverage.out
for pkg in $(go list ./... | grep -v vendor)
do
    go test -v -cover -race -coverprofile=$ROOT_DIR/coverage_tmp.out -timeout=1200s $pkg $@
    if [[ -f "$ROOT_DIR/coverage.out" ]]; then
        tail -n +2 $ROOT_DIR/coverage_tmp.out >> $ROOT_DIR/coverage.out
    else
        cat $ROOT_DIR/coverage_tmp.out > $ROOT_DIR/coverage.out
    fi
done
$GOPATH/bin/goveralls -coverprofile=$ROOT_DIR/coverage.out -service=travis-ci -repotoken $COVERALLS_TOKEN