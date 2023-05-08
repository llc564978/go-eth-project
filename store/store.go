package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/llc564978/ethereum-blockchain-service/db"
)

type Block struct {
	BlockNum   uint64
	BlockHash  string
	BlockTime  uint64
	ParentHash string
	TxHashes   string
}

type Transaction struct {
	TxHash string
	From   string
	To     string
	Nonce  uint64
	Data   string
	Value  string
	Logs   []map[string]interface{}
}

func InsertBlock(blockNum uint64, blockHash string, blockTime uint64, parentHash string, txHashes []string) {
	_, err := db.DB.Exec("INSERT INTO blocks (block_num, block_hash, block_time, parent_hash, tx_hashes) VALUES ($1, $2, $3, $4, $5)", blockNum, blockHash, blockTime, parentHash, pq.Array(txHashes))
	if err != nil {
		log.Fatalf("Failed to insert block: %v", err)
	}
}

func InsertTransaction(txHash string, blockNum uint64, from, to string, nonce uint64, data, value, logsJSON string) {
	_, err := db.DB.Exec("INSERT INTO transactions (tx_hash, block_num, from_address, to_address, nonce, data, value, logs) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", txHash, blockNum, from, to, nonce, data, value, logsJSON)
	if err != nil {
		log.Fatalf("Failed to insert transaction: %v", err)
	}
}

func GetLatestBlocks(limit int) ([]Block, error) {
	rows, err := db.DB.Query("SELECT block_num, block_hash, block_time, parent_hash FROM blocks ORDER BY block_num DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var block Block
		err := rows.Scan(&block.BlockNum, &block.BlockHash, &block.BlockTime, &block.ParentHash)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return blocks, nil
}

func GetBlockByID(blockID uint64) (*Block, error) {
	var block Block
	err := db.DB.QueryRow("SELECT block_num, block_hash, block_time, parent_hash, tx_hashes FROM blocks WHERE block_num = $1", blockID).Scan(&block.BlockNum, &block.BlockHash, &block.BlockTime, &block.ParentHash, &block.TxHashes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

func GetTransactionByHash(txHash string) (*Transaction, error) {
	var tx Transaction
	var logsJSON string
	err := db.DB.QueryRow("SELECT tx_hash, from_address, to_address, nonce, data, value, logs FROM transactions WHERE tx_hash = $1", txHash).Scan(&tx.TxHash, &tx.From, &tx.To, &tx.Nonce, &tx.Data, &tx.Value, &logsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var logs []map[string]interface{}
	if err := json.Unmarshal([]byte(logsJSON), &logs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs: %v", err)
	}
	tx.Logs = logs

	return &tx, nil
}

func BlockExists(blockNum uint64) (bool, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM blocks WHERE block_num = $1", blockNum).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
