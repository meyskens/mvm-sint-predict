#!/bin/bash
if [[ "$GOPATH" == "" ]]; then
    export GOPATH="$HOME/go"
fi
export PATH=$PATH:$GOPATH/bin # fix for GOPATH not being in PATH
protoc -I "$(pwd)" --go_out=plugins=grpc:"$(pwd)" "$(pwd)"/pb/sint.proto
