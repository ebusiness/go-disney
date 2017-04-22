#!/bin/sh

if ! which golint >/dev/null; then
  echo "install golint"
  go get -u github.com/golang/lint/golint
fi

# golint ./...
echo "golint"
golint config/*.go
golint middleware/*.go
golint utils/*.go
golint v1/*.go
golint v1/models/*.go

echo "go fmt"
# go fmt ./...
go fmt *.go
go fmt config/*.go
go fmt middleware/*.go
go fmt utils/*.go
go fmt v1/*.go
go fmt v1/models/*.go


echo "check vendor"
if [[ ! -d /go/src/github.com/ebusiness/go-disney/vendor ]]; then
  echo "get dep"

  apk add --update git
  go get -u github.com/golang/dep/...

  # dep init
  echo "ensure"

  dep ensure
  # dep ensure -update
fi
