#!/bin/bash
export GOROOT=/usr/lib/go
export GOPATH=$HOME/go

function cross-build() {
SUFFIX=
case $1 in
    386)
        BIT=32
        ;;
    amd64)
        BIT=64
        ;;
esac
case $2 in
    windows)
        DIR=win
        SUFFIX=.exe
        ;;
    darwin)
        DIR=macosx
        ;;
    linux)
        DIR=linux
        ;;
esac
echo "building netkits for $DIR$BIT"
OUT=bin/$DIR$BIT/netkits$SUFFIX
GOARCH=$1 GOOS=$2 CGO_ENABLED=0 \
    go build -o $OUT -ldflags "-s -w" netkits
echo "done [save as $PWD/$OUT]"
}

 # mac 64
cross-build amd64 darwin
 
 # win 32
cross-build 386 windows
 
 # win 64
cross-build amd64 windows
 
 # linux 32
cross-build 386 linux
 
 # linux 64
cross-build amd64 linux

