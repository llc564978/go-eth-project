package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/llc564978/ethereum-blockchain-service/handlers"
	"github.com/llc564978/ethereum-blockchain-service/services"
)

func main() {
	router := gin.Default()

	router.GET("/blocks", handlers.GetBlocksHandler)
	router.GET("/blocks/:id", handlers.GetBlockByIDHandler)
	router.GET("/transaction/:txHash", handlers.GetTransactionByHashHandler)

	rpcEndpoint := "https://data-seed-prebsc-2-s3.binance.org:8545/"
	startBlock := uint64(29471119)

	ctx := context.Background()

	go services.IndexBlocks(ctx, rpcEndpoint, startBlock)

	router.Run(":8080")
}
