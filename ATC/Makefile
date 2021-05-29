NAME=captain-atc
VERSION=`git describe --tags`
LDFLAGS=-ldflags="-w -s -X main.Version=$(VERSION)"

build: $(NAME)

$(NAME):
	go build -v $(LDFLAGS) -o $(NAME) github.com/ARMmaster17/Captain/ATC

test: $(NAME)
	sudo mkdir /etc/captain
	sudo mkdir /etc/captain/atc
	sudo chmod 777 -R /etc/captain
	go get -u github.com/jstemmer/go-junit-report
	go test -v -coverprofile=cover.out ./... 2>&1 | go-junit-report > report.xml

clean:
	rm ./$(NAME)
	rm ./pkg/*.deb
	rm ./*.db

install: $(NAME).deb
	$(MAKE) deb
	sudo apt install ./pkg/captain-atc_*_amd64.deb

$(NAME).deb: $(NAME)
	$(MAKE) build
	cp $(NAME) ./pkg/deb/usr/local/bin/$(NAME)
	rm ./pkg/deb/usr/local/bin/.keep
	dpkg-deb --build ./pkg/deb
	mv ./pkg/deb.deb ./pkg/captain-atc_$(VERSION)_amd64.deb

all: install

.PHONY: install clean test build all