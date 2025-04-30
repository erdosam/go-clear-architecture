ARG GO_VERSION=1.23-alpine
# cache module
FROM golang:$GO_VERSION AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:$GO_VERSION AS development
WORKDIR /app
COPY --from=modules /go/pkg $GOPATH/pkg
COPY go.mod go.sum ./
RUN go install github.com/air-verse/air@latest
RUN mv $GOPATH/bin/air .
CMD ["/app/air", "-c", ".air.toml"]

FROM golang:$GO_VERSION AS builder
COPY --from=modules /go/pkg $GOPATH/pkg
# non cache
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /usr/local/bin/app ./cmd/app

FROM alpine:latest

WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/app .
CMD app