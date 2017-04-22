FROM alpine:latest

MAINTAINER Wang Xinguang <wangxinguang@e-business.co.jp>

COPY server /usr/bin/server

ENTRYPOINT ["/usr/bin/server"]
CMD ["/usr/bin/server"]
