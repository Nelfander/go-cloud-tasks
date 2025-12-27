# STEP 1: Build the binary
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Optimization: Copy and download dependencies first (improves build speed)
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build the app specifically for a generic Linux environment
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# STEP 2: Run the binary
FROM alpine:latest
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# NEW: Copy ALL html files (index, login, and register)
# Using a wildcard *.html is safer and cleaner!
COPY --from=builder /app/*.html ./



EXPOSE 8080
CMD ["./main"]      