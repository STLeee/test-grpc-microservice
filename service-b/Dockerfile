# Use the official Golang image as the base image
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY service-b/go.mod service-b/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o out/service-b service-b/main.go

# Use a image for the final container
FROM debian:12.10-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/out/service-b .

# Expose the port the service listens on
EXPOSE 50052

# Command to run the application
CMD ["/app/service-b"]
