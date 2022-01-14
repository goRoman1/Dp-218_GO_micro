FROM golang:1.17-alpine3.13 as builder
WORKDIR /go/src/Dp218Go
COPY ../1/Dp-218_Go .
ENV GOPROXY https://proxy.golang.org,direct
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/scooterapp ./cmd/app

FROM scratch
COPY --from=builder /go/src/Dp218Go/migrations/. /home/Dp218Go/migrations
COPY --from=builder /go/src/Dp218Go/templates/. /home/Dp218Go/templates
COPY --from=builder /go/src/Dp218Go/build/scooterapp /usr/bin/scooterapp
COPY --from=builder /go/src/Dp218Go/microservice/certificates/. /home/certificates
ENTRYPOINT [ "/usr/bin/scooterapp" ]