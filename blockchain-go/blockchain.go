package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Blockchain struct {
	chain               []*Block
	pendingTransactions []*Transaction
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		chain:               make([]*Block, 0),
		pendingTransactions: make([]*Transaction, 0),
	}
	bc.NewBlock("") // Create genesis block
	return bc

}

func (b *Blockchain) NewBlock(previousHash string) *Block {
	block := &Block{
		Index:        len(b.chain),
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Transactions: b.pendingTransactions,
		PreviousHash: previousHash,
	}

	block.Hash = Hash(block)

	b.pendingTransactions = make([]*Transaction, 0)
	b.chain = append(b.chain, block)

	fmt.Printf("Created block %d\n", block.Index)
	return block
}

func Hash(block *Block) string {
	blockBytes, _ := json.Marshal(block)
	hasher := sha256.New()
	hasher.Write(blockBytes)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (b *Blockchain) LastBlock() *Block {
	if len(b.chain) == 0 {
		return nil
	}
	return b.chain[len(b.chain)-1]
}

func (b *Blockchain) NewTransaction(sender, recipient string, amount float64) *Transaction {
	transaction := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	b.pendingTransactions = append(b.pendingTransactions, transaction)
	return transaction
}
