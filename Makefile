default: build

dockerbuild:
	docker build --rm --tag=docker.blnk.io/darksky-api:latest .
build:
	docker run --rm -it -v ${PWD}:/go/src/${PWD} -w /go/src/${PWD} golang:latest make gobuild

gobuild:  *.go
	CGO_ENABLED=0 GOOS=linux go build .

run:
	docker run -d --name darksky_api --network=isolated_nw -e DARKSKY_KEY=${DARKSKY_KEY} docker.blnk.io/darksky-api
