FROM alpine:latest

MAINTAINER Wang Xinguang <wangxinguang@e-business.co.jp>

COPY server /usr/bin/server
COPY asset /asset

# RUN apk add --update tzdata && \
#    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
#    apk del tzdata && \
#    rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/bin/server"]
CMD ["/usr/bin/server"]
