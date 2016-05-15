#!/bin/bash

packages=$(find . -type d | grep -v git)

IFS=$'\n'
for pkg in $packages; do
    pushd $pkg
    go test -coverprofile=.coverprofile -test.v -covermode=count
    popd
done
echo "mode: count" > total.coverprofile
cat $(find . -name .coverprofile) | grep -v "mode: count" >> total.coverprofile

