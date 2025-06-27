# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

```bash
# Run the blockchain demo
go run .

# Build the executable
go build .

# Format code
go fmt ./...

# Run with race detection (for concurrent code)
go run -race .

# Test (when tests are added)
go test ./...
```

## Architecture Overview

This is a basic blockchain implementation in Go consisting of several core components:

### Core Components

- **Block** (`block.go`): Basic block structure containing index, timestamp, transactions, previous hash, current hash, and nonce for proof-of-work
- **Blockchain** (`blockchain.go`): Main blockchain structure managing the chain of blocks and pending transactions, includes validation logic
- **Transaction** (`transaction.go`): Simple transaction structure with sender, recipient, and amount
- **Mining** (`mining.go`): Proof-of-work mining implementation with configurable difficulty (currently set to 4 leading zeros)

### Key Architecture Patterns

- **Proof-of-Work**: Uses SHA-256 hashing with nonce increment until hash has required leading zeros
- **Genesis Block**: First block created automatically with all-zero previous hash
- **Pending Transactions**: Transactions accumulate in blockchain until mined into a block
- **Block Validation**: Validates proof-of-work, block index, and previous hash linkage

### Data Flow

1. Create blockchain with genesis block
2. Add transactions to pending pool via `NewTransaction()`
3. Miner calls `Mine()` to create candidate block with pending transactions
4. Mining loop increments nonce until valid hash found
5. Validated block added to chain via `SubmitBlock()`
6. Pending transactions cleared after successful block addition

### Current Limitations

- No network layer or peer-to-peer communication
- Simple transaction model without digital signatures or UTXO
- Fixed mining difficulty
- Single miner implementation
- No persistence layer

The codebase follows Go conventions with structured logging using emojis for user-friendly output.