# Build stage
FROM golang:1.24-alpine3.20 AS builder

WORKDIR /app
ENV CGO_ENABLED=0 \
    GO111MODULE=on
COPY go.mod go.sum .
RUN go mod download
COPY . .

# Build the Go application binary (stripped)
RUN go build -trimpath -ldflags="-s -w" -o /to-do-list

# Run stage
FROM alpine:latest

WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=builder /to-do-list /to-do-list
USER appuser

# Expose port 8080.
EXPOSE 8080

# Start the Go application.
ENV GIN_MODE=release
CMD ["/to-do-list"]
