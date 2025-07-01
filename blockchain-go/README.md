# Blockchain Go - Educational Cryptocurrency Implementation

A fully functional blockchain implementation in Go featuring UTXO-based transactions, proof-of-work mining, digital signatures, and peer-to-peer networking. This project demonstrates core blockchain concepts with a Bitcoin-inspired architecture.

## 🚀 Quick Start

### Prerequisites
- Go 1.19+ installed
- Basic understanding of blockchain concepts

### Installation
```bash
git clone <repository-url>
cd blockchain-go
go build .
```

## 🏗️ Architecture Overview

This implementation follows Bitcoin's core design patterns:

- **UTXO Model**: Bitcoin-style unspent transaction outputs for state management
- **Proof-of-Work**: SHA-256 based mining with configurable difficulty
- **Digital Signatures**: ECDSA cryptographic authentication
- **P2P Network**: HTTP-based node communication and synchronization
- **Wallet System**: Key pair generation with Base58 address encoding

## 🎯 Testing Scenarios

### Scenario 1: Single Node Operation

#### 1. Start a Blockchain Node
```bash
go run . --node --port 8080
```

**What happens internally:**
- Creates new blockchain with genesis block containing unspendable coinbase transaction
- Generates miner wallet with ECDSA key pair
- Starts HTTP API server for client interactions
- Initializes empty UTXO set for tracking unspent outputs

#### 2. Check Node Status
```bash
go run . --client --cmd status
```

**Theoretical Background:** Blockchain nodes maintain consensus state including current block height, pending transactions, and network peers.

#### 3. Mine Your First Block
```bash
go run . --client --cmd mine
```

**What happens internally:**
1. **Coinbase Transaction Creation**: Miner creates special transaction with no inputs, generating 10.0 new coins
2. **Proof-of-Work**: Node increments nonce until hash starts with 4 zeros (difficulty = 4)
3. **Block Validation**: Validates proof-of-work and transaction structure
4. **UTXO Set Update**: Atomically adds new UTXO to miner's address
5. **Blockchain Extension**: Appends valid block to chain

**Theoretical Background:** Mining implements Nakamoto consensus through computational proof-of-work, ensuring network security and creating new currency supply through coinbase transactions.

#### 4. Check Miner Balance
```bash
go run . --client --cmd balance --address <miner-wallet-address>
```

**What happens internally:**
- **UTXO Lookup**: Scans UTXO set for all unspent outputs belonging to the address
- **Balance Calculation**: Sums values of all UTXOs using `GetBalance()` method
- Returns total spendable balance

**Theoretical Background:** Bitcoin uses the UTXO model instead of account balances. Your "balance" is the sum of all unspent transaction outputs you can spend with your private key.

### Scenario 2: Transaction Lifecycle

#### 1. Send Transaction
```bash
go run . --client --cmd send --address <from-address> --to <to-address> --amount 3.0 --fee 0.1
```

**What happens internally:**
1. **UTXO Selection**: Finds sufficient UTXOs from sender's address (3.1 total needed)
2. **Transaction Construction**: 
   - Creates transaction inputs referencing selected UTXOs
   - Creates output for recipient (3.0 coins)
   - Creates change output back to sender (if applicable)
3. **Digital Signing**: Signs transaction with sender's private key (currently simulated)
4. **Validation**: Ensures inputs ≥ outputs (conservation of money)
5. **Broadcast**: Adds to pending transaction pool and broadcasts to peers

**Theoretical Background:** UTXO transactions are atomic units that consume existing outputs and create new ones. Digital signatures prove ownership without revealing private keys.

#### 2. Mine Block with Transaction
```bash
go run . --client --cmd mine
```

**What happens internally:**
1. **Block Assembly**: Includes coinbase transaction (10.1 coins: 10.0 reward + 0.1 fee) plus pending transactions
2. **Pre-validation**: Validates all transactions before applying changes
3. **UTXO Updates**: 
   - Removes spent UTXOs from sender
   - Creates new UTXOs for recipient and miner
   - Updates performed atomically with rollback protection
4. **Transaction Finality**: Moves transactions from pending to confirmed state

#### 3. Verify Balances
```bash
# Check sender balance (should decrease by 3.1)
go run . --client --cmd balance --address <sender-address>

# Check recipient balance (should increase by 3.0)
go run . --client --cmd balance --address <recipient-address>

# Check miner balance (should increase by 10.1)
go run . --client --cmd balance --address <miner-address>
```

### Scenario 3: Multi-Node P2P Network

#### 1. Start Multiple Nodes
```bash
# Terminal 1
go run . --node --port 8001 --peers localhost:8002,localhost:8003

# Terminal 2
go run . --node --port 8002 --peers localhost:8001,localhost:8003

# Terminal 3
go run . --node --port 8003 --peers localhost:8001,localhost:8002
```

**What happens internally:**
- **Peer Discovery**: Nodes automatically connect and announce themselves
- **Network Formation**: Creates mesh topology for transaction/block propagation
- **Blockchain Synchronization**: Nodes compare chain lengths and adopt longer valid chains
- **Genesis Consensus**: All nodes start with identical deterministic genesis block

#### 2. Create Transaction on Node 1
```bash
go run . --client --node-addr localhost:8001 --cmd send --address <from> --to <to> --amount 2.0 --fee 0.05
```

**What happens internally:**
- **Local Processing**: Node 1 validates and adds transaction to its pending pool
- **Network Broadcast**: Transaction propagates to all connected peers via HTTP POST
- **Duplicate Prevention**: Peers ignore transactions they've already seen

#### 3. Mine on Different Node
```bash
go run . --client --node-addr localhost:8002 --cmd mine
```

**What happens internally:**
- **Block Creation**: Node 2 mines block including transactions from its pending pool
- **Block Broadcast**: Mined block propagates to all peers
- **Consensus Validation**: All nodes validate and accept the block, updating their local state
- **Chain Synchronization**: Nodes automatically adopt longer valid chains (longest chain rule)
- **State Consistency**: UTXO sets synchronized across all nodes
- **Transaction Finality**: Pending transaction pools cleared across network

## 🔧 Implementation Deep Dive

### UTXO Set Management
Our implementation maintains a `map[string][]*UTXO` where:
- **Key**: Wallet address (Base58 encoded)
- **Value**: Array of unspent transaction outputs

This provides O(1) balance lookups and efficient transaction creation.

### Atomic State Updates
```go
// Blockchain ensures consistency with atomic UTXO updates
utxoBackup := b.copyUTXOSet()          // 1. Create backup
if err := b.validateTransactions(block); err != nil {  // 2. Validate
    return err
}
if err := b.applyUTXOChanges(block); err != nil {     // 3. Apply changes
    b.UTXOSet = utxoBackup             // 4. Rollback on failure
    return err
}
```

### Digital Signature Verification
Transactions include ECDSA signatures proving ownership:
- **Signing**: `tx.Sign(privateKey, prevTXs)`
- **Verification**: `tx.Verify(prevTXs)` 
- **Address Derivation**: `address = Base58(RIPEMD160(SHA256(publicKey)))`

### Mining Algorithm
```go
for {
    block.Hash = SHA256(block)
    if strings.HasPrefix(block.Hash, "0000") { // Difficulty = 4
        break
    }
    block.Nonce++
}
```

## 🔍 Key Features Demonstrated

### Blockchain Fundamentals
- ✅ **Immutable Ledger**: Cryptographically linked blocks
- ✅ **Consensus Mechanism**: Proof-of-Work with longest chain rule
- ✅ **Decentralization**: P2P network without central authority
- ✅ **Byzantine Fault Tolerance**: Network continues operating with malicious/failed nodes

### Cryptocurrency Features  
- ✅ **Double-Spend Prevention**: UTXO model ensures each coin spent only once
- ✅ **Cryptographic Security**: ECDSA signatures for transaction authorization
- ✅ **Monetary Policy**: Fixed block rewards with transaction fees

### Advanced Concepts
- ✅ **Atomic Transactions**: All-or-nothing UTXO updates with rollback
- ✅ **Network Consensus**: Automatic block propagation and validation
- ✅ **State Management**: Efficient UTXO set tracking and balance calculation
- ✅ **Chain Synchronization**: Longest chain rule with automatic peer sync
- ✅ **Fork Resolution**: Network automatically converges to single valid chain

## 🏗️ Current Limitations

- **No Persistence**: Blockchain resets on restart (in-memory only)
- **Fixed Difficulty**: Mining difficulty doesn't adjust based on block time
- **Simple P2P**: Basic HTTP-based networking without advanced peer discovery
- **Wallet Management**: Keys generated fresh each time (not persistent)

## 📚 Educational Value

This implementation demonstrates:

1. **UTXO vs Account Model**: Shows Bitcoin's approach vs Ethereum's account balances
2. **Proof-of-Work Security**: How computational puzzles secure the network
3. **P2P Consensus**: How nodes agree on blockchain state without central authority
4. **Cryptographic Primitives**: Real-world application of hash functions and digital signatures
5. **Distributed Systems**: Challenges of maintaining consistency across network nodes

## 🔧 Technical Stack

- **Language**: Go 1.19+
- **Cryptography**: ECDSA (secp256k1), SHA-256, RIPEMD-160
- **Networking**: HTTP/JSON REST API
- **Encoding**: Base58 (Bitcoin-compatible)
- **Concurrency**: Mutex-protected thread-safe operations

## 📖 Further Reading

- [Bitcoin Whitepaper](https://bitcoin.org/bitcoin.pdf) - Original cryptocurrency design
- [Mastering Bitcoin](https://github.com/bitcoinbook/bitcoinbook) - Technical deep dive
- [UTXO vs Account Model](https://medium.com/@sunflora98/utxo-vs-account-balance-model-5e6470f4e0cf)
- [Proof of Work Consensus](https://en.bitcoin.it/wiki/Proof_of_work)

---

*This educational blockchain demonstrates core cryptocurrency concepts through hands-on implementation. Perfect for understanding Bitcoin's architecture and distributed systems principles.*