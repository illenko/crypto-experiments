from typing import Optional, Dict, Any, Union
from web3 import Web3
from web3.types import TxReceipt, TxParams
from web3.exceptions import Web3Exception
from eth_account.signers.local import LocalAccount
from eth_typing import ChecksumAddress
from ..config import Web3Config
from ..utils.display import ConsoleDisplay
from ..utils.security import SecurityManager


class TransactionManager:
    """Enhanced transaction manager with error handling and validation"""
    
    def __init__(self, w3: Web3, display: ConsoleDisplay = None):
        self.w3 = w3
        self.display = display or ConsoleDisplay()
        self.security = SecurityManager()
        self.config = Web3Config()
    
    def send_transaction(self, from_address: ChecksumAddress, to_address: ChecksumAddress,
                        value_in_eth: Union[int, float]) -> Optional[TxReceipt]:
        """Send ETH from one account to another with improved error handling"""
        
        # Validate parameters
        if not self.security.validate_address(to_address):
            self.display.show_error(f"Invalid destination address: {to_address}")
            return None
        
        if value_in_eth <= 0:
            self.display.show_error("Transaction value must be positive")
            return None
        
        try:
            value_in_wei = Web3.to_wei(value_in_eth, 'ether')
            
            # Check sender balance
            sender_balance = self.w3.eth.get_balance(from_address)
            if sender_balance < value_in_wei:
                self.display.show_error(
                    f"Insufficient balance. Required: {value_in_eth} ETH, "
                    f"Available: {Web3.from_wei(sender_balance, 'ether')} ETH"
                )
                return None
            
            self.display.show_transaction_initiated(from_address, to_address, value_in_eth)
            
            # Build transaction
            tx_params = {
                'from': from_address,
                'to': to_address,
                'value': value_in_wei
            }
            
            # Estimate gas
            estimated_gas = self.w3.eth.estimate_gas(tx_params)
            tx_params['gas'] = int(estimated_gas * self.config.DEFAULT_GAS_MULTIPLIER)
            
            # Send transaction
            tx_hash = self.w3.eth.send_transaction(tx_params)
            return self._process_transaction(tx_hash)
            
        except Web3Exception as e:
            self.display.show_error("Transaction failed", e)
            return None
    
    def send_raw_transaction(self, account: LocalAccount, to_address: ChecksumAddress,
                           value_in_wei: int) -> Optional[TxReceipt]:
        """Send a raw transaction from a LocalAccount with validation"""
        
        # Validate parameters
        if not self.security.validate_transaction_params(to_address, value_in_wei):
            self.display.show_error("Invalid transaction parameters")
            return None
        
        try:
            # Check balance
            sender_balance = self.w3.eth.get_balance(account.address)
            if sender_balance < value_in_wei:
                self.display.show_error(
                    f"Insufficient balance. Required: {Web3.from_wei(value_in_wei, 'ether')} ETH, "
                    f"Available: {Web3.from_wei(sender_balance, 'ether')} ETH"
                )
                return None
            
            # Build and sign transaction
            tx = self._build_raw_transaction(account.address, to_address, value_in_wei)
            if not tx:
                return None
            
            signed_tx = account.sign_transaction(tx)
            tx_hash = self.w3.eth.send_raw_transaction(signed_tx.rawTransaction)
            return self._process_transaction(tx_hash)
            
        except Web3Exception as e:
            self.display.show_error("Raw transaction failed", e)
            return None
    
    def _build_raw_transaction(self, from_address: ChecksumAddress, to_address: ChecksumAddress,
                              value: int) -> Optional[TxParams]:
        """Build a raw transaction object with error handling"""
        try:
            # Get current gas price
            gas_price = self.w3.eth.gas_price
            
            # Get nonce
            nonce = self.w3.eth.get_transaction_count(from_address)
            
            # Build transaction
            tx_params = {
                'to': to_address,
                'value': value,
                'gas': self.config.DEFAULT_GAS_LIMIT,
                'gasPrice': gas_price,
                'nonce': nonce
            }
            
            # Try to estimate gas for better accuracy
            try:
                # Create a copy for estimation (without gas field)
                estimation_params = {k: v for k, v in tx_params.items() if k != 'gas'}
                estimated_gas = self.w3.eth.estimate_gas(estimation_params)
                tx_params['gas'] = int(estimated_gas * self.config.DEFAULT_GAS_MULTIPLIER)
            except (Web3Exception, Exception):
                # Fall back to default gas limit
                self.display.show_warning("Using default gas limit due to estimation failure")
            
            return tx_params
            
        except Web3Exception as e:
            self.display.show_error("Failed to build transaction", e)
            return None
    
    def _process_transaction(self, tx_hash: bytes) -> Optional[TxReceipt]:
        """Process and wait for transaction completion with error handling"""
        try:
            self.display.show_transaction_hash(tx_hash)
            self.display.show_transaction_mining()
            
            # Wait for transaction receipt
            receipt = self.w3.eth.wait_for_transaction_receipt(tx_hash)
            
            # Check if transaction was successful
            if receipt.status == 1:
                # Get full transaction details
                transaction = self.w3.eth.get_transaction(tx_hash)
                self.display.show_transaction_completed(dict(transaction))
                return receipt
            else:
                self.display.show_error("Transaction was mined but failed")
                return None
                
        except Web3Exception as e:
            self.display.show_error("Transaction processing failed", e)
            return None
    
    def get_transaction_details(self, tx_hash: str) -> Optional[Dict[str, Any]]:
        """Get transaction details by hash"""
        try:
            transaction = self.w3.eth.get_transaction(tx_hash)
            return dict(transaction)
        except Web3Exception as e:
            self.display.show_error(f"Failed to get transaction {tx_hash}", e)
            return None
    
    def estimate_transaction_cost(self, from_address: ChecksumAddress, 
                                to_address: ChecksumAddress, value_in_wei: int) -> Optional[Dict[str, int]]:
        """Estimate total transaction cost including gas"""
        try:
            tx_params = {
                'from': from_address,
                'to': to_address,
                'value': value_in_wei
            }
            
            gas_estimate = self.w3.eth.estimate_gas(tx_params)
            gas_price = self.w3.eth.gas_price
            gas_cost = gas_estimate * gas_price
            total_cost = value_in_wei + gas_cost
            
            return {
                'value': value_in_wei,
                'gas_estimate': gas_estimate,
                'gas_price': gas_price,
                'gas_cost': gas_cost,
                'total_cost': total_cost
            }
            
        except Web3Exception as e:
            self.display.show_error("Cost estimation failed", e)
            return None