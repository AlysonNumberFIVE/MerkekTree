package main

import (
	"errors"
)

type MerkelProof struct {
	LeafHash   []byte
	ProofList  [][]byte
	Directions []bool
}

// GenerateProof creates a new merkel proof. The logic takes the root node's hash
// and attempts to traverse up the tree from the leaf, building up the leaf's hash
// until it reaches root.
func (merkelTree *MerkelTree) GenerateProof(leafHash []byte) (*MerkelProof, error) {
	if merkelTree.root == nil {
		return nil, errors.New("No root")
	}

	node, err := merkelTree.Lookup(leafHash)
	if err != nil {
		return nil, errors.New("Hash not found")
	}

	proofChain := [][]byte{node.hash}
	pathway := []bool{true}
	current := node
	previous := node.prev
	for previous != nil {
		if previous.left == current {
			proofChain = append(proofChain, previous.right.hash)
			pathway = append(pathway, true)
		} else {

			proofChain = append(proofChain, previous.left.hash)
			pathway = append(pathway, false)
		}
		current = previous
		previous = previous.prev
	}

	return &MerkelProof{
		LeafHash:   leafHash,
		ProofList:  proofChain,
		Directions: pathway,
	}, nil
}

// VerifyProof verifies that a merkel proof is valid and can
// be used to rebuild root's hash.
func VerifyProof(proof *MerkelProof, rootHash []byte) bool {
	var value []byte

	if proof == nil {
		return false
	}

	for index, hashPiece := range proof.ProofList {
		// Right join
		if proof.Directions[index] {
			value = GenerateHash(value, hashPiece)
			// Left  join
		} else {
			value = GenerateHash(hashPiece, value)
		}
	}

	return compareHash(value, rootHash)
}
