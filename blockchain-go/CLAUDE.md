# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

```bash
# Run the original UTXO blockchain demo
go run . --demo

# Run the new Ethereum-like account-based demo
go run . --eth-demo

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

This blockchain implementation is being transformed from a Bitcoin-like UTXO model to an Ethereum-like account-based system with smart contract capabilities.

### Current Core Components (Being Transformed)

- **Block** (`block.go`): Block structure - will be enhanced with state root, gas limit, gas used
- **Blockchain** (`blockchain.go`): Main blockchain structure - transitioning from UTXO to account-based state management
- **Transaction** (`transaction.go`): Transaction structure - being replaced with account-based model (from/to/value/gas)
- **Mining** (`mining.go`): Proof-of-work mining - will be updated for gas fee collection
- **Wallet** (`wallet.go`): Key pair generation - transitioning to Ethereum-style addresses
- **Node** (`node.go`): P2P network node with HTTP API server
- **Client** (`client.go`): HTTP client for blockchain interaction
- **Database** (`database.go`): BadgerDB persistence layer
- **Main** (`main.go`): Command-line interface

### Target Architecture (Ethereum-like)

- **Account-Based State**: Replace UTXO with account balances and state
- **Smart Contracts**: Virtual machine for contract execution
- **Gas System**: Computational resource management
- **Transaction Model**: From/To/Value/Gas structure
- **State Management**: World state with state root
- **Virtual Machine**: Bytecode execution engine

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
- Atomic chain replacement with state reconstruction
- Deterministic genesis blocks ensure network compatibility

**Persistence Layer:**
- BadgerDB embedded NoSQL database for each node
- Per-node database isolation: `./blockchain-data/node-{port}/badger/`
- Automatic persistence of blocks, state, and metadata
- Graceful recovery with state reconstruction from blockchain
- CLI flag: `--data-dir <dir>` for custom storage locations

The codebase follows Go conventions with structured logging using emojis for user-friendly output.

## Ethereum-Like Implementation Plan

### **🎯 Branch: ethereum-like-implementation**

This plan outlines the transformation from the current Bitcoin-like UTXO blockchain to an Ethereum-like account-based system with smart contract capabilities.

### **✅ Phase 1: Account-Based Transaction Model (COMPLETED)**
**Goal**: Replace UTXO model with account-based state management

**Tasks:**
1. **✅ Account State Management**
   - ✅ Create `Account` struct with `balance`, `nonce`, `codeHash`, `storageRoot`
   - ✅ Add `WorldState` with `AccountState` map[string]*Account alongside existing UTXOSet
   - ✅ Implement account creation and balance tracking
   - ✅ Add nonce-based replay protection

2. **✅ Transaction Structure Updates**
   - ✅ Create `EthTransaction` struct with account-based fields:
     - ✅ `from`, `to`, `value`, `gas`, `gasPrice`, `nonce`
     - ✅ `data` field for contract calls/deployment
   - ✅ Update transaction validation for account model
   - ✅ Implement account balance checks and nonce validation

3. **✅ State Management**
   - ✅ Design world state structure with deterministic state root
   - ✅ Implement state root calculation for blocks
   - ✅ Add state transition functions (transfer, balance updates)
   - ✅ Create account-based transaction processing

**✅ Implementation Status:**
- ✅ **Account struct** (`transaction.go`) - Balance, nonce, contract fields
- ✅ **EthTransaction struct** (`transaction.go`) - Ethereum-like transaction model
- ✅ **WorldState management** (`blockchain.go`) - Account state tracking
- ✅ **Transaction validation** - Nonce checking, balance verification, gas validation
- ✅ **Transaction execution** - Gas fee collection, balance transfers, nonce increment
- ✅ **Demo implementation** (`main.go`) - `--eth-demo` command for testing
- ✅ **Gradual integration** - New account-based system alongside existing UTXO system

**🧪 Testing Results:**
- ✅ Account creation and balance management
- ✅ Nonce-based replay protection (0 → 1 → 2)
- ✅ Gas fee calculation and collection (420,000 wei per transaction)
- ✅ State root updates with each transaction
- ✅ Ethereum-style addresses (0x format)
- ✅ Wei-based value handling (100 ETH = 100,000,000,000,000,000,000 wei)

### **Phase 2: Simple Virtual Machine (Week 3-4)**
**Goal**: Build basic execution environment for smart contracts

**Tasks:**
1. **Basic VM Architecture**
   - Create `SimpleVM` struct with stack, memory, storage
   - Implement execution context with gas tracking
   - Add program counter and instruction pointer
   - Design simple instruction set

2. **Core Opcodes Implementation**
   - **Arithmetic**: ADD, SUB, MUL, DIV, MOD
   - **Logic**: AND, OR, XOR, NOT, LT, GT, EQ
   - **Memory**: LOAD, STORE, MLOAD, MSTORE
   - **Stack**: PUSH, POP, DUP, SWAP
   - **Control**: JUMP, JUMPI, PC, STOP
   - **Environment**: ADDRESS, BALANCE, CALLER, CALLVALUE

3. **Bytecode Execution Engine**
   - Implement opcode dispatcher
   - Add gas metering for each operation
   - Handle execution errors and reverts
   - Create execution result handling

### **Phase 3: Smart Contract Support (Week 5-6)**
**Goal**: Enable contract deployment and execution

**Tasks:**
1. **Contract Deployment**
   - Support contract creation transactions (to = nil)
   - Store contract bytecode in account state
   - Implement contract address generation (CREATE opcode)
   - Add constructor execution

2. **Contract Execution**
   - Execute contract code on CALL transactions
   - Handle contract-to-contract calls
   - Implement CALL, DELEGATECALL, STATICCALL opcodes
   - Add return values and revert mechanisms

3. **Contract Storage**
   - Implement persistent contract storage
   - Add SLOAD, SSTORE opcodes
   - Handle storage state changes
   - Add storage root calculation

### **Phase 4: Gas System (Week 7)**
**Goal**: Implement computational resource management

**Tasks:**
1. **Gas Mechanism**
   - Define gas costs for different operations
   - Implement gas limit and gas used tracking
   - Add gas price mechanism for transaction fees
   - Handle out-of-gas scenarios

2. **Fee System Updates**
   - Replace mining rewards with gas fees
   - Implement fee collection for miners
   - Add gas estimation for transactions
   - Handle gas refunds for storage cleanup

### **Phase 5: Enhanced Block Structure (Week 8)**
**Goal**: Update blockchain structure for account-based model

**Tasks:**
1. **Block Structure Updates**
   - Add `stateRoot`, `gasLimit`, `gasUsed` to blocks
   - Implement transaction receipts
   - Add logs and events system
   - Update block validation logic

2. **Enhanced Features**
   - Add event logging and filtering
   - Implement transaction receipts
   - Add basic debugging tools
   - Create development utilities

### **Key Architecture Changes Summary**

**From Bitcoin-like (UTXO) to Ethereum-like (Account-based):**

| Component | Bitcoin-like (Current) | Ethereum-like (Target) |
|-----------|----------------------|----------------------|
| **State Model** | UTXO Set | Account State |
| **Transactions** | Inputs/Outputs | From/To/Value/Gas |
| **Execution** | Script validation | Virtual Machine |
| **Smart Contracts** | None | Bytecode execution |
| **Fees** | Simple fee | Gas system |
| **Addresses** | Bitcoin-style | Ethereum-style |
| **Block Structure** | Simple | StateRoot + Gas tracking |

### **Implementation Guidelines**

1. **Complete Transformation**: Replace Bitcoin-like UTXO model entirely with Ethereum-like account-based system
2. **Incremental Development**: Each phase builds on previous phases
3. **Testing Strategy**: Add comprehensive tests for each component
4. **Documentation**: Update documentation as features are implemented
5. **Performance**: Consider performance implications of state management

### **Success Metrics**

- [ ] Deploy and execute simple smart contracts
- [x] Handle account-based transactions
- [x] Implement gas metering and fees
- [ ] Maintain blockchain consensus
- [ ] Support contract-to-contract calls
- [ ] Generate transaction receipts and logs

### **Phase 1 Achievements**

**✅ Completed Features:**
- [x] **Account-based state model** - Replaced UTXO with persistent account balances
- [x] **Ethereum-style transactions** - From/To/Value/Gas structure implemented
- [x] **Nonce-based security** - Replay attack prevention with account nonces
- [x] **Gas fee system** - Basic gas calculation and miner fee collection
- [x] **World state management** - Deterministic state root calculation
- [x] **Wei precision** - Full 18-decimal place precision for values
- [x] **Dual system support** - Both UTXO and account models working simultaneously

**🎯 Next Phase Ready:** The foundation is complete for implementing a simple virtual machine and smart contract execution.

This plan transforms the blockchain from a simple cryptocurrency to a programmable platform while maintaining the core networking and consensus mechanisms.