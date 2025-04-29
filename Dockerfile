FROM golang:1.23 AS builder

# compile the codebase
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o /build/bitcaskDB-server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /build/bitcaskDB-server /app
ENTRYPOINT ["/app/bitcaskDB-server"]
