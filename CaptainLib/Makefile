NAME=captain-lib
VERSION="v0.3.0-beta"

build: $(NAME)

$(NAME):
	go build -o $(NAME) github.com/ARMmaster17/Captain/CaptainLib

test: $(NAME)
	go get -u github.com/jstemmer/go-junit-report
	go test -v -coverprofile=cover.out ./... 2>&1 | go-junit-report > report.xml

all: $(NAME)

clean:
	rm $(NAME)

.PHONY: clean all test build