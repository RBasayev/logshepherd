#!/bin/bash

VER='0.3'
TS=$(date +%y%m%d.%H%M%S)

# docker run --rm -ti -v "$PWD:/host" -w /host -e CGO_ENABLED=0 -e GOOS=darwin -e GOARCH=arm64 golang go build -o logshepherd.iSilicon-mac -ldflags "-X main.Version=$VER.$TS"
# docker run --rm -ti -v "$PWD:/host" -w /host -e CGO_ENABLED=0 -e GOOS=darwin -e GOARCH=amd64 golang go build -o logshepherd.intel-mac -ldflags "-X main.Version=$VER.$TS"
# docker run --rm -ti -v "$PWD:/host" -w /host golang go build -o logshepherd.linux -ldflags "-X main.Version=$VER.$TS"

export CGO_ENABLED=0

export GOARCH=arm64
export GOOS=darwin
go build -o logshepherd.iSilicon-mac -ldflags "-X main.Version=$VER.$TS"

export GOARCH=amd64
export GOOS=darwin
go build -o logshepherd.intel-mac -ldflags "-X main.Version=$VER.$TS"

export GOARCH=amd64
export GOOS=linux
go build -o logshepherd.linux -ldflags "-X main.Version=$VER.$TS"
