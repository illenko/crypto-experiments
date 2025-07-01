package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("🚀 Blockchain Go - P2P Network")
		fmt.Println("Usage:")
		fmt.Println("  --demo           Run original demo")
		fmt.Println("  --node           Run as blockchain node")
		fmt.Println("  --client         Run as client")
		fmt.Println("  --port <port>    Port for node (default: 8080)")
		fmt.Println("  --peers <peers>  Comma-separated peer addresses")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "--demo":
		runDemo()
	case "--node":
		nodeFlags := flag.NewFlagSet("node", flag.ExitOnError)
		port := nodeFlags.Int("port", 8080, "Port for node to listen on")
		peers := nodeFlags.String("peers", "", "Comma-separated list of peer addresses")

		if len(os.Args) > 2 {
			if err := nodeFlags.Parse(os.Args[2:]); err != nil {
				log.Fatal("❌ Failed to parse node flags:", err)
			}
		}
		runNode(*port, *peers)
	case "--client":
		runClient()
	default:
		fmt.Println("🚀 Blockchain Go - P2P Network")
		fmt.Println("Usage:")
		fmt.Println("  --demo           Run original demo")
		fmt.Println("  --node           Run as blockchain node")
		fmt.Println("  --client         Run as client")
		fmt.Println("  --port <port>    Port for node (default: 8080)")
		fmt.Println("  --peers <peers>  Comma-separated peer addresses")
		os.Exit(1)
	}
}

func runDemo() {
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

func runNode(port int, peers string) {
	node := NewNode(port, peers)
	node.Start()
}

func runClient() {
	clientFlags := flag.NewFlagSet("client", flag.ExitOnError)

	nodeAddr := clientFlags.String("node-addr", "localhost:8080", "Address of the blockchain node")
	command := clientFlags.String("cmd", "", "Command to execute: balance, send, mine, status")
	address := clientFlags.String("address", "", "Wallet address")
	to := clientFlags.String("to", "", "Recipient address")
	amount := clientFlags.Float64("amount", 0, "Amount to send")
	fee := clientFlags.Float64("fee", 0.01, "Transaction fee")

	// Parse remaining arguments after --client
	if len(os.Args) > 2 {
		if err := clientFlags.Parse(os.Args[2:]); err != nil {
			log.Fatal("❌ Failed to parse client flags:", err)
		}
	}

	fmt.Printf("💻 Blockchain Client - connecting to %s\n", *nodeAddr)

	if *command == "" {
		fmt.Println("Available commands:")
		fmt.Println("  --cmd balance --address <addr>         Check wallet balance")
		fmt.Println("  --cmd send --address <from> --to <to> --amount <amount> [--fee <fee>]  Send transaction")
		fmt.Println("  --cmd mine                             Trigger mining")
		fmt.Println("  --cmd status                           Get node status")
		fmt.Println("  --node-addr <addr>                     Specify node address (default: localhost:8080)")
		return
	}

	client := NewClient(*nodeAddr)

	switch *command {
	case "balance":
		if *address == "" {
			fmt.Println("❌ Address required for balance command")
			return
		}
		client.GetBalance(*address)
	case "send":
		if *address == "" || *to == "" || *amount <= 0 {
			fmt.Println("❌ From address, to address, and amount required for send command")
			return
		}
		client.SendTransaction(*address, *to, *amount, *fee)
	case "mine":
		client.Mine()
	case "status":
		client.GetStatus()
	default:
		fmt.Printf("❌ Unknown command: %s\n", *command)
	}
}
