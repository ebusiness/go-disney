#!/bin/sh

# -e GIN_MODE=release
docker rmi ebusinessdocker/disney
docker build --tag ebusinessdocker/disney --no-cache .
docker push ebusinessdocker/disney
docker rmi $(docker images -f dangling=true -q)