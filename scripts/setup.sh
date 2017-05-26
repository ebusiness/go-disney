#!/bin/sh

if ! command -v golint >/dev/null; then

  if ! command -v git >/dev/null; then
    echo "install git"
    # alpine
    apk add --update git

  fi

  echo "install golint"
  go get -u github.com/golang/lint/golint
fi

# command -v golint >/dev/null 2>&1 || { echo >&2 "I require golint but it's not installed.  Aborting."; exit 1; }


# golint ./...
echo "golint"
golint config/*.go
golint middleware/*.go
golint utils/*.go
golint v1/*.go
golint v1/models/*.go
golint v1/algorithms/*.go

echo "go fmt"
# go fmt ./...
# go fmt *.go
# go fmt config/*.go
# go fmt middleware/*.go
# go fmt utils/*.go
# go fmt v1/*.go
# go fmt v1/models/*.go

gofmt -s -w -l  *.go
gofmt -s -w -l  config/*.go
gofmt -s -w -l  middleware/*.go
gofmt -s -w -l  utils/*.go
gofmt -s -w -l  v1/*.go
gofmt -s -w -l  v1/models/*.go
gofmt -s -w -l  v1/algorithms/*.go


echo "check vendor"
if [[ ! -d /go/src/github.com/ebusiness/go-disney/vendor ]]; then
  echo "get dep"

  if ! command -v dep >/dev/null; then

    if ! command -v git >/dev/null; then
      echo "install git"
      # alpine
      apk add --update git

    fi

    echo "install dep"
    go get -u github.com/golang/dep/...
  fi
  # dep init
  echo "ensure"

  dep ensure
  # dep ensure -update
fi
