FROM golang:1.16-buster

ADD . /go/src/github.com/ARMmaster17/Captain/
WORKDIR /go/src/
WORKDIR /go/src/github.com/ARMmaster17/Captain/flightplanner/
RUN go build -o captain-flightplanner

CMD ["/go/src/github.com/ARMmaster17/Captain/builder/captain-flightplanner"]