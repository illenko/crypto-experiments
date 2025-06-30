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
	genesisTx := NewGenesisTransaction()
	fmt.Printf("🎯 Creating genesis transaction: %.2f coins (unspendable)\n", genesisTx.Outputs[0].Value)

	block := &Block{
		Index:        0,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Transactions: []*Transaction{genesisTx},
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

func (b *Blockchain) CreateTransaction(fromAddress, toAddress string, amount, fee float64) *Transaction {
	senderUTXOs := b.FindUTXO(fromAddress)

	var inputs []*TxInput
	var total float64

	totalNeeded := amount + fee

	for _, utxo := range senderUTXOs {
		if total >= totalNeeded {
			break
		}
		input := &TxInput{
			TxID:      utxo.TxID,
			OutIndex:  utxo.OutIndex,
			Signature: nil,
			PubKey:    nil,
		}
		inputs = append(inputs, input)
		total += utxo.Output.Value
	}

	if total < totalNeeded {
		fmt.Printf("❌ Insufficient funds: have %.2f, need %.2f (amount: %.2f + fee: %.2f)\n", total, totalNeeded, amount, fee)
		return nil
	}

	var outputs []*TxOutput

	outputs = append(outputs, &TxOutput{
		Value:   amount,
		Address: toAddress,
	})

	if total > totalNeeded {
		outputs = append(outputs, &TxOutput{
			Value:   total - totalNeeded,
			Address: fromAddress,
		})
	}

	transaction := &Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}
	transaction.SetID()

	b.PendingTransactions = append(b.PendingTransactions, transaction)
	if fee > 0 {
		fmt.Printf("💸 New transaction: %s -> %s: %.2f + %.2f fee (ID: %s)\n", fromAddress, toAddress, amount, fee, transaction.ID[:8])
	} else {
		fmt.Printf("💸 New transaction: %s -> %s: %.2f (ID: %s)\n", fromAddress, toAddress, amount, transaction.ID[:8])
	}
	return transaction
}

func (b *Blockchain) FindUTXO(address string) []*UTXO {
	if utxos, exists := b.UTXOSet[address]; exists {
		return utxos
	}
	return []*UTXO{}
}

func (b *Blockchain) utxoExists(address, txID string, outIndex int) bool {
	utxos, exists := b.UTXOSet[address]
	if !exists {
		return false
	}

	for _, utxo := range utxos {
		if utxo.TxID == txID && utxo.OutIndex == outIndex {
			return true
		}
	}
	return false
}

func (b *Blockchain) copyUTXOSet() map[string][]*UTXO {
	copy := make(map[string][]*UTXO)
	for address, utxos := range b.UTXOSet {
		copyUTXOs := make([]*UTXO, len(utxos))
		for i, utxo := range utxos {
			copyUTXOs[i] = &UTXO{
				TxID:     utxo.TxID,
				OutIndex: utxo.OutIndex,
				Output: &TxOutput{
					Value:     utxo.Output.Value,
					Address:   utxo.Output.Address,
					ScriptPub: utxo.Output.ScriptPub,
				},
			}
		}
		copy[address] = copyUTXOs
	}
	return copy
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

func (b *Blockchain) validateTransactions(block *Block) error {
	fmt.Printf("🔍 Validating transactions for block #%d...\n", block.Index)

	if block.Index == 0 {
		if len(block.Transactions) != 1 {
			return fmt.Errorf("genesis block should have exactly one transaction")
		}
		genesisTx := block.Transactions[0]
		if !genesisTx.IsCoinbase() {
			return fmt.Errorf("genesis block transaction must be coinbase")
		}
		if genesisTx.ID != "genesis-coinbase-transaction" {
			return fmt.Errorf("invalid genesis transaction ID")
		}
		fmt.Printf("✓ Genesis block validated: genesis transaction present (unspendable)\n")
		return nil
	}

	if len(block.Transactions) == 0 {
		return fmt.Errorf("non-genesis blocks must have at least one transaction (coinbase)")
	}

	for i, tx := range block.Transactions {
		if tx.IsCoinbase() {
			if i != 0 {
				return fmt.Errorf("coinbase transaction must be the first transaction in block")
			}

			var totalFees float64
			for j := 1; j < len(block.Transactions); j++ {
				totalFees += block.Transactions[j].CalculateFee(b.UTXOSet)
			}

			expectedReward := MiningReward + totalFees
			if len(tx.Outputs) != 1 || tx.Outputs[0].Value != expectedReward {
				return fmt.Errorf("invalid coinbase transaction reward: expected %.2f, got %.2f", expectedReward, tx.Outputs[0].Value)
			}

			if totalFees > 0 {
				fmt.Printf("✓ Coinbase transaction validated: %.2f coins (%.2f reward + %.2f fees) to %s\n",
					tx.Outputs[0].Value, MiningReward, totalFees, tx.Outputs[0].Address)
			} else {
				fmt.Printf("✓ Coinbase transaction validated: %.2f coins to %s\n", tx.Outputs[0].Value, tx.Outputs[0].Address)
			}
			continue
		}

		if i == 0 {
			return fmt.Errorf("first transaction must be coinbase transaction")
		}

		var inputSum float64
		for _, input := range tx.Inputs {
			ownerAddress := ""
			for addr, addrUTXOs := range b.UTXOSet {
				for _, utxo := range addrUTXOs {
					if utxo.TxID == input.TxID && utxo.OutIndex == input.OutIndex {
						ownerAddress = addr
						inputSum += utxo.Output.Value
						break
					}
				}
				if ownerAddress != "" {
					break
				}
			}

			if ownerAddress == "" {
				return fmt.Errorf("UTXO not found: %s[%d]", input.TxID, input.OutIndex)
			}
		}

		var outputSum float64
		for _, output := range tx.Outputs {
			outputSum += output.Value
		}

		if inputSum < outputSum {
			return fmt.Errorf("transaction %s: insufficient inputs (%.2f) for outputs (%.2f)", tx.ID, inputSum, outputSum)
		}

		fee := inputSum - outputSum
		fmt.Printf("✓ Transaction %s validated: %.2f in, %.2f out, %.2f fee\n", tx.ID[:8], inputSum, outputSum, fee)
	}

	return nil
}

func (b *Blockchain) applyUTXOChanges(block *Block) error {
	for _, tx := range block.Transactions {
		if !tx.IsCoinbase() {
			for _, input := range tx.Inputs {
				ownerAddress := ""
				for addr, addrUTXOs := range b.UTXOSet {
					for _, utxo := range addrUTXOs {
						if utxo.TxID == input.TxID && utxo.OutIndex == input.OutIndex {
							ownerAddress = addr
							break
						}
					}
					if ownerAddress != "" {
						break
					}
				}
				if ownerAddress != "" {
					b.removeUTXO(ownerAddress, input.TxID, input.OutIndex)
				}
			}
		}

		if tx.IsCoinbase() && block.Index == 0 {
			fmt.Printf("🚫 Genesis coinbase ignored: %.2f coins remain unspendable\n", tx.Outputs[0].Value)
			continue
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
	return nil
}

func (b *Blockchain) SubmitBlock(block *Block) error {
	if !b.isValidNewBlock(block) {
		fmt.Println("❌ Invalid block rejected")
		return fmt.Errorf("invalid block submitted")
	}

	utxoBackup := b.copyUTXOSet()

	if err := b.validateTransactions(block); err != nil {
		fmt.Printf("❌ Transaction validation failed: %v\n", err)
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	fmt.Printf("🔄 Updating UTXO set for block #%d...\n", block.Index)

	if err := b.applyUTXOChanges(block); err != nil {
		fmt.Printf("⚠️ UTXO update failed, rolling back...\n")
		b.UTXOSet = utxoBackup
		return fmt.Errorf("UTXO update failed: %v", err)
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

func (b *Blockchain) GetBalance(address string) float64 {
	balance := 0.0
	utxos := b.FindUTXO(address)
	for _, utxo := range utxos {
		balance += utxo.Output.Value
	}
	return balance
}
