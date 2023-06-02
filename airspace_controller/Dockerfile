FROM golang:1.20-alpine AS gobuild
COPY . /src/
WORKDIR /src
RUN go build ./cmd/main.go

FROM node:18.15-alpine AS nodebuild
COPY . /src/
WORKDIR /src/frontend
RUN npm install
RUN npm run build

FROM alpine:3.18
COPY --from=gobuild /src/main /app/main
COPY --from=nodebuild /src/frontend/dist/ /app/frontend/dist/
CMD ["/app/main"]