FROM golang:1.23-alpine as builder

WORKDIR $GOPATH/src/github.com/gilkor/ba-version-2/
COPY ./go.* .
RUN go mod download
COPY . .
RUN go mod tidy
RUN go build -o /usr/local/bin/app

FROM alpine:latest

WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/app .
CMD app