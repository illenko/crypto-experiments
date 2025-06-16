import hashlib
import json
from time import time
import os
import requests
import structlog
from typing import Dict, List, Optional, Any
from uuid import uuid4
from flask import Flask, jsonify, request
from urllib.parse import urlparse

logger = structlog.get_logger()
PORT = int(os.getenv('BLOCKCHAIN_PORT', 8080))
HOST = os.getenv('BLOCKCHAIN_HOST', '0.0.0.0')

class Blockchain:
    DIFFICULTY = "0000"

    def __init__(self):
        self.current_transactions: List[Dict[str, Any]] = []
        self.chain: List[Dict[str, Any]] = []
        self.nodes = set()
        logger.info("🔗 Initializing blockchain")
        self._create_genesis_block()

    def _create_genesis_block(self) -> None:
        logger.info("⚡ Creating genesis block")
        self.new_block(proof=100, previous_hash="0" * 64)

    def register_node(self, address: str) -> None:
        parsed_url = urlparse(address)
        if parsed_url.netloc:
            self.nodes.add(parsed_url.netloc)
            logger.info(f"📡 Registered new node: {parsed_url.netloc}")
        else:
            logger.error(f"⚠️ Invalid node address format: {address}")
            raise ValueError("Invalid node address")

    def new_block(self, proof: int, previous_hash: Optional[str] = None) -> Dict[str, Any]:
        block = {
            'index': len(self.chain) + 1,
            'timestamp': time(),
            'transactions': self.current_transactions,
            'proof': proof,
            'previous_hash': previous_hash or self.hash(self.chain[-1])
        }

        block_string = self._get_block_string(block)
        current_hash = None

        logger.info(f"⛏️ Mining block {block['index']}")
        while not current_hash or not current_hash.startswith(self.DIFFICULTY):
            current_hash = self._calculate_hash(block_string, block['proof'])
            block['proof'] += 1

        block['proof'] -= 1
        block['hash'] = current_hash

        logger.info(f"💎 Block {block['index']} mined with {len(block['transactions'])} transactions")

        self.current_transactions = []
        self.chain.append(block)
        return block

    def new_transaction(self, sender: str, recipient: str, amount: float) -> int:
        self.current_transactions.append({
            'sender': sender,
            'recipient': recipient,
            'amount': amount,
        })
        next_index = self.last_block['index'] + 1
        logger.info(f"💸 New transaction: {amount} coins from {sender} to {recipient}")
        return next_index

    def get_pending_transactions(self) -> List[Dict[str, Any]]:
        return self.current_transactions

    @staticmethod
    def _get_block_string(block: Dict[str, Any]) -> str:
        block_copy = block.copy()
        block_copy.pop('proof')
        block_copy.pop('hash', None)
        return json.dumps(block_copy, sort_keys=True)

    @staticmethod
    def _calculate_hash(block_string: str, proof: int) -> str:
        guess = f'{block_string}{proof}'.encode()
        return hashlib.sha256(guess).hexdigest()

    @staticmethod
    def hash(block: Dict[str, Any]) -> str:
        block_string = Blockchain._get_block_string(block)
        return Blockchain._calculate_hash(block_string, block['proof'])

    @property
    def last_block(self) -> Dict[str, Any]:
        return self.chain[-1]

    def get_block_by_index(self, index: int) -> Optional[Dict[str, Any]]:
        if 0 <= index < len(self.chain):
            return self.chain[index]
        return None

    def is_valid_chain(self, chain: List[Dict[str, Any]]) -> bool:
        for i in range(1, len(chain)):
            current_block = chain[i]
            previous_block = chain[i - 1]

            if current_block['previous_hash'] != self.hash(previous_block):
                logger.error(f"❌ Invalid previous hash at block {current_block['index']}")
                return False

            if not current_block['hash'].startswith(self.DIFFICULTY):
                logger.error(f"❌ Invalid proof of work at block {current_block['index']}")
                return False

            calculated_hash = self.hash(current_block)
            if calculated_hash != current_block['hash']:
                logger.error(f"❌ Invalid block hash at block {current_block['index']}")
                return False

        logger.info("✅ External chain validation successful")
        return True

    def resolve_conflicts(self) -> bool:
        max_length = len(self.chain)
        new_chain = None

        logger.info("🔄 Starting chain resolution with network nodes")
        for node in self.nodes:
            try:
                response = requests.get(f'http://{node}/chain')
                if response.status_code == 200:
                    length = response.json()['length']
                    chain = response.json()['chain']

                    if length > max_length and self.is_valid_chain(chain):
                        max_length = length
                        new_chain = chain
                        logger.info(f"📡 Found longer valid chain from {node}, length: {length}")
            except requests.RequestException as e:
                logger.error(f"❌ Failed to connect to node {node}")

        if new_chain:
            self.chain = new_chain
            logger.info(f"🔁 Chain replaced with longer chain of length {len(new_chain)}")
            return True

        logger.info("✅ Current chain is up to date")
        return False


app = Flask(__name__)
node_identifier = str(uuid4()).replace('-', '')
blockchain = Blockchain()


@app.route('/mine', methods=['GET'])
def mine():
    last_block = blockchain.last_block
    blockchain.new_transaction(
        sender="0",
        recipient=node_identifier,
        amount=1,
    )

    previous_hash = blockchain.hash(last_block)
    block = blockchain.new_block(proof=0, previous_hash=previous_hash)

    response = {
        'message': "New Block Forged",
        'index': block['index'],
        'transactions': block['transactions'],
        'proof': block['proof'],
        'previous_hash': block['previous_hash'],
        'hash': block['hash']
    }
    logger.info(f"✨ Successfully mined block {block['index']}")
    return jsonify(response), 200


@app.route('/transactions', methods=['POST'])
def new_transaction():
    values = request.get_json()

    required = ['sender', 'recipient', 'amount']
    if not all(k in values for k in required):
        logger.error("⚠️ Missing transaction values")
        return jsonify({'error': 'Missing values'}), 400

    index = blockchain.new_transaction(
        values['sender'],
        values['recipient'],
        values['amount']
    )

    response = {'message': f'Transaction will be added to Block {index}'}
    return jsonify(response), 201


@app.route('/transactions', methods=['GET'])
def get_pending_transactions():
    pending = blockchain.get_pending_transactions()
    response = {
        'pending_count': len(pending),
        'transactions': pending,
        'next_block_index': blockchain.last_block['index'] + 1
    }
    logger.info(f"📋 Retrieved {len(pending)} pending transactions")
    return jsonify(response), 200


@app.route('/chain', methods=['GET'])
def full_chain():
    response = {
        'chain': blockchain.chain,
        'length': len(blockchain.chain),
        'valid': blockchain.is_valid_chain(blockchain.chain)
    }
    return jsonify(response), 200


@app.route('/block/<int:index>', methods=['GET'])
def get_block(index):
    block = blockchain.get_block_by_index(index)
    if block is None:
        logger.error(f"🔍 Block {index} not found")
        return jsonify({'error': 'Block not found'}), 404

    logger.info(f"📦 Retrieved block {index}")
    return jsonify(block), 200


@app.route('/nodes/register', methods=['POST'])
def register_nodes():
    values = request.get_json()
    nodes = values.get('nodes')

    if not nodes:
        logger.error("⚠️ No nodes provided for registration")
        return jsonify({'error': 'Please supply a valid list of nodes'}), 400

    registered = []
    failed = []

    for node in nodes:
        try:
            blockchain.register_node(node)
            registered.append(node)
        except ValueError:
            failed.append(node)

    response = {
        'message': '🌐 Nodes registration processed',
        'registered_nodes': registered,
        'total_nodes': list(blockchain.nodes),
        'failed_nodes': failed
    }

    if failed:
        logger.warning(f"⚠️ Failed to register {len(failed)} nodes")
        return jsonify(response), 207  # Partial Content

    logger.info(f"✅ Successfully registered {len(registered)} new nodes")
    return jsonify(response), 201


@app.route('/nodes/resolve', methods=['GET'])
def consensus():
    logger.info("🔄 Starting consensus resolution")
    replaced = blockchain.resolve_conflicts()

    response = {
        'message': '🔄 Chain replaced with network consensus' if replaced else '✅ Local chain is authoritative',
        'chain_replaced': replaced,
        'chain_length': len(blockchain.chain),
        'chain': blockchain.chain
    }

    logger.info(
        "📊 Consensus resolution completed",
        chain_replaced=replaced,
        chain_length=len(blockchain.chain)
    )
    return jsonify(response), 200


if __name__ == '__main__':
    logger.info(
        "🚀 Starting blockchain node",
        host=HOST,
        port=PORT
    )
    app.run(
        host=HOST,
        port=PORT,
        debug=os.getenv('FLASK_DEBUG', 'false').lower() == 'true'
    )

