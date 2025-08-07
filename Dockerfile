# Stage 1: Build the Go binary
FROM golang:1.21.8-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod .
COPY go.sum .

# Install build dependencies
RUN apk add --no-cache git

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o product-service .

# Stage 2: Create minimal runtime image
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/product-service .

ENV PORT=8080

# Add non-root user
RUN adduser -D appuser
USER appuser

CMD ["./product-service"]