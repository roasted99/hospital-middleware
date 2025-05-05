FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hospital-middleware ./cmd/server

# Use a minimal base image for the final image
FROM alpine:3.21

#Add the necessary CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/hospital-middleware .
COPY .env /app/.env

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./hospital-middleware"]

