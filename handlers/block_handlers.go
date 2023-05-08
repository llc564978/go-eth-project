package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/llc564978/ethereum-blockchain-service/store"
)

func GetBlocksHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	blocks, err := store.GetLatestBlocks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blocks": blocks})
}

func GetBlockByIDHandler(c *gin.Context) {
	blockIDStr := c.Param("id")
	blockID, err := strconv.ParseUint(blockIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block ID"})
		return
	}

	block, err := store.GetBlockByID(blockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if block == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
		return
	}

	transactions := strings.Split(block.TxHashes, ",")
	c.JSON(http.StatusOK, gin.H{
		"block_num":    block.BlockNum,
		"block_hash":   block.BlockHash,
		"block_time":   block.BlockTime,
		"parent_hash":  block.ParentHash,
		"transactions": transactions,
	})
}

func GetTransactionByHashHandler(c *gin.Context) {
	txHash := c.Param("txHash")

	tx, err := store.GetTransactionByHash(txHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tx == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tx_hash": tx.TxHash,
		"from":    tx.From,
		"to":      tx.To,
		"nonce":   tx.Nonce,
		"data":    tx.Data,
		"value":   tx.Value,
		"logs":    tx.Logs,
	})
}
