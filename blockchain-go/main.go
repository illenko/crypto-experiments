package main

import "fmt"

func main() {
	bc := NewBlockchain()
	bc.NewTransaction("Alice", "Bob", 50.0)
	bc.NewBlock(bc.LastBlock().Hash)

	fmt.Println(bc)
}
