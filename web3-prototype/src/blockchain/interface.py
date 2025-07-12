from typing import Optional, Dict, Any
from web3 import Web3
from web3.types import BlockData
from web3.exceptions import Web3Exception
from ..config import Web3Config
from ..utils.display import ConsoleDisplay


class EthereumBlockchainInterface:
    """Enhanced blockchain interface with error handling and configuration"""
    
    def __init__(self, provider_type: str = None, display: ConsoleDisplay = None):
        self.config = Web3Config()
        self.display = display or ConsoleDisplay()
        
        try:
            provider_config = self.config.get_provider_config(provider_type)
            provider_class = provider_config['provider_class']
            provider_args = provider_config['provider_args']
            provider_kwargs = provider_config['provider_kwargs']
            
            self.w3 = Web3(provider_class(*provider_args, **provider_kwargs))
            
        except Exception as e:
            self.display.show_error("Failed to initialize Web3 provider", e)
            raise
    
    def check_connection(self) -> bool:
        """Check if connected to Ethereum network with error handling"""
        try:
            is_connected = self.w3.is_connected()
            self.display.show_connection_status(is_connected)
            return is_connected
        except Web3Exception as e:
            self.display.show_error("Connection check failed", e)
            return False
    
    def get_blockchain_info(self) -> Optional[Dict[str, Any]]:
        """Get current blockchain information with error handling"""
        try:
            block_number = self.w3.eth.block_number
            latest_block = self.w3.eth.get_block('latest')
            
            blockchain_info = {
                'block_number': block_number,
                'latest_block': latest_block
            }
            
            return blockchain_info
            
        except Web3Exception as e:
            self.display.show_error("Failed to retrieve blockchain information", e)
            return None
    
    def display_blockchain_info(self) -> None:
        """Display blockchain information"""
        info = self.get_blockchain_info()
        if info:
            self.display.show_blockchain_info(
                info['block_number'], 
                info['latest_block']
            )
    
    def get_accounts(self) -> Optional[tuple]:
        """Get available accounts with error handling"""
        try:
            return self.w3.eth.accounts
        except Web3Exception as e:
            self.display.show_error("Failed to retrieve accounts", e)
            return None
    
    def get_gas_price(self) -> Optional[int]:
        """Get current gas price with error handling"""
        try:
            return self.w3.eth.gas_price
        except Web3Exception as e:
            self.display.show_error("Failed to retrieve gas price", e)
            return None
    
    def estimate_gas(self, transaction: Dict[str, Any]) -> Optional[int]:
        """Estimate gas for transaction with error handling"""
        try:
            estimated = self.w3.eth.estimate_gas(transaction)
            # Apply safety multiplier
            return int(estimated * self.config.DEFAULT_GAS_MULTIPLIER)
        except Web3Exception as e:
            self.display.show_error("Gas estimation failed", e)
            return None
    
    @property
    def one_eth_wei(self) -> int:
        """Get 1 ETH in Wei"""
        return self.config.ONE_ETH_WEI
    
    @property
    def two_eth_wei(self) -> int:
        """Get 2 ETH in Wei"""
        return self.config.TWO_ETH_WEI