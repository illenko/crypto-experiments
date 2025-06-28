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
	miner := NewMiner(minerWallet)

	fmt.Println("\n💎 Mining first block to create spendable coins...")
	block1 := miner.Mine(bc)
	bc.SubmitBlock(block1)

	fmt.Printf("\n💰 Miner now has %.2f spendable coins from mining\n", MiningReward)

	fmt.Println("\n💸 Creating transaction from miner to Alice...")
	tx := bc.CreateTransaction(minerWallet, aliceWallet, 5.0, 0.1)
	if tx != nil {
		wallet := wallets.GetWallet(minerWallet)
		prevTXs := make(map[string]*Transaction)

		for _, input := range tx.Inputs {
			prevTXs[input.TxID] = &Transaction{
				ID: input.TxID,
				Outputs: []*TxOutput{
					{
						Value:     MiningReward,
						Address:   minerWallet,
						ScriptPub: string(wallet.PublicKey),
					},
				},
			}
		}

		tx.Sign(wallet.PrivateKey, prevTXs)

		for _, input := range tx.Inputs {
			input.PubKey = wallet.PublicKey
		}

		fmt.Println("\n💎 Mining second block with the transaction...")
		block2 := miner.Mine(bc)
		bc.SubmitBlock(block2)
	}

	fmt.Println("\n💸 Now Alice can send to Bob...")
	tx2 := bc.CreateTransaction(aliceWallet, bobWallet, 2.0, 0.05)
	if tx2 != nil {
		wallet := wallets.GetWallet(aliceWallet)
		prevTXs := make(map[string]*Transaction)

		for _, input := range tx2.Inputs {
			prevTXs[input.TxID] = &Transaction{
				ID: input.TxID,
				Outputs: []*TxOutput{
					{
						Value:     5.0,
						Address:   aliceWallet,
						ScriptPub: string(wallet.PublicKey),
					},
				},
			}
		}

		tx2.Sign(wallet.PrivateKey, prevTXs)

		for _, input := range tx2.Inputs {
			input.PubKey = wallet.PublicKey
		}

		fmt.Println("\n💎 Mining third block with Alice's transaction...")
		block3 := miner.Mine(bc)
		bc.SubmitBlock(block3)
	}

	fmt.Println("\n🔗 Final Blockchain:")
	fmt.Println(bc)
}
