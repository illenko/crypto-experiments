import warnings
from typing import Optional
from eth_account.signers.local import LocalAccount
from ..config import Web3Config


class SecurityManager:
    """Handles security-related operations and warnings"""
    
    @staticmethod
    def warn_private_key_exposure() -> None:
        """Issue warning about private key exposure"""
        if Web3Config.WARN_ON_PRIVATE_KEY_DISPLAY:
            warnings.warn(
                "Private key is being displayed. This should only be done in development/demo environments.",
                UserWarning,
                stacklevel=2
            )
    
    @staticmethod
    def should_display_private_key(demo_mode: bool = True) -> bool:
        """Determine if private key should be displayed based on environment"""
        # In a real application, this would check environment variables
        # For now, we'll only allow it in demo mode
        return demo_mode
    
    @staticmethod
    def safe_account_info(account: LocalAccount, include_private_key: bool = False) -> dict:
        """Return safe account information"""
        info = {
            'address': account.address,
            'has_private_key': True
        }
        
        if include_private_key:
            SecurityManager.warn_private_key_exposure()
            info['private_key'] = account.key.hex()
        
        return info
    
    @staticmethod
    def validate_address(address: str) -> bool:
        """Validate Ethereum address format"""
        try:
            from web3 import Web3
            return Web3.is_address(address)
        except Exception:
            return False
    
    @staticmethod
    def validate_transaction_params(to_address: str, value: int, gas_limit: Optional[int] = None) -> bool:
        """Validate transaction parameters"""
        if not SecurityManager.validate_address(to_address):
            return False
        
        if value < 0:
            return False
        
        if gas_limit is not None and gas_limit < Web3Config.DEFAULT_GAS_LIMIT:
            return False
        
        return True