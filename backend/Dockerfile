# Backend Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Invalidate cache for source code changes
ARG CACHEBUST=1

# Copy source code
COPY . .

# Build the application with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o tonish-backend main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

# Copy the binary
COPY --from=builder /app/tonish-backend .

EXPOSE 8080

CMD ["./tonish-backend"]
