#!/bin/bash

if [ ! -n "$BASH" ] ;then
  # echo Please run this script $0 with bash;
  if ! command -v bash >/dev/null; then
    echo "install bash"
    apk add --update bash
  fi
  bash ./scripts/setup.sh
  exit 0;
fi


setup::install::git() {
  if ! command -v git >/dev/null; then
    echo "install git"
    # alpine
    apk add --update git
  fi
}
setup::install::golint() {
  setup::install::git
  echo "install golint"
  go get -u github.com/golang/lint/golint
}

setup::install::dep() {
  if ! command -v dep >/dev/null; then
    setup::install::git
    echo "install dep"
    go get -u github.com/golang/dep/...
  fi
}

setup::golint() {
  if ! command -v golint >/dev/null; then
    setup::install::golint
  fi
  echo "golint"
  golint config/*.go
  golint middleware/*.go
  golint utils/*.go
  golint v1/*.go
  golint v1/models/*.go
  golint v1/algorithms/*.go
}

setup::gofmt() {
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
}


setup::vendor() {
  echo "check vendor"
  if [[ ! -d /go/src/github.com/ebusiness/go-disney/vendor ]]; then
    setup::install::dep
    # dep init
    echo "ensure"
    dep ensure
  fi
}

setup::golint
setup::gofmt
setup::vendor
exit 0;
