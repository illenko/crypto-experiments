from typing import Tuple, List, Dict, Any
from web3 import Web3
from web3.types import BlockData
from eth_account.signers.local import LocalAccount
from eth_typing import ChecksumAddress
from ..config import Web3Config


class DisplayFormatter:
    """Handles all display formatting and console output"""
    
    @staticmethod
    def format_eth_balance(balance_wei: int) -> str:
        """Format Wei balance to ETH with proper decimal places"""
        eth_balance = Web3.from_wei(balance_wei, 'ether')
        return f"{eth_balance:,.{Web3Config.ETH_DECIMAL_PLACES}f} ETH"
    
    @staticmethod
    def format_address(address: ChecksumAddress, label: str = "Address") -> str:
        """Format an Ethereum address with label"""
        return f"{label}: {address}"
    
    @staticmethod
    def format_transaction_count(count: int) -> str:
        """Format transaction count"""
        return f"{count} Transactions"


class ConsoleDisplay:
    """Handles console output for Web3 operations"""
    
    def __init__(self, formatter: DisplayFormatter = None):
        self.formatter = formatter or DisplayFormatter()
    
    def show_connection_status(self, is_connected: bool) -> None:
        """Display connection status"""
        status = '✅ Connected' if is_connected else '❌ Not Connected'
        print(f"\n🔌 Connection Status: {status}")
    
    def show_blockchain_info(self, block_number: int, latest_block: BlockData) -> None:
        """Display blockchain information"""
        print("\n🔗 Blockchain Information:")
        print(f"  📦 Current Block Number: {block_number}")
        print(f"  📄 Latest Block: {latest_block}")
    
    def show_accounts_list(self, accounts: Tuple[ChecksumAddress, ...]) -> None:
        """Display list of available accounts"""
        print("\n👥 Available Ethereum Accounts:")
        for idx, account in enumerate(accounts, 1):
            print(f"  🏦 Account #{idx}: {account}")
    
    def show_account_details(self, address: ChecksumAddress, balance_wei: int, 
                           tx_count: int, label: str = None) -> None:
        """Display account balance and transaction history"""
        account_label = label or f"Account {address}"
        print(f"\n💰 Balance of {account_label}:")
        print(f"  💎 {self.formatter.format_eth_balance(balance_wei)}")
        print(f"\n📝 Transaction History of {account_label}:")
        print(f"  📊 {self.formatter.format_transaction_count(tx_count)}")
    
    def show_transaction_initiated(self, from_addr: ChecksumAddress, 
                                 to_addr: ChecksumAddress, amount_eth: float) -> None:
        """Display transaction initiation"""
        print(f"💸 Sending {amount_eth} ETH from {from_addr} to {to_addr}...")
    
    def show_transaction_hash(self, tx_hash: bytes) -> None:
        """Display transaction hash"""
        print(f"  🔍 Transaction Hash: {tx_hash.hex()}")
    
    def show_transaction_mining(self) -> None:
        """Display mining status"""
        print("⏳ Waiting for transaction to be mined...")
    
    def show_transaction_completed(self, transaction_data: Dict[str, Any]) -> None:
        """Display completed transaction"""
        print("✅ Transaction has been mined!")
        print(f"  📄 Transaction Details: {transaction_data}")
    
    def show_account_created(self, account: LocalAccount, show_private_key: bool = False) -> None:
        """Display new account creation with security warning"""
        print(f"\n✨ Created Account: {account.address}")
        
        if show_private_key:
            if Web3Config.WARN_ON_PRIVATE_KEY_DISPLAY:
                print("⚠️  WARNING: Private key displayed for demo purposes only!")
                print("⚠️  Never share private keys in production!")
            print(f"  🔑 Private Key: {account.key.hex()}")
        else:
            print("  🔑 Private Key: [HIDDEN FOR SECURITY]")
    
    def show_demo_start(self, demo_name: str) -> None:
        """Display demo start message"""
        print(f"\n🌟 Starting {demo_name} 🌟")
    
    def show_demo_completed(self) -> None:
        """Display demo completion message"""
        print("\n🎉 Demo Completed Successfully! 🎉")
    
    def show_section_header(self, section_name: str) -> None:
        """Display section header"""
        print(f"\n🔄 {section_name}:")
    
    def show_error(self, error_msg: str, exception: Exception = None) -> None:
        """Display error message"""
        print(f"❌ Error: {error_msg}")
        if exception and hasattr(exception, '__str__'):
            print(f"   Details: {str(exception)}")
    
    def show_warning(self, warning_msg: str) -> None:
        """Display warning message"""
        print(f"⚠️  Warning: {warning_msg}")