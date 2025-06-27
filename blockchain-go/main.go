package main

import "fmt"

func main() {
	bc := NewBlockchain()

	bc.UTXOSet["Alice"] = []*UTXO{
		{
			TxID:     "genesis",
			OutIndex: 0,
			Output: &TxOutput{
				Value:   100.0,
				Address: "Alice",
			},
		},
	}

	miner := NewMiner("miner1")

	tx := bc.CreateTransaction("Alice", "Bob", 50.0)
	if tx != nil {
		newBlock := miner.Mine(bc)
		bc.SubmitBlock(newBlock)
	}

	fmt.Println(bc)
}
