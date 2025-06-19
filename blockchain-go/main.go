package main

import "fmt"

func main() {
	bc := NewBlockchain()

	miner := NewMiner("miner1")

	bc.NewTransaction("Alice", "Bob", 50.0)

	newBlock := miner.Mine(bc)
	bc.Chain = append(bc.Chain, newBlock)
	bc.PendingTransactions = make([]*Transaction, 0)

	fmt.Println(bc)
}
