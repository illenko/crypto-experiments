from typing import Tuple, Dict, Any
from eth_account.signers.local import LocalAccount
from eth_typing import ChecksumAddress
from web3 import Web3
from web3.types import TxReceipt, TxParams


class EthereumBlockchainInterface:
    def __init__(self):
        self.w3 = Web3(Web3.EthereumTesterProvider())
        self.one_eth = Web3.to_wei(1, 'ether')
        self.two_eth = Web3.to_wei(2, 'ether')

    def check_connection(self) -> bool:
        """Check if connected to Ethereum network"""
        is_connected = self.w3.is_connected()
        print(f"\n🔌 Connection Status: {'✅ Connected' if is_connected else '❌ Not Connected'}")
        return is_connected

    def print_blockchain_info(self) -> None:
        """Display current blockchain information"""
        print("\n🔗 Blockchain Information:")
        block_number = self.w3.eth.block_number
        print(f"  📦 Current Block Number: {block_number}")
        block = self.w3.eth.get_block('latest')
        print(f"  📄 Latest Block: {block}")


class AccountManager:
    def __init__(self, w3: Web3):
        self.w3 = w3

    def print_accounts(self, accounts: Tuple[ChecksumAddress]) -> None:
        """Display all available Ethereum accounts"""
        print("\n👥 Available Ethereum Accounts:")
        for idx, account in enumerate(accounts, 1):
            print(f"  🏦 Account #{idx}: {account}")

    def print_account(self, accounts: Tuple[ChecksumAddress], account_index: int = 0) -> None:
        """Display balance and transaction history for a specific account"""
        print(f"\n💰 Balance of Account #{account_index}:")
        balance = self.w3.eth.get_balance(accounts[account_index])
        print(f"  💎 {Web3.from_wei(balance, 'ether'):,.8f} ETH")
        print(f"\n📝 Transaction History of Account #{account_index}:")
        tx_history = self.w3.eth.get_transaction_count(accounts[account_index])
        print(f"  📊 {tx_history} Transactions")

    def print_account_details(self, account: LocalAccount) -> None:
        """Display balance and transaction history for a LocalAccount"""
        print(f"\n💰 Balance of Account {account.address}:")
        balance = self.w3.eth.get_balance(account.address)
        print(f"  💎 {Web3.from_wei(balance, 'ether'):,.8f} ETH")
        print(f"\n📝 Transaction History of Account {account.address}:")
        tx_history = self.w3.eth.get_transaction_count(account.address)
        print(f"  📊 {tx_history} Transactions")

    def create_new_account(self) -> LocalAccount:
        """Create a new Ethereum account"""
        new_account = self.w3.eth.account.create()
        print(f"\n✨ Created Account: {new_account.address}")
        print(f"  🔑 Private Key: {new_account.key.hex()}")
        print(f"  💰 Balance: {Web3.from_wei(self.w3.eth.get_balance(new_account.address), 'ether'):,.8f} ETH")
        return new_account


class TransactionManager:
    def __init__(self, w3: Web3):
        self.w3 = w3

    def send_transaction(self, from_address: ChecksumAddress, to_address: ChecksumAddress,
                         value_in_eth: int) -> TxReceipt:
        """Send ETH from one account to another"""
        print(f"💸 Sending {value_in_eth} ETH from {from_address} to {to_address}...")
        tx_hash = self.w3.eth.send_transaction({
            'from': from_address,
            'to': to_address,
            'value': self.w3.to_wei(value_in_eth, 'ether')
        })
        return self._process_transaction(tx_hash)

    def send_raw_transaction(self, account: LocalAccount, to_address: ChecksumAddress,
                             value_in_wei: int) -> TxReceipt:
        """Send a raw transaction from a LocalAccount"""
        tx = self._build_raw_transaction(account.address, to_address, value_in_wei)
        signed_tx = account.sign_transaction(tx)
        tx_hash = self.w3.eth.send_raw_transaction(signed_tx.rawTransaction)
        return self._process_transaction(tx_hash)

    def _build_raw_transaction(self, from_address: ChecksumAddress, to_address: ChecksumAddress,
                               value: int) -> TxParams:
        """Build a raw transaction object"""
        return {
            'to': to_address,
            'value': value,
            'gas': 21000,
            'gasPrice': self.w3.eth.gas_price,
            'nonce': self.w3.eth.get_transaction_count(from_address)
        }

    def _process_transaction(self, tx_hash: bytes) -> TxReceipt:
        """Process and wait for transaction completion"""
        print(f"  🔍 Transaction Hash: {tx_hash.hex()}")
        print("⏳ Waiting for transaction to be mined...")
        receipt = self.w3.eth.wait_for_transaction_receipt(tx_hash)
        print("✅ Transaction has been mined!")
        transaction = self.w3.eth.get_transaction(tx_hash)
        print(f"  📄 Transaction Details: {transaction}")
        return receipt


def main():
    print("\n🌟 Starting Ethereum Blockchain Interface Demo 🌟")

    eth_interface = EthereumBlockchainInterface()
    account_manager = AccountManager(eth_interface.w3)
    tx_manager = TransactionManager(eth_interface.w3)

    if not eth_interface.check_connection():
        print("❌ Connection failed. Exiting...")
        return
    eth_interface.print_blockchain_info()

    accounts = eth_interface.w3.eth.accounts
    account_manager.print_accounts(accounts)
    account_manager.print_account(accounts, 0)
    account_manager.print_account(accounts, 1)

    print("\n🔄 Starting Transaction Operations:")

    tx_manager.send_transaction(accounts[0], accounts[1], 3)
    eth_interface.print_blockchain_info()
    account_manager.print_account(accounts, 0)
    account_manager.print_account(accounts, 1)

    new_account = account_manager.create_new_account()
    tx_manager.send_transaction(accounts[0], new_account.address, 2)
    eth_interface.print_blockchain_info()
    account_manager.print_account(accounts, 0)
    account_manager.print_account_details(new_account)

    print("\n📝 Sending 1 ETH from new account to account #1 manually")
    tx_manager.send_raw_transaction(new_account, accounts[1], eth_interface.one_eth)
    eth_interface.print_blockchain_info()
    account_manager.print_account(accounts, 0)
    account_manager.print_account(accounts, 1)
    account_manager.print_account_details(new_account)

    print("\n🎉 Demo Completed Successfully! 🎉")


if __name__ == "__main__":
    main()