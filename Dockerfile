# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o potion .

# Runtime stage
FROM alpine:latest
# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/potion .
# Run application
CMD ["./potion"]
