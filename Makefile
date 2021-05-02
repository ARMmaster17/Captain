BINARY=captain
VERSION=`git describe --tags`
LDFLAGS=-ldflags="-w -s -X main.Version=$(VERSION)"

build:
	go build -v $(LDFLAGS) -o $(BINARY) github.com/ARMmaster17/Captain

test:
	$(MAKE) init-data-dir
	go get -u github.com/jstemmer/go-junit-report
	go test -v -coverprofile=cover.out ./... 2>&1 | go-junit-report > report.xml

clean:
	rm ./$(BINARY)

uninstall:
	rm /usr/local/bin/$(BINARY)
	rm /etc/captain/defaults.yaml
	rmdir /etc/captain

init-data-dir:
	sudo mkdir /etc/captain
	sudo chmod 777 -R /etc/captain
	cp defaults.yaml /etc/captain/defaults.yaml

install-service:
	sudo cp captain.service /etc/systemd/system/captain.service
	sudo systemctl daemon-reload
	sudo systemctl enable captain
	sudo systemctl start captain

install:
	$(MAKE) build
	cp $(BINARY) /usr/local/bin/$(BINARY)
	$(MAKE) init-data-dir