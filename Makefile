BINARY=captain
VERSION=`git describe --tags`
LDFLAGS=-ldflags="-w -s -X main.Version=$(VERSION)"

build:
	go build -v $(LDFLAGS) -o $(BINARY) ./...

test:
	$(MAKE) init-data-dir
	go test -v -coverprofile=cover.out ./...

clean:
	rm ./$(BINARY)

uninstall:
	rm /usr/local/bin/$(BINARY)
	rm /etc/captain/defaults.yaml
	rmdir /etc/captain

init-data-dir:
	mkdir /etc/captain
	cp defaults.yaml /etc/captain/defaults.yaml

install:
	$(MAKE) build
	cp $(BINARY) /usr/local/bin/$(BINARY)
	$(MAKE) init-data-dir