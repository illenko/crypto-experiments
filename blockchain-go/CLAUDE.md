# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

```bash
# Run the original blockchain demo
go run . --demo

# Run as a blockchain node (HTTP API server)
go run . --node --port 8080

# Run multiple nodes with peer connections
go run . --node --port 8001 --peers localhost:8002,localhost:8003
go run . --node --port 8002 --peers localhost:8001,localhost:8003
go run . --node --port 8003 --peers localhost:8001,localhost:8002

# Client commands (connect to running node)
go run . --client --cmd status
go run . --client --cmd balance --address <wallet-address>
go run . --client --cmd send --address <from> --to <to> --amount 5.0 --fee 0.1
go run . --client --cmd mine

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
- **Node** (`node.go`): P2P network node with HTTP API server, peer management, transaction/block broadcasting
- **Client** (`client.go`): HTTP client for interacting with blockchain nodes via REST API
- **Database** (`database.go`): BadgerDB persistence layer with per-node isolation and CRUD operations
- **Main** (`main.go`): Command-line interface supporting demo, node, and client modes

### Key Architecture Patterns

- **UTXO Model**: Bitcoin-style unspent transaction outputs with proper coin creation/destruction
- **Digital Signatures**: ECDSA cryptographic authentication for transaction authorization
- **Wallet System**: Key pair generation and Bitcoin-style address derivation
- **Proof-of-Work**: Uses SHA-256 hashing with nonce increment until hash has required leading zeros
- **Genesis Block**: First block created automatically with coinbase transaction distributing initial coins
- **Coinbase Transactions**: Special transactions with no inputs that create new coins for mining rewards and collect transaction fees
- **Pending Transactions**: UTXO-based transactions accumulate in blockchain until mined into a block
- **Block Validation**: Validates proof-of-work, block index, and previous hash linkage
- **UTXO Set Management**: Tracks unspent outputs by address, removes spent UTXOs and adds new ones atomically
- **Transaction Validation**: Pre-validates all transactions before applying changes to prevent corruption
- **Atomic Updates**: All-or-nothing UTXO modifications with rollback capability for consistency
- **Double-Spending Prevention**: Enforces that each UTXO can only be spent once

### Data Flow

1. Create blockchain with genesis block containing coinbase transaction that distributes initial coins
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

- ~~No network layer or peer-to-peer communication~~ ✅ **RESOLVED**
- ~~No full blockchain synchronization (longest chain consensus pending)~~ ✅ **RESOLVED**
- ~~No persistence layer (blockchain resets on restart)~~ ✅ **RESOLVED**
- Fixed mining difficulty (no difficulty adjustment)

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
- **Coinbase transactions for mining rewards**
- **Proper genesis block with coinbase transaction**
- **Transaction fees with miner collection**

- **Balance calculation from UTXOs with GetBalance method**
- **BadgerDB persistence layer** - Blocks, UTXO sets, and metadata persist across restarts
- **Per-node database isolation** - Each node maintains separate database instance
- **Graceful recovery** - UTXO set reconstruction from blockchain data

**🚧 Remaining Work:**
- Advanced wallet management and key storage

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

## P2P Network Development Plan

### **✅ Phase 1: Foundation (COMPLETED)**
1. ✅ **Restructure main.go** - Add command-line flags for node/client modes
2. ✅ **Create Node structure** - Basic node with blockchain, wallet, and peer management
3. ✅ **HTTP API server** - Basic REST endpoints for health/status
4. ✅ **Transaction APIs** - Endpoints for creating and submitting transactions
5. ✅ **Query APIs** - Balance and blockchain information endpoints
6. ✅ **Client CLI** - Command-line interface for user interactions

### **✅ Phase 2: P2P Communication (COMPLETED)**
7. ✅ **Peer discovery** - Simple peer list management and registration
8. ✅ **Transaction broadcasting** - Propagate transactions across network
9. ✅ **Mining coordination** - Prevent simultaneous mining conflicts

### **✅ Phase 3: Consensus (COMPLETED)**
10. ✅ **Blockchain sync** - Synchronize blockchain state between peers
11. ✅ **Longest chain** - Implement consensus mechanism
12. ✅ **Conflict resolution** - Handle competing blockchain versions

### **✅ Phase 4: Persistence (COMPLETED)**
13. ✅ **BadgerDB Integration** - Embedded NoSQL database for each node
14. ✅ **Per-node isolation** - Separate databases by port (`./blockchain-data/node-{port}/badger/`)
15. ✅ **Blockchain persistence** - Blocks, UTXO sets, and metadata survive restarts
16. ✅ **Graceful recovery** - UTXO set reconstruction and error handling

### **📋 Phase 5: Enhancements (OPTIONAL)**
17. **WebSocket updates** - Real-time notifications (optional)
18. **Advanced wallet management** - Persistent key storage and wallet file support

## Current P2P Implementation Status

**✅ COMPLETED FEATURES:**

**Node Architecture:**
- Multi-mode executable: `--node`, `--client`, `--demo`
- HTTP API server with comprehensive endpoints
- Thread-safe operations with mutex protection
- Peer discovery and management system

**P2P Communication:**
- Automatic peer discovery and connection
- Transaction broadcasting across network
- Block broadcasting for mining coordination
- Duplicate prevention for transactions and blocks

**API Endpoints:**
- `GET /health` - Node health check
- `GET /status` - Node status and statistics
- `GET /peers` - List connected peers
- `POST /peers` - Add new peer to network
- `GET /blockchain` - Full blockchain data
- `POST /blockchain/sync` - Synchronize with peer blockchain
- `GET /balance/<address>` - Query wallet balance
- `POST /transaction` - Create and broadcast transaction
- `POST /transaction/broadcast` - Receive broadcasted transactions
- `POST /block/broadcast` - Receive broadcasted blocks
- `POST /mine` - Mine new block and broadcast

**Client Interface:**
- `--cmd status` - Get node status
- `--cmd balance --address <addr>` - Check balance
- `--cmd send --address <from> --to <to> --amount <amount>` - Send transaction
- `--cmd mine` - Trigger mining

**Network Behavior:**
- Nodes automatically announce themselves to peers
- Transactions propagate instantly across the network
- Mined blocks are broadcast to prevent conflicts
- Basic mining coordination implemented

**Consensus Mechanism:**
- Automatic blockchain synchronization on peer connection
- Longest chain rule with full chain validation
- Atomic chain replacement with UTXO state reconstruction
- Deterministic genesis blocks ensure network compatibility

**Persistence Layer:**
- BadgerDB embedded NoSQL database for each node
- Per-node database isolation: `./blockchain-data/node-{port}/badger/`
- Automatic persistence of blocks, UTXO sets, and metadata
- Graceful recovery with UTXO set reconstruction from blockchain
- CLI flag: `--data-dir <dir>` for custom storage locations

## Next Steps for Development

### High Priority (Features & Usability)
1. **Advanced Wallet Management** - Persistent key storage and wallet file support

### Low Priority (Advanced Features)
1. **Transaction Pool Management** - Better pending transaction handling
2. **Dynamic Difficulty** - Adjust mining difficulty based on block time

The codebase follows Go conventions with structured logging using emojis for user-friendly output.