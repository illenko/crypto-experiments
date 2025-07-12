from web3 import Web3
from typing import Dict, Any


class Web3Config:
    """Configuration constants and settings for Web3 operations"""
    
    # Gas settings
    DEFAULT_GAS_LIMIT = 21000
    DEFAULT_GAS_MULTIPLIER = 1.2
    
    # ETH unit constants
    ONE_ETH_WEI = Web3.to_wei(1, 'ether')
    TWO_ETH_WEI = Web3.to_wei(2, 'ether')
    
    # Network settings
    DEFAULT_PROVIDER = 'tester'
    
    # Security settings
    WARN_ON_PRIVATE_KEY_DISPLAY = True
    
    # Display settings
    ETH_DECIMAL_PLACES = 8
    
    @classmethod
    def get_provider_config(cls, provider_type: str = None) -> Dict[str, Any]:
        """Get provider configuration based on type"""
        provider_type = provider_type or cls.DEFAULT_PROVIDER
        
        configs = {
            'tester': {
                'provider_class': Web3.EthereumTesterProvider,
                'provider_args': (),
                'provider_kwargs': {}
            },
            # Future providers can be added here
            # 'infura': {...},
            # 'alchemy': {...},
        }
        
        return configs.get(provider_type, configs['tester'])