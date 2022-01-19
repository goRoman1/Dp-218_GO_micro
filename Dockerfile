FROM golang:1.17-alpine3.13 as builder
WORKDIR /go/src/Dp218GO
COPY . .
ENV GOPROXY https://proxy.golang.org,direct
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/scooterapp ./cmd/app

FROM scratch
COPY --from=builder /go/src/Dp218GO/migrations/. /home/Dp218GO/migrations
COPY --from=builder /go/src/Dp218GO/templates/. /home/Dp218GO/templates
COPY --from=builder /go/src/Dp218GO/build/scooterapp /usr/bin/scooterapp
COPY --from=builder /go/src/Dp218GO/microservice/certificates/. /home/certificates
ENTRYPOINT [ "/usr/bin/scooterapp" ]