package main

import "encoding/json"

type Block struct {
	Index        int            `json:"index"`
	Timestamp    string         `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
	PreviousHash string         `json:"previousHash"`
	Hash         string         `json:"hash"`
	Nonce        int            `json:"nonce"`
}

func (b *Block) String() string {
	bytes, _ := json.MarshalIndent(b, "", "  ")
	return string(bytes)
}
