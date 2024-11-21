FROM golang:1.22-alpine

WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Ensure all dependencies are up-to-date
RUN go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
