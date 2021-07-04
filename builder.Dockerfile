FROM golang:1.15-alpine
RUN apk add gcc musl-dev
COPY ./Builder /go/src/github.com/ARMmaster17/Captain/builder
COPY ./Shared /go/src/github.com/ARMmaster17/Captain/Shared
WORKDIR /go/src/github.com/ARMmaster17/Captain/builder
RUN go build -o captain-builder ./cmd/main.go
EXPOSE 3000
ENV CAPTAIN_DB=TEST
ENV CAPTAIN_BUILDER_API_PORT=3000
CMD ["/go/src/github.com/ARMmaster17/Captain/builder/captain-builder"]