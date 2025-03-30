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
RUN go build -o service-b service-b/main.go

# Use a minimal image for the final container
FROM gcr.io/distroless/base-debian11

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/service-b .

# Expose the port the service listens on
EXPOSE 50052

# Command to run the application
CMD ["/app/service-b"]
