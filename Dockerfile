##
## Build
##

FROM golang:1.16-buster AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY . .

RUN ls; \
    go env -w GOPROXY="https://goproxy.io,direct"; \
    go mod vendor; \
    CGO_ENABLED=0 GOARCH=amd64 go build -o /bin/httpserver


##
## Deploy
##
FROM scratch

LABEL lan="golang" app="httpserver" maintainer="jiac"

WORKDIR /

COPY --from=build /bin/httpserver /bin/httpserver

ENV VERSION=1.0

ENV PORT=

EXPOSE $PORT

ENTRYPOINT ["/bin/httpserver"]

