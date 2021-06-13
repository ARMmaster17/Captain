FROM golang:1.15-alpine
COPY ./Preflight /go/src/github.com/ARMmaster17/Captain/Preflight
COPY ./Shared /go/src/github.com/ARMmaster17/Captain/Shared
WORKDIR /go/src/github.com/ARMmaster17/Captain/Preflight
RUN go build -o captain-preflight .
# TODO: init config here.
EXPOSE 9090
CMD ["/go/src/github.com/ARMmaster17/Captain/Preflight/captain-preflight"]