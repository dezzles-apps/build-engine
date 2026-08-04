# syntax=docker/dockerfile:1

FROM golang:1.26 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build-engine

FROM golang:1.26 AS runtime
WORKDIR /app
COPY --from=builder /build-engine /app/build-engine
COPY config/config.prod.yaml /app/config/config.prod.yaml
COPY config/config.non.yaml /app/config/config.non.yaml
EXPOSE 8080
ENTRYPOINT ["/app/build-engine"]