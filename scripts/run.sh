#!/bin/sh

# docker run -d --name redis redis:alpine
# docker run -d --name mongo mvertes/alpine-mongo

docker run -it --rm \
      --link redis:redis \
      -e REDISSERVER=redis \
      --link mongo:mongo \
      -e MONGOSERVER=mongo \
      -p 8080:443 \
      -v $(pwd):/go/src/github.com/ebusiness/go-disney \
      golang:alpine sh -c '

cd /go/src/github.com/ebusiness/go-disney

echo "format"
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

echo "run"
# dev
go run main.go

# release
# go build -o server main.go
'
