FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -trimpath -ldflags='-s -w' -o /out/regional-authoritative-dns ./cmd/server
FROM alpine:3.20
RUN adduser -D -u 10001 dns
COPY --from=build /out/regional-authoritative-dns /usr/local/bin/regional-authoritative-dns
COPY configs /app/configs
USER dns
WORKDIR /app
EXPOSE 8080 5353/udp 5353/tcp
ENTRYPOINT ["/usr/local/bin/regional-authoritative-dns"]
