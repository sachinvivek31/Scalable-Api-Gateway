# Step 1: Build the binary
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gateway main.go

# Step 2: Create a tiny runtime image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gateway .
COPY --from=builder /app/config.yaml .
EXPOSE 8080
CMD ["./gateway"]