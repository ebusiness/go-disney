#!/bin/sh

# docker run -d --name redis redis:alpine
# docker run -d --name mongo mvertes/alpine-mongo

docker rmi ebusinessdocker/disney:latest
docker run -it --rm -p 443:443 --link redis:redis --link mongo:mongo -e REDISSERVER=redis -e MONGOSERVER=mongo -v $(pwd)/cert:/cert ebusinessdocker/disney:latest
