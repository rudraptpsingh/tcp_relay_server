# Dockerfile
FROM golang:1.22-alpine as builder

WORKDIR /go/tcp-relay-server
COPY . .

RUN go build -o tcp-relay-server ./relay/relay.go

FROM alpine:latest
EXPOSE 9009
EXPOSE 9010
EXPOSE 9011
EXPOSE 9012
EXPOSE 9013
COPY --from=builder /go/tcp-relay-server/tcp-relay-server .
USER nobody

CMD ["./tcp-relay-server"]
