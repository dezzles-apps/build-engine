# syntax=docker/dockerfile:1

FROM golang:1.26 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build-engine

FROM golang:1.26 AS runtime

COPY --from=builder /build-engine /build-engine
COPY config/config.prod.yaml /config/config.yaml
COPY config/config.non.yaml /config/config.non.yaml
EXPOSE 8080
ENTRYPOINT ["/build-engine"]