FROM golang:1.16-buster

ADD ./ /app/
WORKDIR /app/
RUN go build -o captain-builder

CMD ["/app/captain-builder"]