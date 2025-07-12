# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an Ethereum blockchain interface demonstration written in Python. The codebase implements a simple Web3 interface that shows account management, balance checking, and transaction processing using the `web3.py` library with an Ethereum test provider.

## Architecture

The core functionality is organized into three primary modules:

- **EthereumBlockchainInterface** (`src/blockchain/interface.py`): Blockchain connection and information retrieval
- **AccountManager** (`src/accounts/manager.py`): Ethereum accounts, balance checking, and account creation  
- **TransactionManager** (`src/transactions/manager.py`): Standard and raw transaction operations

## Development Commands

### Environment Setup
```bash
# Install dependencies
pip3 install -r requirements.txt

# Run demos
python3 main.py          # Run basic demo (default)
python3 main.py 1        # Run basic demo explicitly
python3 main.py help     # Show available demos
python3 main.py basic    # Run basic demo by name
```

## Project Structure (Refactored)

```
web3-prototype/
├── main.py                    # Demo selector and entry point
├── requirements.txt           # Python dependencies
├── CLAUDE.md                 # This documentation file
└── src/                      # Source code modules
    ├── config.py             # Configuration and constants
    ├── blockchain/           # Blockchain connection and info
    │   └── interface.py      # EthereumBlockchainInterface
    ├── accounts/             # Account management
    │   └── manager.py        # AccountManager  
    ├── transactions/         # Transaction operations
    │   └── manager.py        # TransactionManager
    ├── utils/                # Utility modules
    │   ├── display.py        # Display formatting and console output
    │   └── security.py       # Security utilities and warnings
    └── demos/                # Learning demos
        ├── basic_demo.py     # Phase 1: Core concepts demo
        └── contracts_demo.py # Phase 2: Smart contracts (placeholder)
```

### Dependencies
The project uses:
- `web3~=6.20.4` - Main Web3 interface
- `web3[tester]` - Ethereum test provider
- `eth-typing~=3.5.2` - Ethereum type annotations
- `eth-account~=0.10.0` - Account management utilities

## Key Implementation Details

### Blockchain Connection
- Uses `Web3.EthereumTesterProvider()` for local testing
- Connection status checked via `w3.is_connected()`

### Account Management
- Supports both pre-existing test accounts (`w3.eth.accounts`) and newly created accounts (`w3.eth.account.create()`)
- Balance tracking in both Wei and ETH units
- Transaction history via transaction count

### Transaction Processing
- Standard transactions: Direct `send_transaction()` calls
- Raw transactions: Manual transaction building, signing, and broadcasting
- Transaction receipts and confirmation tracking

## Code Structure Notes

- All monetary values are handled using Web3's Wei conversion utilities
- Transaction processing includes proper gas estimation and nonce management
- Comprehensive logging with emoji indicators for better UX
- Error handling through Web3's built-in transaction receipt waiting

## Code Review & Refactoring Opportunities

### Current Issues Identified:
1. **Separation of Concerns**: Print statements mixed with business logic
2. **Code Duplication**: Similar balance/transaction count display logic in `print_account` and `print_account_details`
3. **Hard-coded Values**: Gas limit (21000) and Wei constants could be configurable
4. **Type Inconsistency**: `value_in_eth` parameter uses `int` but should be `float` for fractional ETH
5. **Missing Error Handling**: No try/catch blocks for Web3 operations
6. **Security Risk**: Private keys printed to console (line 61)

### Suggested Refactors:
1. **Extract Display Logic**: Create separate classes for formatting and display
2. **Consolidate Account Operations**: Merge similar account display methods
3. **Add Configuration**: Create constants file for gas limits, default values
4. **Improve Type Safety**: Use proper numeric types for ETH amounts
5. **Add Error Handling**: Wrap Web3 calls in try/catch blocks
6. **Security Enhancement**: Remove private key logging or add warning

## Web3.py Learning Plan

### Phase 1: Core Concepts (Current Implementation Coverage)
- [x] **Connection Management**: EthereumTesterProvider setup and connection checking
- [x] **Account Operations**: Account creation, balance checking, transaction counting
- [x] **Basic Transactions**: Standard and raw transaction sending
- [x] **Block Information**: Block number and block data retrieval
- [x] **Wei/ETH Conversion**: Using `to_wei()` and `from_wei()` utilities

### Phase 2: Advanced Features (Next Learning Steps)
- [ ] **Event Filtering**: Listen to blockchain events and logs
- [ ] **Smart Contracts**: Deploy and interact with contracts using ABI
- [ ] **Gas Estimation**: Dynamic gas price and limit calculation
- [ ] **Transaction Pools**: Pending transaction monitoring
- [ ] **Multiple Networks**: Connect to mainnet, testnets (Goerli, Sepolia)
- [ ] **Batch Operations**: Multiple transactions in single block
- [ ] **ENS Integration**: Ethereum Name Service resolution

### Phase 3: Production Features
- [ ] **Provider Management**: Infura, Alchemy, or local node connections
- [ ] **Wallet Integration**: MetaMask and hardware wallet support
- [ ] **Error Recovery**: Retry logic and failover mechanisms
- [ ] **Performance Optimization**: Connection pooling and caching
- [ ] **Security Patterns**: Private key management and secure signing
- [ ] **Monitoring**: Transaction status tracking and notifications

### Phase 4: Practical Exercises
1. **Contract Interaction**: Deploy ERC-20 token and perform transfers
2. **Event Monitoring**: Build transaction monitor for specific addresses
3. **Multi-Sig Wallet**: Implement basic multi-signature functionality
4. **DeFi Integration**: Interact with Uniswap or similar protocols
5. **NFT Operations**: Mint and transfer ERC-721 tokens
6. **Cross-Chain**: Bridge assets between different networks

### Learning Resources Integration
- Web3.py official documentation examples
- Ethereum development best practices
- Gas optimization techniques
- Security considerations for production use

## Refactoring Accomplishments

### ✅ Issues Resolved:
1. **Security Enhancement**: Private key logging now includes warnings and is configurable
2. **Separation of Concerns**: Display logic extracted to dedicated `ConsoleDisplay` class
3. **Code Duplication**: Consolidated account display methods with shared logic
4. **Error Handling**: Added comprehensive try/catch blocks around all Web3 operations
5. **Type Safety**: Updated ETH amount parameters to use `Union[int, float]`
6. **Configuration**: Extracted constants to `Web3Config` class for maintainability

### 🏗️ Architectural Improvements:
1. **Modular Design**: Split monolithic `main.py` into focused modules
2. **Utility Classes**: Created reusable display and security utilities
3. **Demo Framework**: Structured for progressive learning with separate demo scripts
4. **Error Recovery**: Graceful fallback mechanisms for gas estimation failures
5. **Extensibility**: Easy addition of new providers, demos, and features

### 📚 Learning Structure:
- **Phase 1**: Basic concepts (implemented in `basic_demo.py`)
- **Phase 2-4**: Placeholders created for advanced features
- **Progressive Complexity**: Each phase builds on previous knowledge
- **Practical Focus**: Emphasis on hands-on exercises and real-world patterns

The refactored codebase now provides a solid foundation for systematic Web3.py learning with production-quality patterns, comprehensive error handling, and clear separation of concerns.