export tag=v4.0
build:
	echo "building httpserver_gin binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver_gin container"
	docker build -t 463045792/apiserver:${tag} .

push: release
	echo "pushing docker image"
	docker push 463045792/apiserver:${tag}