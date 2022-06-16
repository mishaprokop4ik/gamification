# syntax=docker/dockerfile:1
FROM golang:1.18-alpine as builder

ENV GO111MODULE=on

WORKDIR /acheer
COPY .. .
RUN apk --no-cache add git alpine-sdk build-base gcc
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o acheer cmd/acheer/main.go

FROM alpine:3.15.4

RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=builder /acheer/acheer .
EXPOSE 8080

CMD ["./acheer"]