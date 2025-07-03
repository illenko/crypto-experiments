package main

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestNewMerkleTree(t *testing.T) {
	tests := []struct {
		name               string
		data               [][]byte
		expectedLeaves     int
		expectedRootNil    bool
		expectedRootHash   []byte
		shouldHaveChildren bool
	}{
		{
			name:            "empty data",
			data:            [][]byte{},
			expectedLeaves:  0,
			expectedRootNil: true,
		},
		{
			name:               "single element",
			data:               [][]byte{[]byte("transaction1")},
			expectedLeaves:     1,
			expectedRootNil:    false,
			expectedRootHash:   func() []byte { h := sha256.Sum256([]byte("transaction1")); return h[:] }(),
			shouldHaveChildren: false,
		},
		{
			name: "power of two",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
				[]byte("tx4"),
			},
			expectedLeaves:     4,
			expectedRootNil:    false,
			shouldHaveChildren: true,
		},
		{
			name: "odd number",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
			},
			expectedLeaves:     3,
			expectedRootNil:    false,
			shouldHaveChildren: true,
		},
		{
			name: "two elements",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
			},
			expectedLeaves:     2,
			expectedRootNil:    false,
			shouldHaveChildren: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := NewMerkleTree(tt.data)

			if tt.expectedRootNil {
				if tree.Root != nil {
					t.Error("Expected empty tree to have nil root")
				}
				return
			}

			if tree.Root == nil {
				t.Fatal("Expected tree to have a root")
			}

			if len(tree.Leaves) != tt.expectedLeaves {
				t.Errorf("Expected %d leaves, got %d", tt.expectedLeaves, len(tree.Leaves))
			}

			if tt.expectedRootHash != nil {
				if string(tree.Root.Hash) != string(tt.expectedRootHash) {
					t.Error("Root hash does not match expected hash")
				}
			}

			if tt.shouldHaveChildren {
				if tree.Root.Left == nil || tree.Root.Right == nil {
					t.Error("Root should have both left and right children")
				}
			}
		})
	}
}

func TestGenerateProof(t *testing.T) {
	tests := []struct {
		name              string
		data              [][]byte
		target            []byte
		expectedLeafIndex int
		expectError       bool
		errorMessage      string
	}{
		{
			name: "valid data - first element",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
				[]byte("tx4"),
			},
			target:            []byte("tx1"),
			expectedLeafIndex: 0,
			expectError:       false,
		},
		{
			name: "valid data - middle element",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
				[]byte("tx4"),
			},
			target:            []byte("tx2"),
			expectedLeafIndex: 1,
			expectError:       false,
		},
		{
			name: "valid data - last element",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
			},
			target:            []byte("tx3"),
			expectedLeafIndex: 2,
			expectError:       false,
		},
		{
			name: "invalid data - nonexistent",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
			},
			target:       []byte("nonexistent"),
			expectError:  true,
			errorMessage: "data not found in tree",
		},
		{
			name:         "empty tree",
			data:         [][]byte{},
			target:       []byte("tx1"),
			expectError:  true,
			errorMessage: "empty tree",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := NewMerkleTree(tt.data)

			proof, err := tree.GenerateProof(tt.target)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errorMessage != "" && err != nil && err.Error() != tt.errorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMessage, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if proof.LeafIndex != tt.expectedLeafIndex {
				t.Errorf("Expected leaf index %d, got %d", tt.expectedLeafIndex, proof.LeafIndex)
			}

			expectedHash := sha256.Sum256(tt.target)
			if string(proof.LeafHash) != string(expectedHash[:]) {
				t.Error("Proof leaf hash does not match expected")
			}
		})
	}
}

func TestVerifyProof(t *testing.T) {
	tests := []struct {
		name           string
		data           [][]byte
		targetTx       []byte
		modifyProof    func(*MerkleProof)
		modifyRoot     func([]byte) []byte
		verifyData     []byte
		expectedValid  bool
		testStandalone bool
		description    string
	}{
		{
			name: "valid proof - tree method",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
				[]byte("tx4"),
			},
			targetTx:      []byte("tx3"),
			verifyData:    []byte("tx3"),
			expectedValid: true,
			description:   "Valid proof should be accepted",
		},
		{
			name: "corrupted proof - tree method",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
				[]byte("tx4"),
			},
			targetTx: []byte("tx1"),
			modifyProof: func(proof *MerkleProof) {
				if len(proof.Siblings) > 0 {
					proof.Siblings[0][0] ^= 0xFF
				}
			},
			verifyData:    []byte("tx1"),
			expectedValid: false,
			description:   "Corrupted proof should be rejected",
		},
		{
			name: "wrong root hash - tree method",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
			},
			targetTx: []byte("tx1"),
			modifyRoot: func(root []byte) []byte {
				modified := make([]byte, len(root))
				copy(modified, root)
				modified[0] ^= 0xFF
				return modified
			},
			verifyData:    []byte("tx1"),
			expectedValid: false,
			description:   "Wrong root hash should be rejected",
		},
		{
			name: "valid proof - standalone method",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
				[]byte("tx4"),
			},
			targetTx:       []byte("tx2"),
			verifyData:     []byte("tx2"),
			expectedValid:  true,
			testStandalone: true,
			description:    "Standalone verification should work for valid proof",
		},
		{
			name: "wrong data - standalone method",
			data: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
				[]byte("tx4"),
			},
			targetTx:       []byte("tx2"),
			verifyData:     []byte("wrong"),
			expectedValid:  false,
			testStandalone: true,
			description:    "Standalone verification should reject wrong data",
		},
		{
			name: "single element tree",
			data: [][]byte{
				[]byte("onlytx"),
			},
			targetTx:       []byte("onlytx"),
			verifyData:     []byte("onlytx"),
			expectedValid:  true,
			testStandalone: true,
			description:    "Single element tree should verify correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := NewMerkleTree(tt.data)
			rootHash := tree.GetRootHash()

			proof, err := tree.GenerateProof(tt.targetTx)
			if err != nil {
				t.Fatalf("Failed to generate proof: %v", err)
			}

			// Apply proof modifications if specified
			if tt.modifyProof != nil {
				tt.modifyProof(proof)
			}

			// Apply root modifications if specified
			verifyRoot := rootHash
			if tt.modifyRoot != nil {
				verifyRoot = tt.modifyRoot(rootHash)
			}

			var valid bool
			if tt.testStandalone {
				valid = VerifyProofStandalone(tt.verifyData, proof, verifyRoot)
			} else {
				valid = tree.VerifyProof(proof, verifyRoot)
			}

			if valid != tt.expectedValid {
				t.Errorf("%s: expected %t, got %t", tt.description, tt.expectedValid, valid)
			}
		})
	}
}

func TestLargeTree(t *testing.T) {
	// Test with 1000 transactions
	data := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = []byte(fmt.Sprintf("transaction_%d", i))
	}

	tree := NewMerkleTree(data)

	if tree.Root == nil {
		t.Fatal("Large tree should have a root")
	}

	// Test proof generation and verification for middle element
	testTx := []byte("transaction_500")
	proof, err := tree.GenerateProof(testTx)
	if err != nil {
		t.Fatalf("Failed to generate proof for large tree: %v", err)
	}

	valid := VerifyProofStandalone(testTx, proof, tree.GetRootHash())
	if !valid {
		t.Error("Proof verification failed for large tree")
	}

	// Proof should be logarithmic in size
	maxProofSize := 20 // log2(1000) ≈ 10, with some buffer
	if len(proof.Siblings) > maxProofSize {
		t.Errorf("Proof size too large: %d (expected < %d)", len(proof.Siblings), maxProofSize)
	}
}

func TestConsistentHashing(t *testing.T) {
	data := [][]byte{
		[]byte("tx1"),
		[]byte("tx2"),
		[]byte("tx3"),
		[]byte("tx4"),
	}

	tree1 := NewMerkleTree(data)
	tree2 := NewMerkleTree(data)

	if string(tree1.GetRootHash()) != string(tree2.GetRootHash()) {
		t.Error("Same data should produce same root hash")
	}
}

func BenchmarkMerkleTree_Construction(b *testing.B) {
	data := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = []byte(fmt.Sprintf("transaction_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewMerkleTree(data)
	}
}

func BenchmarkMerkleTree_ProofGeneration(b *testing.B) {
	data := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = []byte(fmt.Sprintf("transaction_%d", i))
	}

	tree := NewMerkleTree(data)
	testData := []byte("transaction_500")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GenerateProof(testData)
	}
}

func BenchmarkMerkleTree_ProofVerification(b *testing.B) {
	data := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = []byte(fmt.Sprintf("transaction_%d", i))
	}

	tree := NewMerkleTree(data)
	testData := []byte("transaction_500")
	proof, _ := tree.GenerateProof(testData)
	rootHash := tree.GetRootHash()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyProofStandalone(testData, proof, rootHash)
	}
}
