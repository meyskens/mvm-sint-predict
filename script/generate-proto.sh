#!/bin/bash
if [[ "$GOPATH" == "" ]]; then
    export GOPATH="$HOME/go"
fi
export PATH=$PATH:$GOPATH/bin # fix for GOPATH not being in PATH
protoc -I $GOPATH/src --go_out=plugins=grpc:$GOPATH/src $GOPATH/src/github.com/meyskens/mvm-sint-predict/pb/sint.proto
