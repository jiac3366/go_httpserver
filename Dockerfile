FROM ubuntu

LABEL lan="golang" app="httpserver" maintainer="jiac"

COPY bin/amd64/httpserver_gin /httpserver_gin

ENV PORT=

EXPOSE $PORT

ENTRYPOINT ["/httpserver_gin"]
