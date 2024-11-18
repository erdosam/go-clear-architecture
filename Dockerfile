FROM golang:1.18-alpine as builder

WORKDIR $GOPATH/src/github.com/gilkor/athena/
COPY ./go.* .
RUN go mod download
COPY . .
RUN go mod tidy
RUN go build -o /usr/local/bin/app

FROM alpine:latest

WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/app .
CMD app