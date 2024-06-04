package main

import "errors"

// Node represents a single node in the Merkel Tree.
type Node struct {
	left  *Node
	right *Node
	prev  *Node
	data  []byte
	hash  []byte
}

// CreateNode creates a new Merkel Leaf/branch node.
//   - data can be "nil" if a branch node is being initialized.
func CreateNode(data, hash []byte) (*Node, error) {
	if hash == nil {
		return nil, errors.New("hash data missing")
	}

	return &Node{
		data:  data,
		hash:  hash,
		left:  nil,
		right: nil,
		prev:  nil,
	}, nil
}

func CreateRootBranch(root **Node, data, hash []byte) (*Node, error) {
	currentNode := *root
	branchNode, err := CreateNode([]byte("Y"), []byte(currentNode.hash))
	if err != nil {
		return nil, err
	}
	branchNode.left = currentNode
	right, err := CreateNode(data, hash)
	if err != nil {
		return nil, err
	}
	branchNode.right = right
	currentNode = branchNode
	branchNode.right.prev = branchNode
	branchNode.left.prev = branchNode

	right.prev = branchNode
	*root = branchNode
	return right, nil
}

// InsertNode inserts a node into the merkle tree by;
//  1. Generating a new branch node with no data (represented by the data value "X")
//     This new branch node will replace the previous leaf node that was
//     already present, inheriting its original hash
//  2. Generates a new leaf with the new data. This new leaf will have a new
//     hash created for it.
//  3. The branch node we created in Step 1 will become the new parent of the original leaf node
//     (which will be attached to the right) and the new leaf node we just made at Step 2
//     (attached to the left).
//  4. Lastly, the node that previous pointed to the old leaf node will now point
//     to this new branch node.
func InsertNode(currentNode *Node, prevNodePtr **Node, data, hash []byte) *Node {
	// Create a new branch that will hold our new node and `currentNode`
	// that has data in it.
	newLeaf, err := CreateNode(data, hash)
	if err != nil {
		// absorb the error.
		return currentNode
	}

	newHash := GenerateHash(hash, currentNode.hash)
	newBranch, err := CreateNode([]byte("X"), []byte(newHash))
	if err != nil {
		// absorb the error.
		return currentNode
	}

	newBranch.left = newLeaf
	newBranch.right = currentNode

	// Reorganize prevNodePtr to be in the old CurrentNode's state
	if *prevNodePtr == currentNode {
		*prevNodePtr = newBranch
	} else {
		if (*prevNodePtr).left == currentNode {
			(*prevNodePtr).left = newBranch
			(*prevNodePtr).left.prev = *prevNodePtr
			(*prevNodePtr).right.prev = *prevNodePtr
		} else if (*prevNodePtr).right == currentNode {
			(*prevNodePtr).right = newBranch
			(*prevNodePtr).right.prev = *prevNodePtr
		}
	}

	newLeaf.prev = newBranch

	// TODO: Tidy this logic up. The prev value to the right has its
	// prev value skipped in assignment. I suspect that it's because of
	// the left branch-bias of my logic.
	newLeaf.prev.right.prev = newLeaf.prev.left.prev

	return newLeaf
}
