package main

import (
	"bytes"
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

	go n.discoverPeers()

	log.Fatal(n.Server.ListenAndServe())
}

func (n *Node) setupRoutes() {
	http.HandleFunc("/health", n.handleHealth)
	http.HandleFunc("/status", n.handleStatus)
	http.HandleFunc("/blockchain", n.handleBlockchain)
	http.HandleFunc("/blockchain/sync", n.handleBlockchainSync)
	http.HandleFunc("/balance/", n.handleBalance)
	http.HandleFunc("/transaction", n.handleTransaction)
	http.HandleFunc("/transaction/broadcast", n.handleTransactionBroadcast)
	http.HandleFunc("/block/broadcast", n.handleBlockBroadcast)
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
	defer n.mutex.RUnlock()

	balance := n.Blockchain.GetBalance(address)

	response := map[string]interface{}{
		"address": address,
		"balance": balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	fmt.Printf("📡 Broadcasting transaction %s to peers...\n", tx.ID[:8])
	go n.broadcastToPeers("/transaction/broadcast", tx)

	response := map[string]interface{}{
		"transactionId": tx.ID,
		"status":        "pending",
		"message":       "Transaction created and broadcasted to peers",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (n *Node) handleTransactionBroadcast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method required", http.StatusMethodNotAllowed)
		return
	}

	var tx Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid transaction JSON", http.StatusBadRequest)
		return
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	for _, existingTx := range n.Blockchain.PendingTransactions {
		if existingTx.ID == tx.ID {
			fmt.Printf("📥 Transaction %s already exists in pending pool\n", tx.ID[:8])
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "already_exists"})
			return
		}
	}

	n.Blockchain.PendingTransactions = append(n.Blockchain.PendingTransactions, &tx)
	fmt.Printf("📥 Received broadcasted transaction %s from peer\n", tx.ID[:8])

	response := map[string]interface{}{
		"status":  "accepted",
		"message": "Transaction added to pending pool",
		"txId":    tx.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (n *Node) handleBlockBroadcast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method required", http.StatusMethodNotAllowed)
		return
	}

	var block Block
	if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
		http.Error(w, "Invalid block JSON", http.StatusBadRequest)
		return
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	if block.Index <= len(n.Blockchain.Chain)-1 {
		fmt.Printf("📥 Block #%d already exists or is outdated\n", block.Index)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "already_exists"})
		return
	}

	if block.Index != len(n.Blockchain.Chain) {
		fmt.Printf("📥 Block #%d is not the next expected block (expected #%d)\n", block.Index, len(n.Blockchain.Chain))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "invalid_sequence"})
		return
	}

	if err := n.Blockchain.SubmitBlock(&block); err != nil {
		fmt.Printf("❌ Failed to accept broadcasted block #%d: %v\n", block.Index, err)
		http.Error(w, "Block validation failed", http.StatusBadRequest)
		return
	}

	fmt.Printf("📥 Accepted broadcasted block #%d from peer\n", block.Index)

	response := map[string]interface{}{
		"status":     "accepted",
		"message":    "Block added to blockchain",
		"blockIndex": block.Index,
		"blockHash":  block.Hash,
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

	fmt.Printf("💎 Node %s: Mining new block...\n", n.ID)
	block := n.Miner.Mine(n.Blockchain)
	if err := n.Blockchain.SubmitBlock(block); err != nil {
		http.Error(w, fmt.Sprintf("Failed to submit block: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("📡 Broadcasting mined block #%d to peers...\n", block.Index)
	go n.broadcastToPeers("/block/broadcast", block)

	response := map[string]interface{}{
		"blockIndex": block.Index,
		"blockHash":  block.Hash,
		"nonce":      block.Nonce,
		"message":    "Block mined and broadcasted to peers",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (n *Node) handlePeers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		n.mutex.RLock()
		defer n.mutex.RUnlock()

		response := map[string]interface{}{
			"nodeId": n.ID,
			"peers":  n.Peers,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		var request struct {
			PeerAddress string `json:"peerAddress"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		n.addPeer(request.PeerAddress)

		response := map[string]interface{}{
			"message": "Peer added successfully",
			"peer":    request.PeerAddress,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (n *Node) addPeer(peerAddress string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for _, peer := range n.Peers {
		if peer == peerAddress {
			fmt.Printf("👥 Peer %s already exists\n", peerAddress)
			return
		}
	}

	n.Peers = append(n.Peers, peerAddress)
	fmt.Printf("👥 Added new peer: %s\n", peerAddress)
}

func (n *Node) discoverPeers() {
	time.Sleep(2 * time.Second)

	for _, peer := range n.Peers {
		go n.connectToPeer(peer)
	}
}

func (n *Node) connectToPeer(peerAddress string) {
	client := &http.Client{Timeout: 5 * time.Second}

	url := fmt.Sprintf("http://%s/health", peerAddress)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("⚠️ Failed to connect to peer %s: %v\n", peerAddress, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Successfully connected to peer: %s\n", peerAddress)

		n.announceToPeer(peerAddress)
	} else {
		fmt.Printf("⚠️ Peer %s responded with status: %d\n", peerAddress, resp.StatusCode)
	}
}

func (n *Node) announceToPeer(peerAddress string) {
	myAddress := fmt.Sprintf("localhost:%d", n.Port)

	announcement := map[string]string{
		"peerAddress": myAddress,
	}

	jsonData, _ := json.Marshal(announcement)

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://%s/peers", peerAddress)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("⚠️ Failed to announce to peer %s: %v\n", peerAddress, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("📢 Successfully announced to peer: %s\n", peerAddress)

		go n.syncWithPeer(peerAddress)
	}
}

func (n *Node) handleBlockchainSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method required", http.StatusMethodNotAllowed)
		return
	}

	var peerChain Blockchain
	if err := json.NewDecoder(r.Body).Decode(&peerChain); err != nil {
		http.Error(w, "Invalid blockchain JSON", http.StatusBadRequest)
		return
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	if peerChain.IsLongerThan(n.Blockchain) && peerChain.IsValidChain() {
		fmt.Printf("🔄 Received longer valid chain (%d blocks vs %d), attempting to replace...\n",
			len(peerChain.Chain), len(n.Blockchain.Chain))

		if err := n.Blockchain.ReplaceChain(&peerChain); err != nil {
			fmt.Printf("❌ Failed to replace chain: %v\n", err)
			http.Error(w, "Chain replacement failed", http.StatusBadRequest)
			return
		}

		response := map[string]interface{}{
			"status":    "chain_replaced",
			"message":   "Blockchain updated with longer chain",
			"newLength": len(n.Blockchain.Chain),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := map[string]interface{}{
			"status":        "chain_not_replaced",
			"message":       "Current chain is longer or peer chain is invalid",
			"currentLength": len(n.Blockchain.Chain),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func (n *Node) syncWithPeer(peerAddress string) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("http://%s/blockchain", peerAddress)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("⚠️ Failed to get blockchain from peer %s: %v\n", peerAddress, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("⚠️ Peer %s returned status %d for blockchain request\n", peerAddress, resp.StatusCode)
		return
	}

	var peerChain Blockchain
	if err := json.NewDecoder(resp.Body).Decode(&peerChain); err != nil {
		fmt.Printf("⚠️ Failed to decode blockchain from peer %s: %v\n", peerAddress, err)
		return
	}

	n.mutex.Lock()
	peerLonger := peerChain.IsLongerThan(n.Blockchain)
	currentLength := len(n.Blockchain.Chain)
	peerLength := len(peerChain.Chain)
	n.mutex.Unlock()

	if peerLonger {
		fmt.Printf("🔍 Peer %s has longer chain (%d vs %d), requesting sync...\n",
			peerAddress, peerLength, currentLength)

		jsonData, _ := json.Marshal(peerChain)
		syncURL := fmt.Sprintf("http://localhost:%d/blockchain/sync", n.Port)

		syncResp, err := client.Post(syncURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("⚠️ Failed to sync with peer chain: %v\n", err)
			return
		}
		defer syncResp.Body.Close()

		if syncResp.StatusCode == http.StatusOK {
			fmt.Printf("🔄 Successfully synced with peer %s\n", peerAddress)
		}
	} else {
		fmt.Printf("ℹ️ Peer %s has same or shorter chain (%d vs %d)\n",
			peerAddress, peerLength, currentLength)
	}
}

func (n *Node) broadcastToPeers(endpoint string, data interface{}) {
	jsonData, _ := json.Marshal(data)

	for _, peer := range n.Peers {
		go func(peerAddr string) {
			client := &http.Client{Timeout: 5 * time.Second}
			url := fmt.Sprintf("http://%s%s", peerAddr, endpoint)

			resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Printf("⚠️ Failed to broadcast to peer %s: %v\n", peerAddr, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				fmt.Printf("📡 Broadcast successful to peer: %s\n", peerAddr)
			}
		}(peer)
	}
}
