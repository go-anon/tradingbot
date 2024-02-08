# syntax=docker/dockerfile:1

# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------

FROM golang:1.21.5 AS builder

WORKDIR /go/src/


COPY ../tradingbot ./

RUN go get -d -v golang.org/x/net/html

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd/bot/main .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/main ./tradingbot
ENTRYPOINT ["./tradingbot"]
