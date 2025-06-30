package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	NodeAddress string
	HTTPClient  *http.Client
}

func NewClient(nodeAddress string) *Client {
	return &Client{
		NodeAddress: nodeAddress,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetBalance(address string) {
	url := fmt.Sprintf("http://%s/balance/%s", c.NodeAddress, address)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		fmt.Printf("❌ Failed to connect to node: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Error: %s\n", string(body))
		return
	}

	var response struct {
		Address string  `json:"address"`
		Balance float64 `json:"balance"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		return
	}

	fmt.Printf("💰 Balance for %s: %.2f coins\n", response.Address, response.Balance)
}

func (c *Client) SendTransaction(from, to string, amount, fee float64) {
	url := fmt.Sprintf("http://%s/transaction", c.NodeAddress)

	txRequest := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
		"fee":    fee,
	}

	jsonData, _ := json.Marshal(txRequest)

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Failed to connect to node: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Transaction failed: %s\n", string(body))
		return
	}

	var response struct {
		TransactionID string `json:"transactionId"`
		Status        string `json:"status"`
		Message       string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		return
	}

	fmt.Printf("✅ Transaction created:\n")
	fmt.Printf("   ID: %s\n", response.TransactionID[:16]+"...")
	fmt.Printf("   Status: %s\n", response.Status)
	fmt.Printf("   From: %s\n", from)
	fmt.Printf("   To: %s\n", to)
	fmt.Printf("   Amount: %.2f coins\n", amount)
	fmt.Printf("   Fee: %.2f coins\n", fee)
}

func (c *Client) Mine() {
	url := fmt.Sprintf("http://%s/mine", c.NodeAddress)

	resp, err := c.HTTPClient.Post(url, "application/json", nil)
	if err != nil {
		fmt.Printf("❌ Failed to connect to node: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Mining failed: %s\n", string(body))
		return
	}

	var response struct {
		BlockIndex int    `json:"blockIndex"`
		BlockHash  string `json:"blockHash"`
		Nonce      int    `json:"nonce"`
		Message    string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		return
	}

	fmt.Printf("⛏️ Mining successful:\n")
	fmt.Printf("   Block: #%d\n", response.BlockIndex)
	fmt.Printf("   Hash: %s...\n", response.BlockHash[:16])
	fmt.Printf("   Nonce: %d\n", response.Nonce)
}

func (c *Client) GetStatus() {
	url := fmt.Sprintf("http://%s/status", c.NodeAddress)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		fmt.Printf("❌ Failed to connect to node: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Error: %s\n", string(body))
		return
	}

	var response struct {
		NodeID      string   `json:"nodeId"`
		Port        int      `json:"port"`
		Peers       []string `json:"peers"`
		BlockCount  int      `json:"blockCount"`
		Difficulty  int      `json:"difficulty"`
		PendingTxs  int      `json:"pendingTxs"`
		MinerWallet string   `json:"minerWallet"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		return
	}

	fmt.Printf("📊 Node Status:\n")
	fmt.Printf("   Node ID: %s\n", response.NodeID)
	fmt.Printf("   Port: %d\n", response.Port)
	fmt.Printf("   Blocks: %d\n", response.BlockCount)
	fmt.Printf("   Difficulty: %d\n", response.Difficulty)
	fmt.Printf("   Pending Transactions: %d\n", response.PendingTxs)
	fmt.Printf("   Miner Wallet: %s\n", response.MinerWallet)
	if len(response.Peers) > 0 {
		fmt.Printf("   Peers: %v\n", response.Peers)
	} else {
		fmt.Printf("   Peers: none\n")
	}
}
