package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type TxOutput struct {
	Value     float64 `json:"value"`
	Address   string  `json:"address"`
	ScriptPub string  `json:"scriptPub"`
}

type TxInput struct {
	TxID      string `json:"txId"`
	OutIndex  int    `json:"outIndex"`
	Signature string `json:"signature"`
	PubKey    string `json:"pubKey"`
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
