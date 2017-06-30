#!/bin/sh

# go test -v -covermode=count -coverprofile=coverage.middleware.out middleware/*.go
#
# go test -v -covermode=count -coverprofile=coverage.v1.out v1/*.go

richgo test -v middleware/*.go

richgo test -v v1/*.go

richgo test -v v1/algorithms/*.go

richgo test -v utils/*.go
