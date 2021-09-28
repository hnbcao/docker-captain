#!/bin/sh
export GOOS=linux
export GOARCH=amd64
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
echo "building docker images for ${GOOS}/${GOARCH} ..."

#export CGO_ENABLED=0
go mod tidy && go mod vendor

go build -mod=vendor -o release/linux/${GOARCH}/server ./
