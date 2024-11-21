package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/labstack/echo/v4"
)

const luksoTestnetRPC = "https://rpc.testnet.lukso.network"

func main() {
	e := echo.New()

	e.GET("/block-number", func(c echo.Context) error {
		blockNumber, err := getLatestBlockNumber()
		if err != nil {
			log.Println("Error fetching block number:", err)
			return c.String(http.StatusInternalServerError, "Failed to fetch block number")
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"blockNumber": blockNumber,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

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
