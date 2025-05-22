FROM golang:1.23 AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go mod tidy

WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM alpine:3.20
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/cmd/main .
