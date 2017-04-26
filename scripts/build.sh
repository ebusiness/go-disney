#!/bin/sh

# -e GIN_MODE=release
rm server
docker run -it --rm -v $(pwd):/go/src/github.com/ebusiness/go-disney golang:alpine sh -c 'cd /go/src/github.com/ebusiness/go-disney && sh scripts/setup.sh && go build -o server main.go'
docker build --tag ebusinessdocker/disney .
