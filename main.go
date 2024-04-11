package main

import (
	"github.com/illenko/blockchain/internal"
	bolt "go.etcd.io/bbolt"
	"log"
)

func main() {
	bc := internal.NewBlockchain()
	defer func(DB *bolt.DB) {
		err := DB.Close()
		if err != nil {
			log.Panic(err)
		}
	}(bc.DB)

	cli := internal.CLI{BC: bc}
	cli.Run()
}
