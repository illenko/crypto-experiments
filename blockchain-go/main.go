package main

import "fmt"

func main() {
	wallets := NewWallets()

	aliceWallet := wallets.CreateWallet()
	bobWallet := wallets.CreateWallet()
	minerWallet := wallets.CreateWallet()

	fmt.Printf("👤 Alice: %s\n", aliceWallet)
	fmt.Printf("👤 Bob: %s\n", bobWallet)
	fmt.Printf("⛏️  Miner: %s\n", minerWallet)

	bc := NewBlockchain()

	bc.UTXOSet[aliceWallet] = []*UTXO{
		{
			TxID:     "genesis",
			OutIndex: 0,
			Output: &TxOutput{
				Value:     100.0,
				Address:   aliceWallet,
				ScriptPub: string(wallets.GetWallet(aliceWallet).PublicKey),
			},
		},
	}

	miner := NewMiner(minerWallet)

	tx := bc.CreateTransaction(aliceWallet, bobWallet, 50.0)
	if tx != nil {
		wallet := wallets.GetWallet(aliceWallet)
		prevTXs := make(map[string]*Transaction)

		for _, input := range tx.Inputs {
			prevTXs[input.TxID] = &Transaction{
				ID: input.TxID,
				Outputs: []*TxOutput{
					{
						Value:     100.0,
						Address:   aliceWallet,
						ScriptPub: string(wallet.PublicKey),
					},
				},
			}
		}

		tx.Sign(wallet.PrivateKey, prevTXs)

		for _, input := range tx.Inputs {
			input.PubKey = wallet.PublicKey
		}

		newBlock := miner.Mine(bc)
		bc.SubmitBlock(newBlock)
	}

	fmt.Println("\n🔗 Blockchain:")
	fmt.Println(bc)
}
