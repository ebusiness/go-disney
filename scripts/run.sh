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

sh scripts/setup.sh

echo "run"
# dev
go run main.go

# release
# go build -o server main.go
'
