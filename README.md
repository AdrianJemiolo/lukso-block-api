# LUKSO Block API

This project provides a simple HTTP API to fetch the latest block number from the LUKSO Testnet blockchain. It is built using the [Echo](https://echo.labstack.com/) framework for HTTP routing and the [go-ethereum](https://github.com/ethereum/go-ethereum) library for blockchain interaction.

## Features

- Fetches the latest block number from the LUKSO Testnet using the RPC endpoint.
- Built with Go, leveraging modern frameworks and libraries.
- Dockerized for easy deployment.

---

## Requirements

- Go 1.22 or later
- Docker (optional, for containerized deployment)

---

## Getting Started

### Clone the Repository

```bash
git clone git@github.com:<your-username>/lukso-block-api.git
cd lukso-block-api
```
### Install Dependencies
```bash
go mod tidy
```
### Run Locally
To run the application locally without Docker:
```bash
go run main.go
```
The server will start on port `8080`. Access the API at:
```plaintext
http://localhost:8080/block-number
```
---
### Usage
**API Endpoint**
**GET** `/block-number`
Fetches the latest block number from the LUKSO Testnet.

Response Example:
```json
{
  "blockNumber": 12345678
}
```
---
### Docker Deployment
Build Docker Image:
```bash
docker build -t lukso-block-api .
```
Run Docker Container:
```bash
docker run -p 8080:8080 lukso-block-api
```
Access the API at:
```plaintext
http://localhost:8080/block-number
```
---
### Docker Deployment (Docker Compose)
You can use Docker Compose to simplify the deployment process:

1. Build and run the application:
```bash
docker-compose up -d
```
The application will be accessible at:
```plaintext
http://localhost:8080/block-number
```
2. Stop and remove containers:
```bash
docker-compose down
```
3. Check logs (optional):
```bash
docker-compose logs -f
```

---
### Notes
- The LUKSO Testnet RPC endpoint is hardcoded in the project: https://rpc.testnet.lukso.network.
- Modify the Dockerfile to update dependencies or build parameters if required.
---
### License
This project is licensed under the MIT License 🙂
