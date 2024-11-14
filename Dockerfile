# Stage 1: Build the Go application as a statically linked binary
FROM golang:1.23 AS builder

WORKDIR /app

# Copy all source code and the .env file
COPY . .

# Build the application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Use a minimal base image
FROM gcr.io/distroless/static-debian11

WORKDIR /app

# Copy the statically linked binary and .env file from the builder stage
COPY --from=builder /app /app

ENV PORT=8080  

# Run the application
CMD ["/app/main"]
