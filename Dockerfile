# STEP 1: Build the binary
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Build the app specifically for a generic Linux environment
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# STEP 2: Run the binary
FROM alpine:latest
WORKDIR /root/
# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Copy your static files and .env
COPY --from=builder /app/index.html .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]