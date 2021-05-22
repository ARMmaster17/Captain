BINARY=captain-radar
VERSION=`git describe --tags`

build:
	go build -o $(BINARY) github.com/ARMmaster17/Captain/Radar

test:
	go get -u github.com/jstemmer/go-junit-report
	go test -v -coverprofile=cover.out ./... 2>&1 | go-junit-report > report.xml

install:
	$(MAKE) deb
	sudo apt install ./pkg/captain-radar_*_amd64.deb

deb:
	$(MAKE) build
	cp $(BINARY) ./pkg/deb/usr/local/bin/$(BINARY)
	rm -R ./pkg/deb/etc/captain/radar/templates/
	cp -R ./templates/ ./pkg/deb/etc/captain/radar/templates/
	rm ./pkg/deb/usr/local/bin/.keep
	dpkg-deb --build ./pkg/deb
	mv ./pkg/deb.deb ./pkg/captain-radar_$(VERSION)_amd64.deb