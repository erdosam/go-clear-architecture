ARG GO_VERSION=1.23-alpine
ARG APP_PATH=$GOPATH/src/github.com/gilkor/ba-version-2/

FROM golang:$GO_VERSION as builder
# cache module
WORKDIR $APP_PATH
COPY ./go.* .
RUN go mod download
# non cache
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /usr/local/bin/app ./cmd/app

FROM alpine:latest

WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/app .
CMD app