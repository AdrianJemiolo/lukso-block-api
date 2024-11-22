# Pull the base image
FROM golang:1.22-alpine

# Set the working directory
WORKDIR /app

# Install required tools
RUN apk add --no-cache git

# Install swag tool
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Generate Swagger documentation
RUN swag init

# Build the application
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
