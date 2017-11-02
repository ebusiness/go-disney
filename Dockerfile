FROM golang:alpine AS GOALPINE

# ENV GIN_MODE=release

COPY ./ /go/src/github.com/e-business/go-disney
RUN cd /go/src/github.com/e-business/go-disney && \
    sh scripts/setup.sh && \
    go build -o /server main.go

FROM alpine:latest

LABEL maintainer: Wang Xinguang <wangxinguang@e-business.co.jp>

COPY --from=GOALPINE /server /usr/bin/server
COPY asset /asset

# RUN apk add --update tzdata && \
#    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
#    apk del tzdata && \
#    rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/bin/server"]
CMD ["/usr/bin/server"]
