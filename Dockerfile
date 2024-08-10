# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy the Go Modules manifests and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application (replace `myapp` with your application's binary name)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o onecmsbackend

# Final stage
FROM alpine:3.12

# Install necessary runtime dependencies, if any
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/onecmsbackend .

# Copy the .env file from the project root into the container
COPY --from=builder /app/.env .

# Ensure the application binary is executable
RUN chmod +x ./onecmsbackend

# Command to run the binary
CMD ["./onecmsbackend"]
