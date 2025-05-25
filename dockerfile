# Step 1: Use an official Golang image as a build environment
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Create a smaller image to run the application
FROM alpine:latest  

# Install necessary dependencies for running the app
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the build stage
COPY --from=builder /app/main .

# Expose port 2112 for Prometheus scraping
EXPOSE 2112

# Command to run the executable
CMD ["./main"]

