# syntax=docker/dockerfile:1
FROM golang:1.17-alpine3.13 as builder
WORKDIR /go/src/ScooterClient
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/scooter_client ./cmd

FROM scratch
COPY --from=builder /go/src/scooter_client/build/scooter_client ./usr/bin/scooter_client
ENTRYPOINT [ "/usr/bin/scooter_client" ]