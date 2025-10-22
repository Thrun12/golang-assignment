# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Allow Go to download required toolchain version
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./cmd/seed

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/server .
COPY --from=builder /app/migrate .
COPY --from=builder /app/seed .

# Copy migrations
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

# Copy API specifications
COPY --from=builder /app/api/proto/v1/*.swagger.json ./api/proto/v1/

# Expose ports
EXPOSE 8080 9090

# Run server by default
CMD ["./server"]
