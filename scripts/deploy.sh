#!/bin/sh

# docker run -d --name redis redis:alpine

# mongo 3.4 ~
# docker run -d --name mongo mvertes/alpine-mongo

docker rmi ebusinessdocker/disney:latest
docker run -d --name disney --restart=always -p 443:443 --link redis:redis --link mongo:mongo -e REDISSERVER=redis -e MONGOSERVER=mongo -v $(pwd)/cert:/cert ebusinessdocker/disney:latest


docker run -d --name disney --restart=always -p 443:443 --link redis:redis --link mongo:mongo -e REDISSERVER=redis -e MONGOSERVER=mongo -v $(pwd)/ssl/.lego/certificates/api.dev.genbatomo.com.crt:/cert/cert.pem -v $(pwd)/ssl/.lego/certificates/api.dev.genbatomo.com.key:/cert/key.pem ebusinessdocker/disney:latest
