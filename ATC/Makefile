BINARY=captain-atc
VERSION=`git describe --tags`
LDFLAGS=-ldflags="-w -s -X main.Version=$(VERSION)"

build:
	go build -v $(LDFLAGS) -o $(BINARY) github.com/ARMmaster17/Captain/ATC

test:
	$(MAKE) init-data-dir
	go get -u github.com/jstemmer/go-junit-report
	go test -v -coverprofile=cover.out ./... 2>&1 | go-junit-report > report.xml

clean:
	rm ./$(BINARY)

init-data-dir:
	sudo mkdir /etc/captain
	sudo chmod 777 -R /etc/captain

install:
	$(MAKE) deb
	sudo apt install ./pkg/captain-atc_*_amd64.deb
deb:
	$(MAKE) build
	cp $(BINARY) ./pkg/deb/usr/local/bin/$(BINARY)
	rm ./pkg/deb/usr/local/bin/.keep
	dpkg-deb --build ./pkg/deb
	mv ./pkg/deb.deb ./pkg/captain-atc_$(VERSION)_amd64.deb