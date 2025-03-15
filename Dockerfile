FROM golang:alpine AS builder
WORKDIR /build

COPY go.mod .
RUN go mod tidy

COPY . .
RUN env GOOS=linux GOARCH=arm64 go build -mod=mod -o main ./cmd/cdn/main.go

FROM alpine

WORKDIR /usr/app

COPY --from=builder /build/main /usr/app
ENV PRODUCTION=TRUE

ENTRYPOINT ./main 