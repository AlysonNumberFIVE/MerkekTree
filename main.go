package main

import (
	"errors"
	"fmt"
	"log"
)

// Mapping is used to help with managing and maintaining search
// functionality. When a hash is updated, it isn't overwritten,
// version history of the hashes are kept inside of hashUpdateHistroy
// to make querying a node with an "out-of-date" hash still possible.
type Mapping struct {
	node              *Node
	hashUpdateHistroy [][]byte
}

// MerkelTree holds the root node as well as a list for easy lookup for
// searching/updating nodes.
type MerkelTree struct {
	root           *Node
	lookupNodeList map[string]*Mapping
}

type NodeDepth struct {
	depth int
	node  *Node
}

// NavigateTree scans the entire merkel tree and returns all the leaf nodes
// along with their respective depths. This is very important for determining
// which branch will be appended to when a new record is created to
// keep th tree balanced
//
//	   i.e
//	This tree:
//	          O
//	        /   \
//	       O     C
//	      / \
//	     A   B
//
//	Returns:
//	   [A, 2], [B, 2], [C, 1]
//
//	And so C will be prioritized when adding new data.
func (merkelTree *MerkelTree) navigateTree() []NodeDepth {
	if merkelTree.root == nil {
		return nil
	}
	traversalNode := merkelTree.root
	leafDepths := []NodeDepth{NodeDepth{depth: 0, node: traversalNode}}
	depths := []NodeDepth{}

	for len(leafDepths) > 0 {
		currentNodeDepth := leafDepths[len(leafDepths)-1]
		leafDepths = leafDepths[:len(leafDepths)-1]
		traversalNode = currentNodeDepth.node

		if traversalNode.right == nil && traversalNode.left == nil {
			depths = append(depths, currentNodeDepth)
		}
		if traversalNode.right != nil {
			leafDepths = append(leafDepths, NodeDepth{node: traversalNode.right, depth: currentNodeDepth.depth + 1})
		}
		if traversalNode.left != nil {
			leafDepths = append(leafDepths, NodeDepth{node: traversalNode.left, depth: currentNodeDepth.depth + 1})
		}
	}

	return depths
}

// InitMerkelTree initializes a new Merkel Tree. Init function is standalone
// in case a hash function needs to be changed/replaced.
func InitMerkelTree() *MerkelTree {
	return &MerkelTree{
		root:           nil,
		lookupNodeList: map[string]*Mapping{},
	}
}

// findHash determines if a hash exists.
//
//	NOTE: The nature of this tree is specificity. The logic in this search
//	function is akin to the logic in Lookup but is designed
//	solely for searching. Lookup must only be used for querying nodes
func (merkelTree *MerkelTree) findHash(hash []byte) bool {
	_, ok := merkelTree.lookupNodeList[string(hash)]

	if !ok {
		// Regular lookup failed. We
		for _, updateHistory := range merkelTree.lookupNodeList {
			for _, updatedHash := range updateHistory.hashUpdateHistroy {
				if compareHash(hash, updatedHash) {
					return true
				}
			}
		}
		return false
	}

	return true
}

// Lookup scans the hashmap of the merkel tree and returns the existing node.
func (merkelTree *MerkelTree) Lookup(hash []byte) (*Node, error) {
	block, ok := merkelTree.lookupNodeList[string(hash)]

	if !ok {
		// Regular lookup failed. We
		for oldHash, updateHistory := range merkelTree.lookupNodeList {

			for _, updatedHash := range updateHistory.hashUpdateHistroy {
				if compareHash(hash, updatedHash) {
					block, ok := merkelTree.lookupNodeList[oldHash]
					if ok {
						return block.node, nil
					}
				}
			}
		}

		return nil, errors.New(fmt.Sprintf("hash (%v) not found", hash))
	}

	return block.node, nil
}

// newHash creates a new hash for the specified node.
//
//	IMPORTANT: HashUpdateHistory is only accessible when
//	an existing node data object is being interacted with.
func (merkelTree *MerkelTree) newHash(node *Node, hash []byte) error {
	if merkelTree.findHash(hash) {
		return errors.New("Hash already exists. Use Update() to update an existing hash")
	}

	merkelTree.lookupNodeList[string(hash)] = &Mapping{
		node:              node,
		hashUpdateHistroy: [][]byte{},
	}

	return nil
}

// Insert creates a new leaf node and adds it to the merkel tree, expanding
// It's behaviour is designed to accommodate the following conditions uniquely
//
//  1. A new node is created on an empty merkel tree
//  2. A new child node is created, shifting the initial root node to
//     being a leaf alongside the new child node just created.
//  3. Every new insert after the initial 2 unique cases.
func (merkelTree *MerkelTree) Insert(data []byte) ([]byte, error) {
	hash := Hash128(data)

	// First check if this hash exists
	if merkelTree.findHash(hash) {
		return nil, errors.New("Hash already exists. Use Update() to update an existing hash")
	}

	var newNode *Node
	// If the root is nil, make a new node
	if merkelTree.root == nil {
		newNode, err := CreateNode(data, hash)
		if err != nil {
			return nil, err
		}
		err = merkelTree.newHash(newNode, hash)
		if err != nil {
			return nil, err
		}
		merkelTree.root = newNode
		// If we're at the first node, initialize it's children
	} else if merkelTree.root.left == nil && merkelTree.root.right == nil {
		newNode, err := CreateRootBranch(&merkelTree.root, data, hash)
		if err != nil {
			return nil, err
		}

		traverse := newNode.prev
		for traverse != nil {
			traverse.hash = GenerateHash(traverse.left.hash, traverse.right.hash)
			traverse = traverse.prev
		}
		merkelTree.newHash(newNode, hash)
		// Scenario after first and second inserts. Find all leaf heights and
		// only start adding nodes to the shallowest to ensure we prioritize
		// evening out/leveling the heights to keep our tree balanced.
	} else {

		// Get all the deepest nodes and get the first shallow node you find.
		deepestNodes := merkelTree.navigateTree()
		targetNode := deepestNodes[0]
		for _, node := range deepestNodes {
			if node.depth < targetNode.depth {
				targetNode = node
			}
		}

		// Insert at this shallow node.
		newNode = InsertNode(targetNode.node, &targetNode.node.prev, data, hash)
		traverse := newNode.prev
		for traverse != nil {
			traverse.hash = GenerateHash(traverse.left.hash, traverse.right.hash)
			traverse = traverse.prev
		}
		merkelTree.newHash(newNode, hash)
	}

	//	merkelTree.NewHash(newNode, hash)
	// Update all the parent hashs all the way up the stack.

	return hash, nil
}

// updateHashVersionHistory adds the newly created hash to the old hash's hash chain.
func (merkelTree *MerkelTree) updateHashVersionHistory(oldHash, newHash []byte) error {
	_, ok := merkelTree.lookupNodeList[string(oldHash)]
	if !ok {
		for hash, node := range merkelTree.lookupNodeList {
			for _, hashHistory := range node.hashUpdateHistroy {
				if compareHash(oldHash, hashHistory) {
					oldHash = []byte(hash)
					break
				}
			}
		}

	}
	// This check should absolutely never fail, so if it does, foul play is likely
	// a cause; i.e a forged/duplicate hash.
	for _, hash := range merkelTree.lookupNodeList[string(oldHash)].hashUpdateHistroy {
		if compareHash(newHash, hash) {
			return errors.New("duplicate hash update insert detected")
		}
	}

	merkelTree.lookupNodeList[string(oldHash)].hashUpdateHistroy = append(
		merkelTree.lookupNodeList[string(oldHash)].hashUpdateHistroy, newHash,
	)

	return nil
}

// Update takes in the old hash with new data and replaces it.
func (merkelTree *MerkelTree) Update(newData, hash []byte) ([]byte, error) {
	node, err := merkelTree.Lookup(hash)
	if err != nil {
		return nil, errors.New("Hash not found")
	}

	newHash := Hash128(newData)
	node.data = newData
	node.hash = newHash
	node = node.prev

	for node != nil {
		node.hash = GenerateHash(node.left.hash, node.right.hash)
		node = node.prev
	}

	merkelTree.updateHashVersionHistory(hash, newHash)

	return newHash, nil
}

// Visualizer is the MerkelTree version of treeDebug. As an endpoint, this seems
// useful to have implemented
func (merkelTree *MerkelTree) Visualizer(node *Node, prefix string, isLeft bool) {
	if node == nil {
		return
	}

	if node.right != nil {
		newPrefix := prefix
		if isLeft {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		merkelTree.Visualizer(node.right, newPrefix, false)
	}

	fmt.Printf("%s", prefix)
	if isLeft {
		fmt.Printf("├── ")
	} else {
		fmt.Printf("└── ")
	}
	fmt.Printf("%s\n", node.data)

	if node.left != nil {
		newPrefix := prefix
		if isLeft {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		merkelTree.Visualizer(node.left, newPrefix, true)
	}
}

func main() {

	tree := InitMerkelTree()
	fmt.Println("First -- insert A")
	hashA, err := tree.Insert([]byte("A"))
	if err != nil {
		log.Fatal(err)
	}
	tree.Visualizer(tree.root, "", false)

	fmt.Println("Second -- insert B")
	hashB, err := tree.Insert([]byte("B"))
	if err != nil {
		log.Fatal(err)
	}
	tree.Visualizer(tree.root, "", false)

	fmt.Println("Third -- insert duplicate A")
	_, err = tree.Insert([]byte("B"))
	fmt.Printf("Error : %+v\n", err)
	tree.Visualizer(tree.root, "", false)

	fmt.Println("Fourth -- insert C, D, E")
	tree.Insert([]byte("C"))
	tree.Insert([]byte("D"))
	tree.Insert([]byte("E"))
	tree.Visualizer(tree.root, "", false)

	fmt.Println("Fifth -- update A to F")
	hashF, err := tree.Update([]byte("F"), hashA)
	if err != nil {
		log.Fatal(err)
	}
	tree.Visualizer(tree.root, "", false)

	fmt.Println("Sixth -- update F to A again (use A's old hash)")
	tree.Update([]byte("A"), hashA)
	tree.Visualizer(tree.root, "", false)

	fmt.Println("Seventh -- update A to M (use F's hash this time)")
	_, err = tree.Update([]byte("M"), hashF)
	if err != nil {
		log.Fatal(err)
	}
	tree.Visualizer(tree.root, "", false)

	fmt.Println("Eighth -- lookup M using A's hash")
	nodeM, err := tree.Lookup(hashA)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data is ", string(nodeM.data))

	fmt.Println("Ninth -- lookup M using F's hash")
	nodeM, err = tree.Lookup(hashF)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data is ", string(nodeM.data))

	nodeD, err := tree.Lookup(Hash128([]byte("D")))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data is ", string(nodeD.data))

	proofA, err := tree.GenerateProof(hashF)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("proofA is ", VerifyProof(proofA, tree.root.hash))

	proofA, err = tree.GenerateProof(hashA)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("proofA is ", VerifyProof(proofA, tree.root.hash))

	proofNull, _ := tree.GenerateProof(Hash128([]byte("4444")))
	fmt.Println("proofNull is ", VerifyProof(proofNull, tree.root.hash))

	proofB, err := tree.GenerateProof(hashB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("proofB is ", VerifyProof(proofB, tree.root.hash))

	tree.GenerateProof(nil)
	tree.Insert(nil)
	tree.Update(nil, nil)
	tree.Lookup(nil)
	VerifyProof(nil, nil)
}
