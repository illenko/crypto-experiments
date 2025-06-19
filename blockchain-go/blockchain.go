package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Blockchain struct {
	Chain               []*Block       `json:"chain"`
	PendingTransactions []*Transaction `json:"pendingTransactions"`
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Chain:               make([]*Block, 0),
		PendingTransactions: make([]*Transaction, 0),
	}
	fmt.Println("🌱 Creating new blockchain...")
	bc.Chain = append(bc.Chain, createGenesisBlock())
	fmt.Println("⛓️ Genesis block created!")

	return bc
}

func createGenesisBlock() *Block {
	block := &Block{
		Index:        0,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Transactions: make([]*Transaction, 0),
		PreviousHash: strings.Repeat("0", 64),
		Nonce:        0,
	}

	for {
		block.Hash = Hash(block)
		if strings.HasPrefix(block.Hash, strings.Repeat("0", DIFFICULTY)) {
			break
		}
		block.Nonce++
	}

	return block
}

func Hash(block *Block) string {
	blockBytes, _ := json.Marshal(block)
	hasher := sha256.New()
	hasher.Write(blockBytes)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (b *Blockchain) NewTransaction(sender, recipient string, amount float64) *Transaction {
	transaction := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	b.PendingTransactions = append(b.PendingTransactions, transaction)
	fmt.Printf("💸 New transaction: %s -> %s: %.2f\n", sender, recipient, amount)
	return transaction
}

func (b *Blockchain) String() string {
	bytes, _ := json.MarshalIndent(b, "", "  ")
	return string(bytes)
}

func (b *Blockchain) SubmitBlock(block *Block) error {
	if !b.isValidNewBlock(block) {
		fmt.Println("❌ Invalid block rejected")
		return fmt.Errorf("invalid block submitted")
	}

	b.Chain = append(b.Chain, block)
	b.PendingTransactions = make([]*Transaction, 0)

	fmt.Printf("✅ Block #%d accepted with hash: %s\n", block.Index, block.Hash[:8])
	return nil
}

func (b *Blockchain) isValidNewBlock(block *Block) bool {
	if !isValidProof(block) {
		return false
	}

	if block.Index != len(b.Chain) {
		return false
	}

	lastBlock := b.LastBlock()
	if lastBlock != nil && block.PreviousHash != lastBlock.Hash {
		return false
	}

	return true
}

func (b *Blockchain) LastBlock() *Block {
	if len(b.Chain) == 0 {
		return nil
	}
	return b.Chain[len(b.Chain)-1]
}
