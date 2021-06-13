FROM golang:1.15-alpine
COPY ./Builder /go/src/github.com/ARMmaster17/Captain/Builder
COPY ./Shared /go/src/github.com/ARMmaster17/Captain/Shared
WORKDIR /go/src/github.com/ARMmaster17/Captain/Builder
RUN go build -o captain-builder .
# TODO: init config here.
EXPOSE 9090
CMD ["/go/src/github.com/ARMmaster17/Captain/Builder/captain-builder"]