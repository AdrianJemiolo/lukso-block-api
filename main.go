package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

const luksoTestnetRPC = "https://rpc.testnet.lukso.network"

var (
	privateKey1 *ecdsa.PrivateKey // First private key loaded from .env
	multiSigABI string            // ABI of the MultiSigWallet contract
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load the first private key from .env
	privateKeyHex := os.Getenv("PRIVATE_KEY_1")
	if privateKeyHex == "" {
		log.Fatal("PRIVATE_KEY_1 not set in .env")
	}

	privateKey1, err = crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	// Load ABI from JSON file
	abiBytes, err := ioutil.ReadFile("multisig_abi.json")
	if err != nil {
		log.Fatalf("Failed to load ABI file: %v", err)
	}

	multiSigABI = string(abiBytes)
	log.Println("ABI loaded successfully")
}

func main() {
	e := echo.New()

	// Endpoint: Get the latest block number
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

	// Endpoint: Submit a transaction signed by two keys
	e.POST("/submit-transaction", func(c echo.Context) error {
		type SubmitRequest struct {
			To    string `json:"to"`
			Value string `json:"value"` // Value in wei
			Data  string `json:"data"`  // Optional calldata
		}

		req := new(SubmitRequest)
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, "Invalid request body")
		}

		// Dynamically generate the second private key
		privateKey2, secondOwnerAddress, err := generatePrivateKey()
		if err != nil {
			log.Println("Error generating second private key:", err)
			return c.String(http.StatusInternalServerError, "Failed to generate second private key")
		}

		// Send the transaction
		txHash, err := sendTransaction(req.To, req.Value, req.Data, privateKey2, secondOwnerAddress)
		if err != nil {
			log.Println("Error submitting transaction:", err)
			return c.String(http.StatusInternalServerError, "Failed to submit transaction")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"txHash":              txHash,
			"secondOwnerAddress":  secondOwnerAddress,
			"generatedPrivateKey": hex.EncodeToString(crypto.FromECDSA(privateKey2)),
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func getLatestBlockNumber() (uint64, error) {
	client, err := ethclient.DialContext(context.Background(), luksoTestnetRPC)
	if err != nil {
		return 0, err
	}

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func generatePrivateKey() (*ecdsa.PrivateKey, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, "", err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return privateKey, address, nil
}

func checkIsOwner(client *ethclient.Client, contractAddress common.Address, parsedABI abi.ABI, addressToCheck common.Address) (bool, error) {
	callData, err := parsedABI.Pack("isOwner", addressToCheck)
	if err != nil {
		return false, fmt.Errorf("failed to pack call data for isOwner: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	output, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return false, fmt.Errorf("failed to call contract: %w", err)
	}

	var isOwner bool
	err = parsedABI.UnpackIntoInterface(&isOwner, "isOwner", output)
	if err != nil {
		return false, fmt.Errorf("failed to unpack isOwner response: %w", err)
	}

	return isOwner, nil
}

func getOwners(client *ethclient.Client, contractAddress common.Address, parsedABI abi.ABI) ([]common.Address, error) {
	callData, err := parsedABI.Pack("getOwners")
	if err != nil {
		return nil, fmt.Errorf("failed to pack call data for getOwners: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	output, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	var owners []common.Address
	err = parsedABI.UnpackIntoInterface(&owners, "getOwners", output)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack owners: %w", err)
	}

	return owners, nil
}

func sendTransaction(to, value, data string, privateKey2 *ecdsa.PrivateKey, secondOwnerAddress string) (string, error) {
	client, err := ethclient.DialContext(context.Background(), luksoTestnetRPC)
	if err != nil {
		return "", fmt.Errorf("failed to connect to LUKSO Testnet: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(multiSigABI))
	if err != nil {
		return "", fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	contractAddress := common.HexToAddress(os.Getenv("MULTISIG_ADDRESS"))

	owners, err := getOwners(client, contractAddress, parsedABI)
	if err != nil {
		log.Fatalf("Failed to get owners: %v", err)
	}
	log.Printf("Owners: %v", owners)

	fromAddress := crypto.PubkeyToAddress(privateKey1.PublicKey)
	isOwnerContract, err := checkIsOwner(client, contractAddress, parsedABI, fromAddress)
	if err != nil {
		log.Fatalf("Error checking isOwner in contract: %v", err)
	}
	log.Printf("isOwner in contract for %s: %v", fromAddress.Hex(), isOwnerContract)

	if !isOwnerContract {
		return "", fmt.Errorf("address %s is not an owner", fromAddress.Hex())
	}

	valueInWei, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return "", fmt.Errorf("invalid value: %s", value)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to fetch chain ID: %w", err)
	}

	suggestedTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to fetch gas tip cap: %w", err)
	}
	latestBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest block: %w", err)
	}
	baseFee := latestBlock.BaseFee()
	maxFeePerGas := new(big.Int).Add(baseFee, suggestedTipCap)

	toAddress := common.HexToAddress(to)
	callData, err := parsedABI.Pack("submitTransaction", toAddress, valueInWei, common.Hex2Bytes(data))
	if err != nil {
		return "", fmt.Errorf("failed to pack call data: %w", err)
	}

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	})
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas: %w", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to fetch nonce: %w", err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &contractAddress,
		Value:     big.NewInt(0),
		Gas:       gasLimit,
		GasFeeCap: maxFeePerGas,
		GasTipCap: suggestedTipCap,
		Data:      callData,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey1)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	log.Printf("Transaction submitted. Hash: %s", signedTx.Hash().Hex())

	return signedTx.Hash().Hex(), nil
}
