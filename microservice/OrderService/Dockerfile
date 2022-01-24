# syntax=docker/dockerfile:1
FROM golang:1.17-alpine3.13 as builder
WORKDIR /go/src/OrderService
COPY . .
ENV GO111MODULE=on
ENV GOPROXY https://proxy.golang.org,direct
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/order_service

FROM scratch
COPY --from=builder /go/src/order_service/build/order_service ./usr/bin/order_service
ENTRYPOINT [ "/usr/bin/order_service" ]