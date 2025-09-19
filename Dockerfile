# syntax=docker/dockerfile:1.7

FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/app ./cmd/main.go

FROM gcr.io/distroless/static-debian12 AS runtime
WORKDIR /app

COPY --from=builder /out/app /app/app
COPY config/ config/

EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/app/app"]