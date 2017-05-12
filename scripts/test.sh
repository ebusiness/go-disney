#!/bin/sh

# go test -v -covermode=count -coverprofile=coverage.middleware.out middleware/*.go
#
# go test -v -covermode=count -coverprofile=coverage.v1.out v1/*.go

go test -v middleware/*.go

go test -v v1/*.go

go test -v v1/algorithms/*.go

go test -v utils/*.go
