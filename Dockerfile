# Stage 1: Build the Go application
FROM golang:alpine as builder

# Set the working directory to the project root
WORKDIR /build

# Copy the project files to the working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create the runner image
FROM alpine:latest

# Set the working directory to the project root
WORKDIR /app

# Copy the Go binary and public folder from the builder image
COPY --from=builder /build/main /app/main
COPY --from=builder /build/public /app/public

# Copy the database.accdb file to the working directory
COPY database.accdb /app/database.accdb

# Expose port 8080
EXPOSE 8080

# Run the Go application
CMD ["./main"]
