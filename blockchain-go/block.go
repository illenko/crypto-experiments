package main

type Block struct {
	Index        int            `json:"index"`
	Timestamp    string         `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
	PreviousHash string         `json:"previous_hash"`
	Hash         string         `json:"hash"`
}
