package main

import (
	"context"
	"fmt"
	"net/http"

	_ "lukso-block-api/docs" // Import wygenerowanej dokumentacji

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger" // Import Swagger
)

const luksoTestnetRPC = "https://rpc.testnet.lukso.network"

// @title LUKSO Block API
// @version 1.0
// @description API for fetching the latest block number from LUKSO Testnet.
// @host localhost:8080
// @BasePath /

// @Schemes http
func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Swagger endpoint
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "UP",
		})
	})

	// Block number endpoint
	e.GET("/block-number", func(c echo.Context) error {
		blockNumber, err := getLatestBlockNumber()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch block number")
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"blockNumber": blockNumber,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

// getLatestBlockNumber fetches the latest block number from LUKSO Testnet.
func getLatestBlockNumber() (uint64, error) {
	client, err := rpc.DialContext(context.Background(), luksoTestnetRPC)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	var blockNumberHex string
	err = client.CallContext(context.Background(), &blockNumberHex, "eth_blockNumber")
	if err != nil {
		return 0, err
	}

	var blockNumber uint64
	_, err = fmt.Sscanf(blockNumberHex, "0x%x", &blockNumber)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}
