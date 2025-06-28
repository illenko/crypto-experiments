package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Node struct {
	ID         string
	Port       int
	Blockchain *Blockchain
	Wallets    *Wallets
	Miner      *Miner
	Peers      []string
	Server     *http.Server
	mutex      sync.RWMutex
}

func NewNode(port int, peers string) *Node {
	nodeID := fmt.Sprintf("node-%d", port)

	wallets := NewWallets()
	minerWallet := wallets.CreateWallet()

	blockchain := NewBlockchain()
	miner := NewMiner(minerWallet)

	var peerList []string
	if peers != "" {
		peerList = strings.Split(peers, ",")
		for i, peer := range peerList {
			peerList[i] = strings.TrimSpace(peer)
		}
	}

	node := &Node{
		ID:         nodeID,
		Port:       port,
		Blockchain: blockchain,
		Wallets:    wallets,
		Miner:      miner,
		Peers:      peerList,
	}

	return node
}

func (n *Node) Start() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	fmt.Printf("🔗 Starting node %s on port %d\n", n.ID, n.Port)
	fmt.Printf("⛏️  Miner wallet: %s\n", n.Miner.Address)

	if len(n.Peers) > 0 {
		fmt.Printf("👥 Peers: %v\n", n.Peers)
	}

	n.setupRoutes()

	addr := fmt.Sprintf(":%d", n.Port)
	n.Server = &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("🚀 Node %s ready and listening on %s\n", n.ID, addr)
	log.Fatal(n.Server.ListenAndServe())
}

func (n *Node) setupRoutes() {
	http.HandleFunc("/health", n.handleHealth)
	http.HandleFunc("/status", n.handleStatus)
	http.HandleFunc("/blockchain", n.handleBlockchain)
	http.HandleFunc("/balance/", n.handleBalance)
	http.HandleFunc("/transaction", n.handleTransaction)
	http.HandleFunc("/mine", n.handleMine)
	http.HandleFunc("/peers", n.handlePeers)
}

func (n *Node) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "healthy",
		"nodeId": n.ID,
		"port":   n.Port,
		"time":   time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (n *Node) handleStatus(w http.ResponseWriter, r *http.Request) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	response := map[string]interface{}{
		"nodeId":      n.ID,
		"port":        n.Port,
		"peers":       n.Peers,
		"blockCount":  len(n.Blockchain.Chain),
		"difficulty":  DIFFICULTY,
		"pendingTxs":  len(n.Blockchain.PendingTransactions),
		"minerWallet": n.Miner.Address,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (n *Node) handleBlockchain(w http.ResponseWriter, r *http.Request) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n.Blockchain)
}

func (n *Node) handleBalance(w http.ResponseWriter, r *http.Request) {
	address := strings.TrimPrefix(r.URL.Path, "/balance/")
	if address == "" {
		http.Error(w, "Address required", http.StatusBadRequest)
		return
	}

	n.mutex.RLock()
	balance := n.calculateBalance(address)
	n.mutex.RUnlock()

	response := map[string]interface{}{
		"address": address,
		"balance": balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (n *Node) calculateBalance(address string) float64 {
	balance := 0.0
	utxos := n.Blockchain.UTXOSet[address]
	for _, utxo := range utxos {
		balance += utxo.Output.Value
	}
	return balance
}

func (n *Node) handleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method required", http.StatusMethodNotAllowed)
		return
	}

	var txRequest struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
		Fee    float64 `json:"fee"`
	}

	if err := json.NewDecoder(r.Body).Decode(&txRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	n.mutex.Lock()
	tx := n.Blockchain.CreateTransaction(txRequest.From, txRequest.To, txRequest.Amount, txRequest.Fee)
	n.mutex.Unlock()

	if tx == nil {
		http.Error(w, "Failed to create transaction", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"transactionId": tx.ID,
		"status":        "pending",
		"message":       "Transaction created and added to pending pool",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (n *Node) handleMine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method required", http.StatusMethodNotAllowed)
		return
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	if len(n.Blockchain.PendingTransactions) == 0 {
		http.Error(w, "No pending transactions to mine", http.StatusBadRequest)
		return
	}

	fmt.Printf("💎 Node %s: Mining new block...\n", n.ID)
	block := n.Miner.Mine(n.Blockchain)
	n.Blockchain.SubmitBlock(block)

	response := map[string]interface{}{
		"blockIndex": block.Index,
		"blockHash":  block.Hash,
		"nonce":      block.Nonce,
		"message":    "Block mined successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (n *Node) handlePeers(w http.ResponseWriter, r *http.Request) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	response := map[string]interface{}{
		"nodeId": n.ID,
		"peers":  n.Peers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
