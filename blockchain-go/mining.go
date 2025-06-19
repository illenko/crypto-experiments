package main

import (
	"fmt"
	"strings"
	"time"
)

const DIFFICULTY = 4

type Miner struct {
	Address string
}

func NewMiner(address string) *Miner {
	return &Miner{
		Address: address,
	}
}

func (m *Miner) Mine(blockchain *Blockchain) *Block {
	fmt.Printf("⛏️ Miner %s started mining new block...\n", m.Address)

	previousHash := ""
	if lastBlock := blockchain.LastBlock(); lastBlock != nil {
		previousHash = lastBlock.Hash
	}

	candidateBlock := &Block{
		Index:        len(blockchain.Chain),
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Transactions: blockchain.PendingTransactions,
		PreviousHash: previousHash,
		Nonce:        0,
	}

	startTime := time.Now()
	for {
		candidateBlock.Hash = Hash(candidateBlock)
		if isValidProof(candidateBlock) {
			break
		}
		candidateBlock.Nonce++
	}
	duration := time.Since(startTime)

	fmt.Printf("💎 Block mined! Nonce: %d, Time: %v\n", candidateBlock.Nonce, duration)
	return candidateBlock
}

func isValidProof(block *Block) bool {
	return strings.HasPrefix(block.Hash, strings.Repeat("0", DIFFICULTY))
}
