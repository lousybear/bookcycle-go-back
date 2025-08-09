FROM golang:1.24-alpine

WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the app
RUN go build -o server ./cmd/server

# Expose app port
EXPOSE 8080

# Run the app
CMD ["./server"]