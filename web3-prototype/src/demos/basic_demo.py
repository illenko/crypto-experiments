#!/usr/bin/env python3
"""
Basic Web3.py demonstration using the refactored modular structure.
This demo covers fundamental Web3 operations: connections, accounts, and transactions.
"""

from ..blockchain.interface import EthereumBlockchainInterface
from ..accounts.manager import AccountManager
from ..transactions.manager import TransactionManager
from ..utils.display import ConsoleDisplay


def run_basic_demo():
    """Run the basic Web3.py demonstration"""
    
    display = ConsoleDisplay()
    display.show_demo_start("Ethereum Blockchain Interface Demo")
    
    # Initialize components
    try:
        eth_interface = EthereumBlockchainInterface(display=display)
        account_manager = AccountManager(eth_interface.w3, display=display)
        tx_manager = TransactionManager(eth_interface.w3, display=display)
    except Exception as e:
        display.show_error("Failed to initialize components", e)
        return False
    
    # Check connection
    if not eth_interface.check_connection():
        display.show_error("Connection failed. Exiting...")
        return False
    
    # Display blockchain information
    eth_interface.display_blockchain_info()
    
    # Get and display accounts
    accounts = eth_interface.get_accounts()
    if not accounts:
        display.show_error("No accounts available")
        return False
    
    account_manager.display_accounts_list(accounts)
    
    # Display account details
    account_manager.display_indexed_account(accounts, 0)
    account_manager.display_indexed_account(accounts, 1)
    
    # Transaction operations
    display.show_section_header("Starting Transaction Operations")
    
    # Send transaction between existing accounts
    receipt = tx_manager.send_transaction(accounts[0], accounts[1], 3.0)
    if receipt:
        eth_interface.display_blockchain_info()
        account_manager.display_indexed_account(accounts, 0)
        account_manager.display_indexed_account(accounts, 1)
    
    # Create new account and fund it
    new_account = account_manager.create_new_account(display_private_key=True)
    if new_account:
        receipt = tx_manager.send_transaction(accounts[0], new_account.address, 2.0)
        if receipt:
            eth_interface.display_blockchain_info()
            account_manager.display_indexed_account(accounts, 0)
            account_manager.display_account_details(new_account.address, "New Account")
        
        # Send raw transaction from new account
        display.show_section_header("Sending raw transaction from new account to account #2")
        receipt = tx_manager.send_raw_transaction(
            new_account, 
            accounts[1], 
            eth_interface.one_eth_wei
        )
        if receipt:
            eth_interface.display_blockchain_info()
            account_manager.display_indexed_account(accounts, 0)
            account_manager.display_indexed_account(accounts, 1)
            account_manager.display_account_details(new_account.address, "New Account")
    
    display.show_demo_completed()
    return True


if __name__ == "__main__":
    success = run_basic_demo()
    exit(0 if success else 1)