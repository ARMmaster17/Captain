NAME=captain-radar
VERSION=`git describe --tags`

build: $(NAME)

$(NAME):
	go build -o $(NAME) github.com/ARMmaster17/Captain/Radar

test: $(NAME)
	go get -u github.com/jstemmer/go-junit-report
	go test -v -coverprofile=cover.out ./... 2>&1 | go-junit-report > report.xml

$(NAME).deb: $(NAME)
	cp $(NAME) ./pkg/deb/usr/local/bin/$(NAME)
	rm -R ./pkg/deb/etc/captain/radar/templates/ || true
	cp -R ./templates/ ./pkg/deb/etc/captain/radar/templates/
	rm ./pkg/deb/usr/local/bin/.keep || true
	dpkg-deb --build ./pkg/deb
	mv ./pkg/deb.deb ./pkg/captain-radar_$(VERSION)_amd64.deb

install: $(NAME).deb
	sudo apt install ./pkg/captain-radar_*_amd64.deb

all: install

clean:
	rm $(NAME) || true
	rm ./pkg/*.deb || true

.PHONY: clean all install test build