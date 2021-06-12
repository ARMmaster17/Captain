FROM golang:1.15-alpine
COPY . /go/src/github.com/ARMmaster17/Captain/Builder
WORKDIR /go/src/github.com/ARMmaster17/Captain/Builder
RUN go build -o captain-builder .
# TODO: init config here.
EXPOSE 5000
CMD ["/go/src/github.com/ARMmaster17/Captain/Builder/captain-builder"]