BINARY=captain
VERSION=`git describe --tags`
LDFLAGS=-ldflags="-w -s -X main.Version=${VERSION}"

build:
	go build -v ${LDFLAGS} -o ${BINARY} ./...

test:
	go test -v -coverprofile=cover.out ./...