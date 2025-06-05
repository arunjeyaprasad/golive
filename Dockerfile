# Build stage
FROM golang:1.24.3-alpine3.21 AS builder

# Add necessary build tools
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /app/golive main.go

# Setup the permissions
RUN chmod +x /app/golive

# Final stage
FROM alpine:3.21

# Install ffmpeg and necessary runtime dependencies
RUN apk add --no-cache ffmpeg fontconfig ttf-dejavu ttf-droid ttf-freefont ttf-liberation

# Add basic security updates and CA certificates
RUN apk --no-cache add ca-certificates && \
    adduser -D -H -h /app appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/golive .

# Expose the port your app runs on
EXPOSE 9090

# Run the binary
CMD ["./golive"]