from typing import Optional, Tuple
from web3 import Web3
from web3.exceptions import Web3Exception
from eth_account.signers.local import LocalAccount
from eth_typing import ChecksumAddress
from ..utils.display import ConsoleDisplay
from ..utils.security import SecurityManager


class AccountManager:
    """Enhanced account manager with security and error handling"""
    
    def __init__(self, w3: Web3, display: ConsoleDisplay = None):
        self.w3 = w3
        self.display = display or ConsoleDisplay()
        self.security = SecurityManager()
    
    def get_balance(self, address: ChecksumAddress) -> Optional[int]:
        """Get account balance with error handling"""
        try:
            return self.w3.eth.get_balance(address)
        except Web3Exception as e:
            self.display.show_error(f"Failed to get balance for {address}", e)
            return None
    
    def get_transaction_count(self, address: ChecksumAddress) -> Optional[int]:
        """Get transaction count with error handling"""
        try:
            return self.w3.eth.get_transaction_count(address)
        except Web3Exception as e:
            self.display.show_error(f"Failed to get transaction count for {address}", e)
            return None
    
    def display_accounts_list(self, accounts: Tuple[ChecksumAddress, ...]) -> None:
        """Display list of available accounts"""
        if accounts:
            self.display.show_accounts_list(accounts)
        else:
            self.display.show_warning("No accounts available")
    
    def display_account_details(self, address: ChecksumAddress, label: str = None) -> bool:
        """Display account details with error handling"""
        balance = self.get_balance(address)
        tx_count = self.get_transaction_count(address)
        
        if balance is not None and tx_count is not None:
            self.display.show_account_details(address, balance, tx_count, label)
            return True
        return False
    
    def display_indexed_account(self, accounts: Tuple[ChecksumAddress, ...], 
                              account_index: int) -> bool:
        """Display account details by index with bounds checking"""
        if not accounts:
            self.display.show_error("No accounts available")
            return False
        
        if account_index < 0 or account_index >= len(accounts):
            self.display.show_error(f"Account index {account_index} out of range")
            return False
        
        address = accounts[account_index]
        label = f"Account #{account_index + 1}"
        return self.display_account_details(address, label)
    
    def create_new_account(self, display_private_key: bool = False) -> Optional[LocalAccount]:
        """Create a new Ethereum account with security considerations"""
        try:
            new_account = self.w3.eth.account.create()
            
            # Display account creation with security warnings
            self.display.show_account_created(
                new_account, 
                show_private_key=display_private_key and self.security.should_display_private_key()
            )
            
            # Display balance
            balance = self.get_balance(new_account.address)
            if balance is not None:
                print(f"  💰 Initial Balance: {self.display.formatter.format_eth_balance(balance)}")
            
            return new_account
            
        except Exception as e:
            self.display.show_error("Failed to create new account", e)
            return None
    
    def validate_account_access(self, address: ChecksumAddress) -> bool:
        """Validate that account can be accessed"""
        try:
            balance = self.w3.eth.get_balance(address)
            return balance is not None
        except Web3Exception:
            return False