# Step 1: Use an official Golang image as the base image
FROM golang:1.21 AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy the rest of the application source code
COPY . .

# Step 5: Build the Go application
RUN go build -o receipt-processor .

# Step 6: Use a lightweight runtime image with glibc
FROM debian:stable-slim

# Step 7: Set working directory for the runtime container
WORKDIR /app

# Step 8: Create logs directory
RUN mkdir -p /app/logs

# Step 9: Copy the compiled binary from the builder stage
COPY --from=builder /app/receipt-processor .

# Step 10: Expose the application port
EXPOSE 8080

# Step 11: Define the default command to run the application
CMD ["./receipt-processor"]