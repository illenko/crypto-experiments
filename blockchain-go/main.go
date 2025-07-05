package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("🚀 Blockchain Go - P2P Network")
		fmt.Println("Usage:")
		fmt.Println("  --demo           Run original UTXO demo")
		fmt.Println("  --eth-demo       Run new Ethereum-like demo")
		fmt.Println("  --node           Run as blockchain node")
		fmt.Println("  --client         Run as client")
		fmt.Println("  --port <port>    Port for node (default: 8080)")
		fmt.Println("  --peers <peers>  Comma-separated peer addresses")
		fmt.Println("  --data-dir <dir> Data directory for blockchain storage (default: ./blockchain-data)")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "--demo":
		runDemo()
	case "--eth-demo":
		runEthDemo()
	case "--node":
		nodeFlags := flag.NewFlagSet("node", flag.ExitOnError)
		port := nodeFlags.Int("port", 8080, "Port for node to listen on")
		peers := nodeFlags.String("peers", "", "Comma-separated list of peer addresses")
		dataDir := nodeFlags.String("data-dir", "./blockchain-data", "Data directory for blockchain storage")

		if len(os.Args) > 2 {
			if err := nodeFlags.Parse(os.Args[2:]); err != nil {
				log.Fatal("❌ Failed to parse node flags:", err)
			}
		}
		runNode(*port, *peers, *dataDir)
	case "--client":
		runClient()
	default:
		fmt.Println("🚀 Blockchain Go - P2P Network")
		fmt.Println("Usage:")
		fmt.Println("  --demo           Run original UTXO demo")
		fmt.Println("  --eth-demo       Run new Ethereum-like demo")
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

	bc := NewBlockchain(nil) // Demo mode without persistence
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

func runNode(port int, peers string, dataDir string) {
	node := NewNode(port, peers, dataDir)
	defer node.Shutdown() // Ensure database is closed properly
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

// runEthDemo demonstrates the new Ethereum-like account-based blockchain
func runEthDemo() {
	fmt.Println("🚀 Ethereum-like Account-Based Blockchain Demo")
	fmt.Println("==============================================")

	// Create blockchain with Ethereum-like features
	bc := NewBlockchain(nil) // Demo mode without persistence

	// Initialize some accounts for testing
	alice := "0x1234567890123456789012345678901234567890"
	bob := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	miner := "0x9876543210987654321098765432109876543210"

	fmt.Printf("👤 Alice: %s\n", alice)
	fmt.Printf("👤 Bob: %s\n", bob)
	fmt.Printf("⛏️  Miner: %s\n", miner)

	// Give Alice some initial balance (100 ETH in wei)
	initialBalance := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
	bc.WorldState.AddToBalance(alice, initialBalance)

	fmt.Printf("\n💰 Alice initial balance: %s wei (100 ETH)\n",
		bc.GetEthBalance(alice).String())
	fmt.Printf("💰 Bob initial balance: %s wei\n",
		bc.GetEthBalance(bob).String())
	fmt.Printf("💰 Miner initial balance: %s wei\n",
		bc.GetEthBalance(miner).String())

	// Create a transaction from Alice to Bob
	value := new(big.Int).Mul(big.NewInt(5), big.NewInt(1e18)) // 5 ETH in wei
	gasPrice := big.NewInt(20e9)                               // 20 Gwei

	fmt.Println("\n💸 Creating transaction from Alice to Bob...")
	tx := bc.CreateEthTransaction(alice, bob, value, gasPrice)

	if tx != nil {
		fmt.Printf("✅ Transaction created successfully!\n")
		fmt.Printf("   From: %s\n", tx.From)
		fmt.Printf("   To: %s\n", tx.To)
		fmt.Printf("   Value: %s wei (5 ETH)\n", tx.Value.String())
		fmt.Printf("   Gas: %d\n", tx.Gas)
		fmt.Printf("   Gas Price: %s wei (20 Gwei)\n", tx.GasPrice.String())
		fmt.Printf("   Nonce: %d\n", tx.Nonce)
		fmt.Printf("   Fee: %s wei\n", tx.CalculateFeeEth().String())
		fmt.Printf("   Hash: %s\n", tx.Hash)

		// Execute the transaction
		fmt.Println("\n⚡ Executing transaction...")
		if err := bc.ExecuteEthTransaction(tx, miner); err != nil {
			fmt.Printf("❌ Transaction execution failed: %v\n", err)
		} else {
			fmt.Println("✅ Transaction executed successfully!")

			// Show updated balances
			fmt.Println("\n💰 Updated balances:")
			fmt.Printf("   Alice: %s wei\n", bc.GetEthBalance(alice).String())
			fmt.Printf("   Bob: %s wei\n", bc.GetEthBalance(bob).String())
			fmt.Printf("   Miner: %s wei\n", bc.GetEthBalance(miner).String())

			// Show account states
			fmt.Println("\n🏦 Account states:")
			aliceAccount := bc.WorldState.GetAccount(alice)
			bobAccount := bc.WorldState.GetAccount(bob)
			minerAccount := bc.WorldState.GetAccount(miner)

			fmt.Printf("   Alice - Balance: %s, Nonce: %d\n",
				aliceAccount.Balance.String(), aliceAccount.Nonce)
			fmt.Printf("   Bob - Balance: %s, Nonce: %d\n",
				bobAccount.Balance.String(), bobAccount.Nonce)
			fmt.Printf("   Miner - Balance: %s, Nonce: %d\n",
				minerAccount.Balance.String(), minerAccount.Nonce)

			// Show world state root
			fmt.Printf("\n🌍 World State Root: %s\n", bc.WorldState.StateRoot)
		}
	}

	// Try another transaction from Alice to Bob with higher nonce
	fmt.Println("\n💸 Creating second transaction from Alice to Bob...")
	value2 := new(big.Int).Mul(big.NewInt(3), big.NewInt(1e18)) // 3 ETH in wei
	tx2 := bc.CreateEthTransaction(alice, bob, value2, gasPrice)

	if tx2 != nil {
		fmt.Printf("✅ Second transaction created!\n")
		fmt.Printf("   Value: %s wei (3 ETH)\n", tx2.Value.String())
		fmt.Printf("   Nonce: %d (incremented)\n", tx2.Nonce)

		// Execute the second transaction
		fmt.Println("\n⚡ Executing second transaction...")
		if err := bc.ExecuteEthTransaction(tx2, miner); err != nil {
			fmt.Printf("❌ Second transaction execution failed: %v\n", err)
		} else {
			fmt.Println("✅ Second transaction executed successfully!")

			// Show final balances
			fmt.Println("\n💰 Final balances:")
			fmt.Printf("   Alice: %s wei\n", bc.GetEthBalance(alice).String())
			fmt.Printf("   Bob: %s wei\n", bc.GetEthBalance(bob).String())
			fmt.Printf("   Miner: %s wei\n", bc.GetEthBalance(miner).String())

			// Convert to ETH for readability
			aliceEth := new(big.Float).Quo(new(big.Float).SetInt(bc.GetEthBalance(alice)), big.NewFloat(1e18))
			bobEth := new(big.Float).Quo(new(big.Float).SetInt(bc.GetEthBalance(bob)), big.NewFloat(1e18))
			minerEth := new(big.Float).Quo(new(big.Float).SetInt(bc.GetEthBalance(miner)), big.NewFloat(1e18))

			fmt.Println("\n💰 Final balances in ETH:")
			fmt.Printf("   Alice: %s ETH\n", aliceEth.Text('f', 6))
			fmt.Printf("   Bob: %s ETH\n", bobEth.Text('f', 6))
			fmt.Printf("   Miner: %s ETH\n", minerEth.Text('f', 6))
		}
	}

	fmt.Println("\n🎉 Ethereum-like demo completed!")
	fmt.Println("Key differences from UTXO model:")
	fmt.Println("✓ Account-based state management")
	fmt.Println("✓ Nonce-based replay protection")
	fmt.Println("✓ Gas-based fee system")
	fmt.Println("✓ World state root for verification")
	fmt.Println("✓ Direct balance modifications")
}
