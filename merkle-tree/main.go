package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Merkle Tree Blockchain Example ===")
	fmt.Println()

	transactions := [][]byte{
		[]byte("Alice sends 10 BTC to Bob"),
		[]byte("Bob sends 5 BTC to Charlie"),
		[]byte("Charlie sends 3 BTC to David"),
		[]byte("David sends 1 BTC to Eve"),
		[]byte("Eve sends 2 BTC to Frank"),
		[]byte("Frank sends 4 BTC to Grace"),
		[]byte("Grace sends 6 BTC to Henry"),
		[]byte("Henry sends 8 BTC to Alice"),
	}

	fmt.Printf("Creating Merkle tree from %d transactions...\n", len(transactions))

	start := time.Now()
	tree := NewMerkleTree(transactions)
	constructionTime := time.Since(start)

	fmt.Printf("Tree construction took: %v\n\n", constructionTime)

	tree.PrintTree()
	fmt.Println()

	fmt.Println("=== Merkle Proof Demonstration ===")
	targetTx := []byte("Charlie sends 3 BTC to David")
	fmt.Printf("Generating proof for transaction: %s\n\n", string(targetTx))

	start = time.Now()
	proof, err := tree.GenerateProof(targetTx)
	proofGenTime := time.Since(start)

	if err != nil {
		fmt.Printf("Error generating proof: %v\n", err)
		return
	}

	fmt.Printf("Proof generation took: %v\n", proofGenTime)
	fmt.Println(proof.String())

	start = time.Now()
	rootHash := tree.GetRootHash()
	isValid := VerifyProofStandalone(targetTx, proof, rootHash)
	verificationTime := time.Since(start)

	fmt.Printf("Proof verification took: %v\n", verificationTime)
	fmt.Printf("Proof is valid: %t\n\n", isValid)

	fmt.Println("=== Efficiency Comparison ===")

	largeBlockchain := make([][]byte, 1000000)
	for i := 0; i < 1000000; i++ {
		largeBlockchain[i] = []byte(fmt.Sprintf("Transaction_%d: User_%d sends %d coins to User_%d",
			i, i%1000, (i%100)+1, (i+1)%1000))
	}

	fmt.Printf("Creating Merkle tree from %d transactions...\n", len(largeBlockchain))
	start = time.Now()
	largeTree := NewMerkleTree(largeBlockchain)
	largeBuildTime := time.Since(start)
	fmt.Printf("Large tree construction took: %v\n", largeBuildTime)

	targetLargeTx := largeBlockchain[500000]

	start = time.Now()
	largeProof, err := largeTree.GenerateProof(targetLargeTx)
	largeProofGenTime := time.Since(start)

	if err != nil {
		fmt.Printf("Error generating proof for large tree: %v\n", err)
		return
	}

	fmt.Printf("Proof generation for large tree took: %v\n", largeProofGenTime)
	fmt.Printf("Proof size: %d hashes (log₂(%d) ≈ %d)\n",
		len(largeProof.Siblings), len(largeBlockchain), approximateLog2(len(largeBlockchain)))

	start = time.Now()
	largeRootHash := largeTree.GetRootHash()
	largeIsValid := VerifyProofStandalone(targetLargeTx, largeProof, largeRootHash)
	largeVerificationTime := time.Since(start)

	fmt.Printf("Proof verification for large tree took: %v\n", largeVerificationTime)
	fmt.Printf("Large proof is valid: %t\n\n", largeIsValid)

	fmt.Println("=== Naive vs Merkle Tree Verification ===")

	start = time.Now()
	naiveValid := naiveVerification(targetLargeTx, largeBlockchain)
	naiveTime := time.Since(start)

	fmt.Printf("Naive verification (O(n)): %v, valid: %t\n", naiveTime, naiveValid)
	fmt.Printf("Merkle verification (O(log n)): %v, valid: %t\n", largeVerificationTime, largeIsValid)
	fmt.Printf("Speedup: %.2fx\n\n", float64(naiveTime)/float64(largeVerificationTime))

	fmt.Println("=== Security Demonstration ===")

	tamperedTx := []byte("Charlie sends 300 BTC to David") // Changed amount
	tamperedValid := VerifyProofStandalone(tamperedTx, proof, rootHash)
	fmt.Printf("Tampered transaction verification: %t (should be false)\n", tamperedValid)

	// Try to use proof with wrong root hash
	wrongRoot := make([]byte, len(rootHash))
	copy(wrongRoot, rootHash)
	wrongRoot[0] ^= 0xFF // Flip some bits

	wrongRootValid := VerifyProofStandalone(targetTx, proof, wrongRoot)
	fmt.Printf("Valid transaction with wrong root: %t (should be false)\n", wrongRootValid)

	fmt.Println("\n=== Summary ===")
	fmt.Printf("✓ Merkle trees enable O(log n) verification vs O(n) naive approach\n")
	fmt.Printf("✓ Proof size grows logarithmically with blockchain size\n")
	fmt.Printf("✓ Cryptographic security prevents tampering\n")
	fmt.Printf("✓ Enables light clients to verify transactions without downloading entire blockchain\n")
}

func naiveVerification(target []byte, transactions [][]byte) bool {
	for _, tx := range transactions {
		if string(tx) == string(target) {
			return true
		}
	}
	return false
}

func approximateLog2(n int) int {
	log := 0
	for n > 1 {
		n /= 2
		log++
	}
	return log
}
