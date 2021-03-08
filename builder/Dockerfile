FROM golang:1.16-buster

ADD . /go/src/github.com/ARMmaster17/Captain/
ADD ./builder/conf/ /etc/captain/builder/
WORKDIR /go/src/
#ENV GOPATH=/go/src/
#ENV GO111MODULE=auto
WORKDIR /go/src/github.com/ARMmaster17/Captain/builder/
RUN go build -o captain-builder

CMD ["/go/src/github.com/ARMmaster17/Captain/builder/captain-builder"]