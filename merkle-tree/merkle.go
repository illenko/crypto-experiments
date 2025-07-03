package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type Node struct {
	Hash   []byte
	Left   *Node
	Right  *Node
	Parent *Node
	Data   []byte
}

type MerkleTree struct {
	Root         *Node
	Leaves       []*Node
	merkleHashes [][]byte
}

type MerkleProof struct {
	LeafIndex int
	LeafHash  []byte
	Siblings  [][]byte
	Path      []bool // true = right, false = left
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	if len(data) == 0 {
		return &MerkleTree{}
	}

	tree := &MerkleTree{}

	// Create leaf nodes
	leaves := make([]*Node, len(data))
	for i, d := range data {
		hash := sha256.Sum256(d)
		leaves[i] = &Node{
			Hash: hash[:],
			Data: d,
		}
	}
	tree.Leaves = leaves

	// Build tree bottom-up
	currentLevel := leaves

	for len(currentLevel) > 1 {
		nextLevel := make([]*Node, 0, (len(currentLevel)+1)/2)

		for i := 0; i < len(currentLevel); i += 2 {
			left := currentLevel[i]
			var right *Node

			if i+1 < len(currentLevel) {
				right = currentLevel[i+1]
			} else {
				// Odd number of nodes - duplicate the last one
				right = &Node{
					Hash: left.Hash,
					Data: left.Data,
				}
			}

			// Create parent node
			combined := append(left.Hash, right.Hash...)
			parentHash := sha256.Sum256(combined)

			parent := &Node{
				Hash:  parentHash[:],
				Left:  left,
				Right: right,
			}

			left.Parent = parent
			right.Parent = parent

			nextLevel = append(nextLevel, parent)
		}

		currentLevel = nextLevel
	}

	if len(currentLevel) > 0 {
		tree.Root = currentLevel[0]
	}

	return tree
}

func (mt *MerkleTree) GenerateProof(data []byte) (*MerkleProof, error) {
	if mt.Root == nil {
		return nil, fmt.Errorf("empty tree")
	}

	targetHash := sha256.Sum256(data)
	leafIndex := -1

	// Find the leaf with matching data
	for i, leaf := range mt.Leaves {
		if string(leaf.Hash) == string(targetHash[:]) {
			leafIndex = i
			break
		}
	}

	if leafIndex == -1 {
		return nil, fmt.Errorf("data not found in tree")
	}

	proof := &MerkleProof{
		LeafIndex: leafIndex,
		LeafHash:  targetHash[:],
		Siblings:  [][]byte{},
		Path:      []bool{},
	}

	current := mt.Leaves[leafIndex]

	// Traverse up the tree collecting sibling hashes and path directions
	for current.Parent != nil {
		parent := current.Parent

		if parent.Left == current {
			// Current is left child, sibling is right
			proof.Siblings = append(proof.Siblings, parent.Right.Hash)
			proof.Path = append(proof.Path, true) // right sibling
		} else {
			// Current is right child, sibling is left
			proof.Siblings = append(proof.Siblings, parent.Left.Hash)
			proof.Path = append(proof.Path, false) // left sibling
		}

		current = parent
	}

	return proof, nil
}

func (mt *MerkleTree) VerifyProof(proof *MerkleProof, rootHash []byte) bool {
	if proof == nil || mt.Root == nil {
		return false
	}

	if string(rootHash) != string(mt.Root.Hash) {
		return false
	}

	currentHash := proof.LeafHash

	// Reconstruct the path to root using sibling hashes
	for i, sibling := range proof.Siblings {
		if proof.Path[i] {
			// Sibling is on the right
			combined := append(currentHash, sibling...)
			hash := sha256.Sum256(combined)
			currentHash = hash[:]
		} else {
			// Sibling is on the left
			combined := append(sibling, currentHash...)
			hash := sha256.Sum256(combined)
			currentHash = hash[:]
		}
	}

	return string(currentHash) == string(rootHash)
}

func VerifyProofStandalone(data []byte, proof *MerkleProof, rootHash []byte) bool {
	leafHash := sha256.Sum256(data)

	if string(leafHash[:]) != string(proof.LeafHash) {
		return false
	}

	currentHash := leafHash[:]

	for i, sibling := range proof.Siblings {
		if proof.Path[i] {
			combined := append(currentHash, sibling...)
			hash := sha256.Sum256(combined)
			currentHash = hash[:]
		} else {
			combined := append(sibling, currentHash...)
			hash := sha256.Sum256(combined)
			currentHash = hash[:]
		}
	}

	return string(currentHash) == string(rootHash)
}

func (mt *MerkleTree) GetRootHash() []byte {
	if mt.Root == nil {
		return nil
	}
	return mt.Root.Hash
}

func (mt *MerkleTree) PrintTree() {
	if mt.Root == nil {
		fmt.Println("Empty tree")
		return
	}

	fmt.Println("Merkle Tree Structure:")
	mt.printNode(mt.Root, "", true)
}

func (mt *MerkleTree) printNode(node *Node, prefix string, isLast bool) {
	if node == nil {
		return
	}

	connector := "├── "
	if isLast {
		connector = "└── "
	}

	hashStr := fmt.Sprintf("%x", node.Hash)
	if len(hashStr) > 16 {
		hashStr = hashStr[:16] + "..."
	}

	dataStr := ""
	if node.Data != nil {
		dataStr = fmt.Sprintf(" [%s]", string(node.Data))
	}

	fmt.Printf("%s%s%s%s\n", prefix, connector, hashStr, dataStr)

	if node.Left != nil || node.Right != nil {
		childPrefix := prefix
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}

		if node.Right != nil {
			mt.printNode(node.Right, childPrefix, node.Left == nil)
		}
		if node.Left != nil {
			mt.printNode(node.Left, childPrefix, true)
		}
	}
}

func (proof *MerkleProof) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Merkle Proof:\n"))
	sb.WriteString(fmt.Sprintf("  Leaf Index: %d\n", proof.LeafIndex))
	sb.WriteString(fmt.Sprintf("  Leaf Hash: %x\n", proof.LeafHash))
	sb.WriteString(fmt.Sprintf("  Proof Length: %d\n", len(proof.Siblings)))

	for i, sibling := range proof.Siblings {
		direction := "left"
		if proof.Path[i] {
			direction = "right"
		}
		sb.WriteString(fmt.Sprintf("  Level %d: %x (%s)\n", i, sibling, direction))
	}

	return sb.String()
}
