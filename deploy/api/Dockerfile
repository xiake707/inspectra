FROM golang:1.26.5 AS builder
WORKDIR /inspector
COPY ./inspector ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /inspector/go-inspector main.go
WORKDIR /api
COPY ./api ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /api/goapi .
FROM debian:13.6
WORKDIR /Go
COPY --from=builder /inspector/go-inspector /Go/go-inspector
COPY --from=builder /api/goapi /Go/goapi
CMD ["/Go/goapi"]
