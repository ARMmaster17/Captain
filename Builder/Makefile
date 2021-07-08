NAME=captain-builder
VERSION=`git describe --tags`
LDFLAGS=-ldflags="-w -s -X metadata.Version=$(VERSION)"

test:
	go get -u github.com/jstemmer/go-junit-report
	go test -v -coverprofile=cover.out ./... 2>&1 | go-junit-report > report.xml

clean:
	rm ./$(NAME) || true
	rm ./pkg/*.deb || true
	rm ./*.db || true

.PHONY: clean test