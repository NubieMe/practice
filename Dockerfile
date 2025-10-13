# ========= Stage 1: Build =========
FROM golang:1.22-alpine AS builder

# Biar binary statis (tidak tergantung libc)
ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

# Copy go.mod dan go.sum dulu biar cache jalan
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build ke binary bernama 'app'
RUN go build -o app .

# ========= Stage 2: Run =========
FROM alpine:3.20

# Buat user non-root (lebih aman)
RUN adduser -D -u 10001 appuser

WORKDIR /app

# Copy binary hasil build
COPY --from=builder /app/app .

# Port default Fiber = 3000 (ubah kalau beda di main.go)
EXPOSE 3000

USER appuser

ENTRYPOINT ["./app"]
