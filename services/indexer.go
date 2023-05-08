package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llc564978/ethereum-blockchain-service/store"
)

func IndexBlocks(ctx context.Context, rpcEndpoint string, startBlock uint64) {
	client, err := ethclient.DialContext(ctx, rpcEndpoint)
	if err != nil {
		panic(fmt.Errorf("failed to connect to Ethereum node: %v", err))
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		panic(fmt.Errorf("failed to retrieve latest header: %v", err))
	}

	latestBlock := header.Number.Uint64()
	fmt.Printf("Latest block number: %d\n", latestBlock)

	for blockNumber := startBlock; blockNumber <= latestBlock; blockNumber++ {
		block, err := client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err != nil {
			fmt.Printf("Failed to fetch block %d: %v\n", blockNumber, err)
			continue
		}

		if block == nil {
			fmt.Printf("Block %d is nil\n", blockNumber)
			continue
		}

		fmt.Printf("Processing block %d (hash=%v)\n", blockNumber, block.Hash().String())
		processBlock(client, block)
	}
}

func processBlock(client *ethclient.Client, block *types.Block) {
	blockNum := block.NumberU64()
	blockHash := block.TxHash().Hex()
	blockTime := uint64(block.Time())
	parentHash := block.ParentHash().Hex()

	// Check block exists
	exists, err := store.BlockExists(blockNum)
	if err != nil {
		log.Fatalf("Failed to check if block exists: %v", err)
	}

	if !exists {

		txHashes := []string{}
		for _, tx := range block.Transactions() {
			txHashes = append(txHashes, tx.Hash().Hex())
		}

		// Insert Data to Block Table
		store.InsertBlock(blockNum, blockHash, blockTime, parentHash, txHashes)

		for _, tx := range block.Transactions() {
			txHash := tx.Hash().Hex()
			from, _ := client.TransactionSender(context.Background(), tx, block.Hash(), 0)
			to := tx.To()
			toStr := "Contract creation"
			if to != nil {
				toStr = to.Hex()
			}
			value := tx.Value().String()
			nonce := tx.Nonce()
			data := hex.EncodeToString(tx.Data())

			// Get TransactionReceipt
			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				panic(fmt.Errorf("failed to retrieve transaction receipt: %v", err))
			}
			logs := []map[string]interface{}{}
			for i, log := range receipt.Logs {
				logs = append(logs, map[string]interface{}{
					"index": i,
					"data":  hex.EncodeToString(log.Data),
				})
			}
			logsJSON, err := json.Marshal(logs)
			if err != nil {
				panic(fmt.Errorf("failed to marshal logs: %v", err))
			}

			// Insert Data to Transaction Table
			store.InsertTransaction(txHash, blockNum, from.Hex(), toStr, nonce, data, value, string(logsJSON))
		}
	}
}
