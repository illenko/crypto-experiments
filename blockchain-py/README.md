# Blockchain Node

A simple blockchain implementation with Flask REST API.

## Requests

### 1. Get Full Chain

```shell
curl -X GET http://localhost:8080/chain \
  -H "Content-Type: application/json"
```

### 2. Add New Transaction
```shell
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "sender": "8527147fe1f5426f9dd545de4b27ee00",
    "recipient": "a77f5cdfa2934df3954a5c7c7da5df1f",
    "amount": 5.0
}'
```

### 3. Mine a New Block
```shell
curl -X GET http://localhost:8080/mine \
  -H "Content-Type: application/json"
```

### 4. Register New Nodes
```shell
curl -X POST http://localhost:8080/nodes/register \
  -H "Content-Type: application/json" \
  -d '{
    "nodes": [
        "http://localhost:5001",
        "http://localhost:5002"
    ]
}'
```

### 5. Resolve Consensus
```shell
curl -X GET http://localhost:8080/nodes/resolve \
  -H "Content-Type: application/json"
```

## Complete flow
```shell
# Register nodes with each other
curl -X POST http://localhost:5001/nodes/register -H "Content-Type: application/json" -d '{"nodes": ["http://localhost:5002", "http://localhost:5003"]}'
curl -X POST http://localhost:5002/nodes/register -H "Content-Type: application/json" -d '{"nodes": ["http://localhost:5001", "http://localhost:5003"]}'
curl -X POST http://localhost:5003/nodes/register -H "Content-Type: application/json" -d '{"nodes": ["http://localhost:5001", "http://localhost:5002"]}'

# Add some transactions on node 1
curl -X POST http://localhost:5001/transactions -H "Content-Type: application/json" -d '{"sender": "node1", "recipient": "node2", "amount": 5.0}'

# Mine on node 1
curl -X GET http://localhost:5001/mine

# Add different transaction on node 2
curl -X POST http://localhost:5002/transactions -H "Content-Type: application/json" -d '{"sender": "node2", "recipient": "node3", "amount": 3.0}'

# Mine on node 2
curl -X GET http://localhost:5002/mine

# Resolve consensus on all nodes
curl -X GET http://localhost:5001/nodes/resolve
curl -X GET http://localhost:5002/nodes/resolve
curl -X GET http://localhost:5003/nodes/resolve

# Check chains are synchronized
curl -X GET http://localhost:5001/chain
curl -X GET http://localhost:5002/chain
curl -X GET http://localhost:5003/chain
```
