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
	Chain               []*Block           `json:"chain"`
	PendingTransactions []*Transaction     `json:"pendingTransactions"`
	UTXOSet             map[string][]*UTXO `json:"utxoSet"`
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Chain:               make([]*Block, 0),
		PendingTransactions: make([]*Transaction, 0),
		UTXOSet:             make(map[string][]*UTXO),
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

func (b *Blockchain) CreateTransaction(fromAddress, toAddress string, amount float64) *Transaction {
	senderUTXOs := b.FindUTXO(fromAddress)

	var inputs []*TxInput
	var total float64

	for _, utxo := range senderUTXOs {
		if total >= amount {
			break
		}
		input := &TxInput{
			TxID:      utxo.TxID,
			OutIndex:  utxo.OutIndex,
			Signature: "",
			PubKey:    fromAddress,
		}
		inputs = append(inputs, input)
		total += utxo.Output.Value
	}

	if total < amount {
		fmt.Printf("❌ Insufficient funds: have %.2f, need %.2f\n", total, amount)
		return nil
	}

	var outputs []*TxOutput

	outputs = append(outputs, &TxOutput{
		Value:   amount,
		Address: toAddress,
	})

	if total > amount {
		outputs = append(outputs, &TxOutput{
			Value:   total - amount,
			Address: fromAddress,
		})
	}

	transaction := &Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}
	transaction.SetID()

	b.PendingTransactions = append(b.PendingTransactions, transaction)
	fmt.Printf("💸 New transaction: %s -> %s: %.2f (ID: %s)\n", fromAddress, toAddress, amount, transaction.ID[:8])
	return transaction
}

func (b *Blockchain) FindUTXO(address string) []*UTXO {
	if utxos, exists := b.UTXOSet[address]; exists {
		return utxos
	}
	return []*UTXO{}
}

func (b *Blockchain) removeUTXO(address, txID string, outIndex int) {
	utxos, exists := b.UTXOSet[address]
	if !exists {
		return
	}

	for i, utxo := range utxos {
		if utxo.TxID == txID && utxo.OutIndex == outIndex {
			b.UTXOSet[address] = append(utxos[:i], utxos[i+1:]...)
			txIDShort := txID
			if len(txID) > 8 {
				txIDShort = txID[:8]
			}
			fmt.Printf("🗑️ Removed UTXO: %s[%d] from %s\n", txIDShort, outIndex, address)
			break
		}
	}

	if len(b.UTXOSet[address]) == 0 {
		delete(b.UTXOSet, address)
	}
}

func (b *Blockchain) addUTXO(address string, utxo *UTXO) {
	b.UTXOSet[address] = append(b.UTXOSet[address], utxo)
	txIDShort := utxo.TxID
	if len(utxo.TxID) > 8 {
		txIDShort = utxo.TxID[:8]
	}
	fmt.Printf("➕ Added UTXO: %s[%d] to %s (%.2f)\n",
		txIDShort, utxo.OutIndex, address, utxo.Output.Value)
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

	fmt.Printf("🔄 Updating UTXO set for block #%d...\n", block.Index)

	for _, tx := range block.Transactions {
		if !tx.IsCoinbase() {
			for _, input := range tx.Inputs {
				b.removeUTXO(input.PubKey, input.TxID, input.OutIndex)
			}
		}

		for i, output := range tx.Outputs {
			utxo := &UTXO{
				TxID:     tx.ID,
				OutIndex: i,
				Output:   output,
			}
			b.addUTXO(output.Address, utxo)
		}
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
