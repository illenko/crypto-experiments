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
