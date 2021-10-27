NAME=captain-atc
VERSION="v0.3.0-beta"
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
	rm ./$(NAME) || true
	rm ./pkg/*.deb || true
	rm ./*.db || true

install: $(NAME).deb
	sudo apt install ./pkg/captain-atc_*_amd64.deb -y

$(NAME).deb: $(NAME)
	cp $(NAME) ./pkg/deb/usr/local/bin/$(NAME)
	rm ./pkg/deb/usr/local/bin/.keep || true
	dpkg-deb --build ./pkg/deb
	mv ./pkg/deb.deb ./pkg/captain-atc_$(VERSION)_amd64.deb

all: install

.PHONY: install clean test build all