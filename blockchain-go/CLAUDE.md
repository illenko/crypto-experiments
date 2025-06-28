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
- **Blockchain** (`blockchain.go`): Main blockchain structure managing the chain of blocks, pending transactions, and UTXO set. Includes validation logic and UTXO management
- **Transaction** (`transaction.go`): UTXO-based transaction structure with inputs (TxInput), outputs (TxOutput), and transaction ID. Includes digital signature support and transaction validation
- **Mining** (`mining.go`): Proof-of-work mining implementation with configurable difficulty (currently set to 4 leading zeros)
- **Wallet** (`wallet.go`): ECDSA key pair generation and Bitcoin-style address derivation with Base58 encoding
- **Base58** (`base58.go`): Bitcoin-compatible Base58 encoding/decoding for human-readable addresses

### Key Architecture Patterns

- **UTXO Model**: Bitcoin-style unspent transaction outputs with proper coin creation/destruction
- **Digital Signatures**: ECDSA cryptographic authentication for transaction authorization
- **Wallet System**: Key pair generation and Bitcoin-style address derivation
- **Proof-of-Work**: Uses SHA-256 hashing with nonce increment until hash has required leading zeros
- **Genesis Block**: First block created automatically with all-zero previous hash
- **Pending Transactions**: UTXO-based transactions accumulate in blockchain until mined into a block
- **Block Validation**: Validates proof-of-work, block index, and previous hash linkage
- **UTXO Set Management**: Tracks unspent outputs by address, removes spent UTXOs and adds new ones atomically
- **Transaction Validation**: Pre-validates all transactions before applying changes to prevent corruption
- **Atomic Updates**: All-or-nothing UTXO modifications with rollback capability for consistency
- **Double-Spending Prevention**: Enforces that each UTXO can only be spent once

### Data Flow

1. Create blockchain with genesis block and initialize empty UTXO set
2. Create UTXO-based transactions via `CreateTransaction()` which:
   - Finds sender's available UTXOs
   - Validates sufficient funds
   - Creates transaction inputs (spending UTXOs) and outputs (new UTXOs)
   - Handles change back to sender
3. Sign transactions with wallet private keys via `Sign()` method
4. Miner calls `Mine()` to create candidate block with pending transactions
5. Mining loop increments nonce until valid hash found (proof-of-work)
6. Validated block added to chain via `SubmitBlock()` which:
   - Pre-validates all transactions (UTXO existence, input/output balance, signatures)
   - Creates backup of UTXO set for potential rollback
   - Atomically removes spent UTXOs and adds new ones
   - Rolls back changes if any step fails
   - Clears pending transactions only on success

### Current Limitations

- No network layer or peer-to-peer communication
- Fixed mining difficulty (no difficulty adjustment)
- Single miner implementation
- No persistence layer (blockchain resets on restart)
- Manual UTXO initialization in main.go (no proper genesis transaction)
- No coinbase transactions for mining rewards

### UTXO Implementation Status

**✅ Completed:**
- UTXO data structures (TxInput, TxOutput, UTXO)
- UTXO set management and tracking
- Proper coin destruction and creation
- Change handling in transactions
- Double-spending prevention
- **Transaction validation with consistency guarantees**
- **Atomic UTXO updates with rollback**
- **Pre-validation of all transactions before state changes**
- **Conservation of money enforcement (input sum >= output sum)**
- **Digital signatures for transaction authentication**
- **ECDSA key pair generation and wallet system**
- **Bitcoin-style address derivation with Base58 encoding**
- **Transaction signing and verification**

**🚧 Remaining Work:**
- Coinbase transactions for mining rewards
- Balance calculation from UTXOs
- Transaction fees

### UTXO Consistency Guarantees

**🔒 Security Features Implemented:**
- **Double-Spending Prevention**: Each UTXO can only be spent once
- **Digital Signature Authentication**: Only private key holders can spend UTXOs
- **Cryptographic Address Derivation**: Addresses derived from public key hashes
- **Atomic Updates**: All-or-nothing modifications to UTXO set
- **Pre-validation**: All transactions validated before any state changes
- **Rollback Protection**: Failed validations don't corrupt UTXO state
- **Conservation Laws**: Input values must be >= output values
- **Deep Copy Backup**: Complete UTXO set backup for rollback scenarios

**📊 Validation Process:**
1. `validateTransactions()` - Check UTXO existence and money conservation
2. `copyUTXOSet()` - Create backup for potential rollback
3. `applyUTXOChanges()` - Atomically modify UTXO set
4. Rollback to backup if any step fails

## Next Steps for Development

### High Priority (Security & Core Functionality)
1. **Create Coinbase Transactions** - Add mining rewards:
   - Special transaction type with no inputs (creates new coins)
   - Miner address receives block reward
   - Integrate with mining process

### Medium Priority (Features & Usability)
1. **Balance Calculation** - Replace manual UTXO management with:
   - `GetBalance(address)` method that sums UTXOs
   - Better transaction creation with automatic UTXO selection

2. **Transaction Fees** - Implement fee mechanism:
   - Calculate fees as input-output difference
   - Allocate fees to miners via coinbase transaction

3. **Proper Genesis Block** - Replace manual UTXO initialization:
   - Create genesis transaction in genesis block
   - Distribute initial coins through proper transaction outputs

### Low Priority (Advanced Features)
1. **Transaction Pool Management** - Better pending transaction handling
2. **Network Layer** - P2P communication between nodes
3. **Persistence** - Save/load blockchain state from disk
4. **Dynamic Difficulty** - Adjust mining difficulty based on block time

The codebase follows Go conventions with structured logging using emojis for user-friendly output.