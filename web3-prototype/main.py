from typing import Tuple

from eth_typing import ChecksumAddress
from web3 import Web3


def print_blockchain():
    print("\n Blockchain Information:")
    block_number = w3.eth.block_number
    print(f"  Current Block Number: {block_number}")
    block = w3.eth.get_block('latest')
    print(f"  Latest Block: {block}")

def print_account(accounts: Tuple[ChecksumAddress], account_index: int = 0):
    print(f"\n Balance of Account #{account_index}:")
    balance = w3.eth.get_balance(accounts[account_index])
    print(f"  {Web3.from_wei(balance, 'ether'):,.8f} ETH")
    print(f"\n Transaction History of Account #{account_index}:")
    tx_history = w3.eth.get_transaction_count(accounts[account_index])
    print(f"  {tx_history} Transactions")


def print_accounts(accounts: Tuple[ChecksumAddress]):
    print("\n👥 Available Ethereum Accounts:")
    for idx, account in enumerate(accounts, 1):
        print(f"  Account #{idx}: {account}")


w3 = Web3(Web3.EthereumTesterProvider())

wei_amount = Web3.to_wei(1, 'ether')
print(f"\n📈 1 ETH in Wei: {wei_amount:,} Wei")

eth_amount = Web3.from_wei(500000000, 'ether')
print(f"📉 500,000,000 Wei in ETH: {eth_amount} ETH")

is_connected = w3.is_connected()
print(f"\n🔌 Connection Status: {'✅ Connected' if is_connected else '❌ Not Connected'}")

print_blockchain()
print_accounts(w3.eth.accounts)

print_account(w3.eth.accounts, 0)
print_account(w3.eth.accounts, 1)

print("Sending 3 ETH from 0 to 1...")
tx_hash = w3.eth.send_transaction({
    'from': w3.eth.accounts[0],
    'to': w3.eth.accounts[1],
    'value': w3.to_wei(3, 'ether')
})

print(f"  Transaction Hash: {tx_hash.hex()}")

print("Waiting for transaction to be mined...")
w3.eth.wait_for_transaction_receipt(tx_hash)

transaction = w3.eth.get_transaction(tx_hash)

print("Transaction has been mined!")
print(f"  Transaction Details: {transaction}")

print_blockchain()
print_account(w3.eth.accounts, 0)
print_account(w3.eth.accounts, 1)

