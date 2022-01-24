# syntax=docker/dockerfile:1
FROM golang:1.17-alpine3.13 as builder
WORKDIR /go/src/ScooteServer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/scooter_server ./cmd

FROM scratch
COPY --from=builder /go/src/scooter_server/templates/. /home/scooter_server/templates
COPY --from=builder /go/src/scooter_server/build/scooter_server ./usr/bin/scooter_server
ENTRYPOINT [ "/usr/bin/scooter_server" ]