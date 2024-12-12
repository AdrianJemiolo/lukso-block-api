
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
# Pobierz obraz bazowy
FROM golang:1.22-alpine

# Ustaw katalog roboczy
WORKDIR /app

# Zainstaluj wymagane narzędzia
RUN apk add --no-cache git

# Zainstaluj narzędzie swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Skopiuj pliki go.mod i go.sum i pobierz zależności
COPY go.mod go.sum ./
RUN go mod download

# Skopiuj cały kod źródłowy
COPY . .

# Wygeneruj dokumentację Swagger
RUN swag init

# Zbuduj aplikację
RUN go build -o main .

# Wystaw port aplikacji
EXPOSE 8080

# Uruchom aplikację
CMD ["./main"]
