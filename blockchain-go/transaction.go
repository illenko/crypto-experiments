package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
)

const MiningReward = 10.0

type TxOutput struct {
	Value     float64 `json:"value"`
	Address   string  `json:"address"`
	ScriptPub string  `json:"scriptPub"`
}

type TxInput struct {
	TxID      string `json:"txId"`
	OutIndex  int    `json:"outIndex"`
	Signature []byte `json:"signature"`
	PubKey    []byte `json:"pubKey"`
}

type Transaction struct {
	ID      string      `json:"id"`
	Inputs  []*TxInput  `json:"inputs"`
	Outputs []*TxOutput `json:"outputs"`
}

type UTXO struct {
	TxID     string    `json:"txId"`
	OutIndex int       `json:"outIndex"`
	Output   *TxOutput `json:"output"`
}

func (tx *Transaction) SetID() {
	txBytes, _ := json.Marshal(tx)
	hash := sha256.Sum256(txBytes)
	tx.ID = hex.EncodeToString(hash[:])
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].TxID == "" && tx.Inputs[0].OutIndex == -1
}

func (tx *Transaction) String() string {
	return fmt.Sprintf("{id: %s, inputs: %d, outputs: %d}", tx.ID[:8], len(tx.Inputs), len(tx.Outputs))
}

func (utxo *UTXO) String() string {
	return fmt.Sprintf("{txId: %s, outIndex: %d, value: %.2f, address: %s}",
		utxo.TxID[:8], utxo.OutIndex, utxo.Output.Value, utxo.Output.Address)
}

func (tx *Transaction) Hash() []byte {
	txCopy := *tx

	for _, input := range txCopy.Inputs {
		input.Signature = nil
		input.PubKey = nil
	}

	txBytes, _ := json.Marshal(txCopy)
	hash := sha256.Sum256(txBytes)
	return hash[:]
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, input := range tx.Inputs {
		if prevTXs[input.TxID] == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, input := range txCopy.Inputs {
		prevTx := prevTXs[input.TxID]
		input.Signature = nil
		input.PubKey = []byte(prevTx.Outputs[input.OutIndex].ScriptPub)

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, privateKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inID].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() *Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput

	for _, input := range tx.Inputs {
		inputs = append(inputs, &TxInput{input.TxID, input.OutIndex, nil, nil})
	}

	for _, output := range tx.Outputs {
		outputs = append(outputs, &TxOutput{output.Value, output.Address, output.ScriptPub})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}
	return &txCopy
}

func (tx *Transaction) Verify(prevTXs map[string]*Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, input := range tx.Inputs {
		if prevTXs[input.TxID] == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, input := range tx.Inputs {
		prevTx := prevTXs[input.TxID]
		input := txCopy.Inputs[inID]
		input.Signature = nil
		input.PubKey = []byte(prevTx.Outputs[input.OutIndex].ScriptPub)

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r := big.Int{}
		s := big.Int{}
		sigLen := len(tx.Inputs[inID].Signature)
		r.SetBytes(tx.Inputs[inID].Signature[:(sigLen / 2)])
		s.SetBytes(tx.Inputs[inID].Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(tx.Inputs[inID].PubKey)
		x.SetBytes(tx.Inputs[inID].PubKey[:(keyLen / 2)])
		y.SetBytes(tx.Inputs[inID].PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, []byte(dataToSign), &r, &s) {
			return false
		}
	}

	return true
}

func (tx *Transaction) CalculateFee(utxoSet map[string][]*UTXO) float64 {
	if tx.IsCoinbase() {
		return 0.0
	}

	var inputSum float64
	for _, input := range tx.Inputs {
		for _, addrUTXOs := range utxoSet {
			for _, utxo := range addrUTXOs {
				if utxo.TxID == input.TxID && utxo.OutIndex == input.OutIndex {
					inputSum += utxo.Output.Value
					break
				}
			}
		}
	}

	var outputSum float64
	for _, output := range tx.Outputs {
		outputSum += output.Value
	}

	return inputSum - outputSum
}

func NewCoinbaseTransaction(toAddress string, fees float64) *Transaction {
	input := &TxInput{
		TxID:      "",
		OutIndex:  -1,
		Signature: nil,
		PubKey:    nil,
	}

	output := &TxOutput{
		Value:   MiningReward + fees,
		Address: toAddress,
	}

	tx := &Transaction{
		Inputs:  []*TxInput{input},
		Outputs: []*TxOutput{output},
	}
	tx.SetID()

	return tx
}

func NewGenesisTransaction() *Transaction {
	input := &TxInput{
		TxID:      "",
		OutIndex:  -1,
		Signature: []byte("Educational Blockchain 2025 - Genesis Block Created"),
		PubKey:    nil,
	}

	output := &TxOutput{
		Value:     50.0,
		Address:   "1GenesisBlockUnspendableAddressXXXXXXXXXXXXXX",
		ScriptPub: "Genesis Block - These coins are unspendable by design",
	}

	tx := &Transaction{
		ID:      "genesis-coinbase-transaction",
		Inputs:  []*TxInput{input},
		Outputs: []*TxOutput{output},
	}

	return tx
}

// ========================================
// NEW: Ethereum-like Account-based Model
// ========================================

// Account represents an Ethereum-like account with balance, nonce, and contract data
type Account struct {
	Balance     *big.Int `json:"balance"`     // Account balance in wei (smallest unit)
	Nonce       uint64   `json:"nonce"`       // Transaction counter for replay protection
	CodeHash    string   `json:"codeHash"`    // Hash of contract code (empty for EOA)
	StorageRoot string   `json:"storageRoot"` // Root of contract storage trie (empty for EOA)
}

// EthTransaction represents an Ethereum-like transaction
type EthTransaction struct {
	From     string   `json:"from"`     // Sender address
	To       string   `json:"to"`       // Recipient address (empty for contract creation)
	Value    *big.Int `json:"value"`    // Amount to transfer in wei
	Gas      uint64   `json:"gas"`      // Gas limit
	GasPrice *big.Int `json:"gasPrice"` // Gas price in wei
	Nonce    uint64   `json:"nonce"`    // Sender's transaction nonce
	Data     []byte   `json:"data"`     // Contract bytecode or call data
	Hash     string   `json:"hash"`     // Transaction hash

	// Signature fields (simplified for now)
	V *big.Int `json:"v"` // Recovery ID
	R *big.Int `json:"r"` // Signature R
	S *big.Int `json:"s"` // Signature S
}

// NewEOA creates a new Externally Owned Account with initial balance
func NewEOA(initialBalance *big.Int) *Account {
	return &Account{
		Balance:     new(big.Int).Set(initialBalance),
		Nonce:       0,
		CodeHash:    "", // Empty for EOA
		StorageRoot: "", // Empty for EOA
	}
}

// IsContract checks if account is a smart contract
func (a *Account) IsContract() bool {
	return a.CodeHash != ""
}

// GetAccountHash calculates hash of account for state root calculation
func (a *Account) GetAccountHash() string {
	accountBytes, _ := json.Marshal(a)
	hash := sha256.Sum256(accountBytes)
	return hex.EncodeToString(hash[:])
}

// CalculateEthHash calculates transaction hash for Ethereum-like transaction
func (tx *EthTransaction) CalculateEthHash() string {
	// Create copy without hash and signature for hashing
	txCopy := EthTransaction{
		From:     tx.From,
		To:       tx.To,
		Value:    new(big.Int).Set(tx.Value),
		Gas:      tx.Gas,
		GasPrice: new(big.Int).Set(tx.GasPrice),
		Nonce:    tx.Nonce,
		Data:     tx.Data,
	}

	txBytes, _ := json.Marshal(txCopy)
	hash := sha256.Sum256(txBytes)
	return hex.EncodeToString(hash[:])
}

// IsContractCreation checks if transaction creates a new contract
func (tx *EthTransaction) IsContractCreation() bool {
	return tx.To == ""
}

// CalculateFeeEth calculates transaction fee (gas * gasPrice)
func (tx *EthTransaction) CalculateFeeEth() *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas), tx.GasPrice)
}

// SetEthID sets the transaction hash
func (tx *EthTransaction) SetEthID() {
	tx.Hash = tx.CalculateEthHash()
}
