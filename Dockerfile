# Build stage
FROM golang:1.21-alpine3.18 AS builder

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .

# Build the Go application binary.
RUN go build -o /to-do-list

# Run stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /to-do-list /to-do-list

# Expose port 8080.
EXPOSE 8080

# Start the Go application.
CMD ["/to-do-list"]
