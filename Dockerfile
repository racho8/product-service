# Dockerfile
FROM golang:1.24.5 AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Final image
FROM gcr.io/distroless/base-debian11

WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080

CMD ["/app/main"]
