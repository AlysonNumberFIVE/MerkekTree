package main

import (
	"fmt"
	"testing"
)

func treeDebug_HashPrint(node *Node, prefix string, isLeft bool) {
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
		treeDebug(node.right, newPrefix, false)
	}

	fmt.Printf("%s-", prefix)
	if isLeft {
		fmt.Printf("├── ")
	} else {
		fmt.Printf("└── ")
	}
	fmt.Printf("%s\n", string(node.data))

	if node.left != nil {
		newPrefix := prefix
		if isLeft {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		treeDebug(node.left, newPrefix, true)
	}
}

func treeDebug(node *Node, prefix string, isLeft bool) {
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
		treeDebug(node.right, newPrefix, false)
	}

	fmt.Printf("%s", prefix)
	if isLeft {
		fmt.Printf("├── ")
	} else {
		fmt.Printf("└── ")
	}
	fmt.Printf("%s\n", node.data)
	// fmt.Printf("%s\n", node.data)
	if node.left != nil {
		newPrefix := prefix
		if isLeft {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		treeDebug(node.left, newPrefix, true)
	}
}

func Test_Insert(t *testing.T) {
	testMerkelTree := InitMerkelTree()

	t.Run("Insert 3 Nodes", func(t *testing.T) {
		hash, err := testMerkelTree.Insert([]byte("O"))
		if err != nil {
			t.Errorf("Insert error %+v\n", err)
		}
		if !compareHash(Hash128([]byte("O")), hash) {
			t.Errorf("Error: Insert: mismatch in created hash: Expected: %+v, Actual: %+v\n", Hash128([]byte("O")), hash)
		}

		testNode := &Node{
			data: []byte("O"),
			hash: Hash128([]byte("O")),
		}
		if !compareHash(testMerkelTree.root.data, testNode.data) {
			t.Errorf("error: Insert: Mismatch in root node data created: Expected: %+v, Actual: %+v", testNode.data, testMerkelTree.root.data)
		}
		if !compareHash(testMerkelTree.root.hash, testNode.hash) {
			t.Errorf("error: Insert: Mismatch in root node hash created: Expected: %+v, Actual: %+v", testNode.hash, testMerkelTree.root.hash)
		}

		hash, err = testMerkelTree.Insert([]byte("U"))
		if err != nil {
			t.Errorf("Insert error %+v\n", err)
		}
		if !compareHash(Hash128([]byte("U")), hash) {
			t.Errorf("Error: Insert: mismatch in created hash: Expected: %+v, Actual: %+v\n", Hash128([]byte("U")), hash)
		}

		fmt.Println("test second insert")
		testNode = &Node{
			data: []byte("Y"),
			hash: GenerateHash(Hash128([]byte("O")), Hash128([]byte("U"))),
			left: &Node{
				data: []byte("O"),
				hash: Hash128([]byte("O")),
			},
			right: &Node{
				data: []byte("U"),
				hash: Hash128([]byte("U")),
			},
		}
		if !compareHash(testMerkelTree.root.data, testNode.data) {
			t.Errorf("error: Insert: Mismatch in root node data created: Expected: %+v, Actual: %+v",
				testNode.data, testMerkelTree.root.data)
		}
		if !compareHash(testMerkelTree.root.hash, testNode.hash) {
			t.Errorf("error: Insert: Mismatch in root node hash created: Expected: %+v, Actual: %+v",
				testNode.hash, testMerkelTree.root.hash)
		}
		// left turn
		if !compareHash(testMerkelTree.root.left.data, testNode.left.data) {
			t.Errorf("error: Insert: Mismatch in node.left data created: Expected: %+v, Actual: %+v",
				testNode.left.data, testMerkelTree.root.left.data)
		}
		if !compareHash(testMerkelTree.root.left.data, testNode.left.data) {
			t.Errorf("error: Insert: Mismatch in node.right data created: Expected: %+v, Actual: %+v",
				testNode.left.data, testMerkelTree.root.left.data)
		}

		// right turn
		if !compareHash(testMerkelTree.root.right.data, testNode.right.data) {
			t.Errorf("error: Insert: Mismatch in node.left data created: Expected: %+v, Actual: %+v",
				testNode.right.data, testMerkelTree.root.right.data)
		}
		if !compareHash(testMerkelTree.root.right.data, testNode.right.data) {
			t.Errorf("error: Insert: Mismatch in node.right data created: Expected: %+v, Actual: %+v",
				testNode.right.data, testMerkelTree.root.right.data)
		}

		hash, err = testMerkelTree.Insert([]byte("V"))
		if err != nil {
			t.Errorf("Insert error %+v\n", err)
		}
		if !compareHash(Hash128([]byte("V")), hash) {
			t.Errorf("Error: Insert: mismatch in created hash: Expected: %+v, Actual: %+v\n", Hash128([]byte("V")), hash)
		}

		testNode = &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(Hash128([]byte("V")), Hash128([]byte("O"))),
				Hash128([]byte("U")),
			),
			right: &Node{
				data: []byte("U"),
				hash: Hash128([]byte("U")),
			},
			left: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("V")), Hash128([]byte("V"))),
				right: &Node{
					data: []byte("O"),
					hash: Hash128([]byte("O")),
				},
				left: &Node{
					data: []byte("V"),
					hash: Hash128([]byte("V")),
				},
			},
		}
		// Updated root node
		if !compareHash(testMerkelTree.root.hash, testNode.hash) {
			t.Errorf("error: Insert: Mismatch in root node hash updated: Expected: %+v, Actual: %+v", testNode.hash, testMerkelTree.root.hash)
		}

		// left turn
		if !compareHash(testMerkelTree.root.left.left.data, testNode.left.left.data) {
			t.Errorf("error: Insert: Mismatch in node.left data created: Expected: %+v, Actual: %+v",
				testNode.left.left.data, testMerkelTree.root.left.left.data)
		}
		if !compareHash(testMerkelTree.root.left.left.hash, testNode.left.left.hash) {
			t.Errorf("error: Insert: Mismatch in node.left.left.hash data created: Expected: %+v, Actual: %+v",
				testNode.left.left.hash, testMerkelTree.root.left.left.hash)
		}

		if !compareHash(testMerkelTree.root.left.right.data, testNode.left.right.data) {
			t.Errorf("error: Insert: Mismatch in left.right.data data created: Expected: %+v, Actual: %+v",
				testNode.left.right.data, testMerkelTree.root.left.right.data)
		}

		if !compareHash(testMerkelTree.root.left.right.hash, testNode.left.right.hash) {
			t.Errorf("error: Insert: Mismatch in node.right data created: Expected: %+v, Actual: %+v",
				testNode.left.right.data, testMerkelTree.root.left.right.data)
		}

		// right turn
		if !compareHash(testMerkelTree.root.right.data, testNode.right.data) {
			t.Errorf("error: Insert: Mismatch in node.left data created: Expected: %+v, Actual: %+v",
				testNode.right.data, testMerkelTree.root.right.data)
		}
		if !compareHash(testMerkelTree.root.right.hash, testNode.right.hash) {
			t.Errorf("error: Insert: Mismatch in node.right data created: Expected: %+v, Actual: %+v",
				testNode.right.hash, testMerkelTree.root.right.hash)
		}
	})

}

func Test_Update(t *testing.T) {

	t.Run("Test Update", func(t *testing.T) {
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		testMerkelTree.Insert([]byte("B"))
		testMerkelTree.Insert([]byte("C"))
		testMerkelTree.Insert([]byte("D"))
		//      changing this  ------------>  to this
		//            O                          O
		//          /   \                      /   \
		//         O     O                    O     O
		//       /  \   /  \                /  \   /  \
		//      C    A D    B              C    E D    B
		ptrA := testMerkelTree.root.left.right
		if !compareHash(ptrA.hash, Hash128([]byte("A"))) {
			t.Error("Error: Update: Insert node error")
		}
		hashE, err := testMerkelTree.Update([]byte("E"), Hash128([]byte("A")))
		if err != nil {
			t.Errorf("Error: Update: %+v\n", err)
		}
		if !compareHash(ptrA.hash, hashE) {
			t.Errorf("Error: Update: hash unmodified; expected %+v, actual: %+v\n", Hash128([]byte("E")), ptrA.hash)
		}

	})

	t.Run("Test parental hash updates after modification", func(t *testing.T) {
		//
		//      changing this  ------------>  to this
		//            O                          O
		//          /   \                      /   \
		//         O     O                    O     O
		//       /  \   /  \                /  \   /  \
		//      C    A D    B              C    E D    B
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		testMerkelTree.Insert([]byte("B"))
		testMerkelTree.Insert([]byte("C"))
		testMerkelTree.Insert([]byte("D"))

		_, err := testMerkelTree.Update([]byte("E"), Hash128([]byte("A")))
		if err != nil {
			t.Errorf("Error: Update: %+v\n", err)
		}
		//   checking this path of hashs
		//            O
		//          /
		//         O     O
		//          \
		//      C    E D    B
		testNode := &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(Hash128([]byte("C")), Hash128([]byte("E"))),
				GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
			),
			left: &Node{
				hash: GenerateHash(Hash128([]byte("C")), Hash128([]byte("E"))),
				left: &Node{
					hash: Hash128([]byte("C")),
				},
				right: &Node{
					hash: Hash128([]byte("E")),
				},
			},
			right: &Node{
				hash: GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
				left: &Node{
					hash: Hash128([]byte("D")),
				},
				right: &Node{
					hash: Hash128([]byte("B")),
				},
			},
		}
		// set up backtracks
		testNode.left.right.prev = testNode.left
		testNode.left.prev = testNode

		root := testMerkelTree.root.left.right
		testTraverse := testNode.left.right
		steps := 0
		for testTraverse != nil {
			if !compareHash(root.hash, testTraverse.hash) {
				t.Errorf("Error: Update: Traversing up the tree: hash mismatch at step %d. Expected: %+v, Actual: %+v\n",
					steps, string(testTraverse.hash), string(root.hash))
			}
			testTraverse = testTraverse.prev
			root = root.prev
			steps += 1
		}
	})

	t.Run("Update the same node twice", func(t *testing.T) {
		//
		//      changing this  ------------>  to this
		//            O                          O
		//          /   \                      /   \
		//         O     O                    O     O
		//       /  \   /  \                /  \   /  \
		//      C    A D    B              C    E D    B
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		testMerkelTree.Insert([]byte("B"))
		testMerkelTree.Insert([]byte("C"))
		testMerkelTree.Insert([]byte("D"))

		hashE, err := testMerkelTree.Update([]byte("E"), Hash128([]byte("A")))
		if err != nil {
			t.Errorf("Error: Update (same node twice): %+v\n", err)
		}
		hashF, err := testMerkelTree.Update([]byte("F"), hashE)
		if err != nil {
			t.Errorf("Error: Update (same node twice): %+v\n", err)
		}
		//   checking this path of hashs
		//            O
		//          /
		//         O     O
		//          \
		//      C    F D    B
		testNode := &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(Hash128([]byte("C")), Hash128([]byte("F"))),
				GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
			),
			left: &Node{
				hash: GenerateHash(Hash128([]byte("C")), hashF),
				left: &Node{
					hash: Hash128([]byte("C")),
				},
				right: &Node{
					hash: Hash128([]byte("F")),
				},
			},
			right: &Node{
				hash: GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
				left: &Node{
					hash: Hash128([]byte("D")),
				},
				right: &Node{
					hash: Hash128([]byte("B")),
				},
			},
		}
		// set up backtracks
		testNode.left.right.prev = testNode.left
		testNode.left.prev = testNode

		root := testMerkelTree.root.left.right
		testTraverse := testNode.left.right
		steps := 0
		for testTraverse != nil {
			if !compareHash(root.hash, testTraverse.hash) {
				t.Errorf("Error: Update (same node twice): Traversing up the tree: hash mismatch at step %d. Expected: %+v, Actual: %+v\n",
					steps, string(testTraverse.hash), string(root.hash))
			}
			testTraverse = testTraverse.prev
			root = root.prev
			steps += 1
		}
	})

	t.Run("Update fail; node doesn't exist", func(t *testing.T) {
		//   The tree
		//            O
		//          /   \
		//         O     O
		//       /  \  /   \
		//      C    A D    B
		//
		//   Node P doesn't exist
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		testMerkelTree.Insert([]byte("B"))
		testMerkelTree.Insert([]byte("C"))
		testMerkelTree.Insert([]byte("D"))

		// P doesn't exist
		hashE, err := testMerkelTree.Update([]byte("P"), Hash128([]byte("P")))
		if err == nil {
			t.Errorf("Error: Update: Nonexitent error not triggering")
		}
		if hashE != nil {
			t.Errorf("Error: Update: Nonexistent hash error")
		}
	})

	t.Run("Update the root on a tree that's 1 node long", func(t *testing.T) {
		//   The tree
		//            O
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		hashB, err := testMerkelTree.Update([]byte("B"), Hash128([]byte("A")))
		if err != nil {
			t.Errorf("Error: Update: Nonexitent error not triggering")
		}
		testNode := &Node{
			data: []byte("B"),
			hash: Hash128([]byte("B")),
		}
		if !compareHash(testNode.hash, hashB) {
			t.Error("Error: Update: Root node modification failed")
		}
		if !compareHash(testNode.hash, hashB) {
			t.Error("Error: Update: Root node modification failed")
		}
	})

	t.Run("Update the root on a tree that's 2 node long", func(t *testing.T) {
		//   The tree
		//            O
		//          /   \
		//         A     B
		//
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		testMerkelTree.Insert([]byte("B"))
		hashC, err := testMerkelTree.Update([]byte("C"), Hash128([]byte("A")))
		if err != nil {
			t.Errorf("Error: Update: Nonexitent error not triggering")
		}

		testNode := &Node{
			data: []byte("Y"),
			hash: GenerateHash(Hash128([]byte("C")), Hash128([]byte("B"))),
			left: &Node{
				data: []byte("C"),
				hash: Hash128([]byte("C")),
			},
			right: &Node{
				data: []byte("B"),
				hash: Hash128([]byte("B")),
			},
		}

		if !compareHash(testNode.left.hash, hashC) {
			t.Error("Error: Update: 2nd node modification failed")
		}
		if !compareHash(testNode.hash, testMerkelTree.root.hash) {
			t.Error("Error: Update: 2nd node: Root node modification failed")
		}
	})

}

func Test_InsertNode(t *testing.T) {
	t.Run("Insert a single Node", func(t *testing.T) {
		//   The tree
		//            O
		//          /   \
		//         A     B
		//
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		testMerkelTree.Insert([]byte("B"))

		testNode := &Node{
			data: []byte("Y"),
			hash: GenerateHash(Hash128([]byte("B")), Hash128([]byte("A"))),
			left: &Node{
				data: []byte("B"),
				hash: Hash128([]byte("B")),
			},
			right: &Node{
				data: []byte("A"),
				hash: Hash128([]byte("A")),
			},
		}
		prevNodePtr := testNode
		currentPtr := testNode.left
		newLeaf := InsertNode(currentPtr, &prevNodePtr, []byte("C"), Hash128([]byte("C")))
		expectedNewLeaf := &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(Hash128([]byte("C")), Hash128([]byte("B"))),
				Hash128([]byte("A")),
			),
			right: &Node{
				data: []byte("A"),
				hash: Hash128([]byte("A")),
			},
			left: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("C")), Hash128([]byte("B"))),
				right: &Node{
					data: []byte("B"),
					hash: Hash128([]byte("B")),
				},
				left: &Node{
					data: []byte("C"),
					hash: Hash128([]byte("C")),
				},
			},
		}
		expectedNewLeaf.left.left.prev = expectedNewLeaf.left
		expectedNewLeaf.left.right.prev = expectedNewLeaf.left
		if !compareHash(newLeaf.prev.hash, expectedNewLeaf.left.left.prev.hash) {
			t.Errorf("Error: InsertNode: parent hash mismatch. Expected: %+v, Actual: %+v\n", expectedNewLeaf.left.left.prev.hash, newLeaf.prev.hash)
		}
		if !compareHash(newLeaf.prev.data, expectedNewLeaf.left.left.prev.data) {
			t.Errorf("Error: InsertNode: parent data mismatch. Expected: %+v, Actual: %+v\n", expectedNewLeaf.left.left.prev.data, newLeaf.prev.data)
		}
	})
}

func Test_navigateTree(t *testing.T) {
	t.Run("Get all the deepest nodes in the tree", func(t *testing.T) {
		// Tree structure:
		// 				  O
		//              /  \
		//			  /      \
		//           O        O
		//         /  \      /  \
		//        E    F    /    \
		//                 O      O
		//                / \    /  \
		//               C  D   A    V

		testMerkelTree := InitMerkelTree()
		testNode := &Node{
			right: &Node{
				right: &Node{
					right: &Node{
						data: []byte("B"),
					},
					left: &Node{
						data: []byte("A"),
					},
				},
				left: &Node{
					right: &Node{
						data: []byte("D"),
					},
					left: &Node{
						data: []byte("C"),
					},
				},
			},
			left: &Node{
				right: &Node{
					data: []byte("E"),
				},
				left: &Node{
					data: []byte("F"),
				},
			},
		}
		testMerkelTree.root = testNode
		expectedNodeList := []NodeDepth{
			{node: &Node{
				data: []byte("F"),
			}, depth: 2},
			{node: &Node{
				data: []byte("E"),
			}, depth: 2},
			{node: &Node{
				data: []byte("C"),
			}, depth: 3},
			{node: &Node{
				data: []byte("D"),
			}, depth: 3},
			{node: &Node{
				data: []byte("A"),
			}, depth: 3},
			{node: &Node{
				data: []byte("B"),
			}, depth: 3},
		}
		actualNodeDepth := testMerkelTree.navigateTree()

		if len(actualNodeDepth) != len(expectedNodeList) {
			t.Errorf("Error: navigateTree: Number of depths mismatch: Expected: %+v, Actual: %+v\n", len(expectedNodeList), len(actualNodeDepth))
		}
		for index, depthNode := range actualNodeDepth {
			if !compareHash(expectedNodeList[index].node.data, depthNode.node.data) {
				t.Errorf(
					"expected: %s\n depthNode: %s\n", string(expectedNodeList[index].node.data), string(depthNode.node.data),
				)
			}
			if expectedNodeList[index].depth != depthNode.depth {
				t.Errorf(
					"depth mismatch: Expect: %d, Actual: %d\n", expectedNodeList[index].depth, depthNode.depth,
				)
			}
		}
	})

}

func Test_Lookup(t *testing.T) {
	t.Run("Test Lookup of existing data", func(t *testing.T) {
		//   The tree
		//            O
		//          /   \
		//         O     O
		//       /  \  /   \
		//      C    A D    B
		//
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		testMerkelTree.Insert([]byte("B"))
		testMerkelTree.Insert([]byte("C"))
		testMerkelTree.Insert([]byte("D"))

		node, err := testMerkelTree.Lookup(Hash128([]byte("C")))
		if err != nil {
			t.Errorf("Error: Lookup: existing data error: %+v\n", err)
		}

		if !compareHash(node.hash, Hash128([]byte("C"))) {
			t.Error("Error: Lookup: incorrect hash value returned")
		}

		hashE, err := testMerkelTree.Update([]byte("E"), Hash128([]byte("C")))

		node, err = testMerkelTree.Lookup(hashE)
		if err != nil {
			t.Errorf("Error: Lookup: existing data error: %+v\n", err)
		}

		if !compareHash(node.hash, hashE) {
			t.Error("Error: Lookup: incorrect hash value returned")
		}

		hashH, err := testMerkelTree.Update([]byte("H"), hashE)
		node, err = testMerkelTree.Lookup(hashH)
		if err != nil {
			t.Errorf("Error: Lookup: existing data error: %+v\n", err)
		}

		if !compareHash(node.hash, hashH) {
			t.Error("Error: Lookup: incorrect hash value returned")
		}
	})

	t.Run("Query the root itself", func(t *testing.T) {
		// The tree
		//    	 O
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		node, err := testMerkelTree.Lookup(Hash128([]byte("A")))
		if err != nil {
			t.Errorf("Error: Lookup: existing data error: %+v\n", err)
		}

		if !compareHash(node.hash, Hash128([]byte("A"))) {
			t.Error("Error: Lookup: incorrect hash value returned")
		}
		if !compareHash(node.data, []byte("A")) {
			t.Errorf("Error: Lookup: incorrect data returned")
		}
	})

	t.Run("query a node with a stale hash", func(t *testing.T) {
		//   The tree
		//            O
		//          /   \
		//         O     O
		//       /  \  /   \
		//      C    A D    B
		//
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		testMerkelTree.Insert([]byte("B"))
		testMerkelTree.Insert([]byte("C"))
		testMerkelTree.Insert([]byte("D"))

		hashE, err := testMerkelTree.Update([]byte("E"), Hash128([]byte("C")))
		hashF, err := testMerkelTree.Update([]byte("F"), hashE)
		hashG, err := testMerkelTree.Update([]byte("G"), hashF)
		node, err := testMerkelTree.Lookup(Hash128([]byte("C")))
		if err != nil {
			t.Errorf("Error: Lookup: existing data error: %+v\n", err)
		}
		if !compareHash(node.hash, hashG) {
			t.Error("Error: Lookup: incorrect hash value returned")
		}

		node, err = testMerkelTree.Lookup(hashF)
		if err != nil {
			t.Errorf("Error: Lookup: existing data error: %+v\n", err)
		}
		if !compareHash(node.hash, hashG) {
			t.Error("Error: Lookup: incorrect hash value returned")
		}
	})

}

func Test_GenerateProof(t *testing.T) {
	tests := []struct {
		name          string
		expectedProof MerkelProof
		testData      [][]byte
	}{
		{
			name: "Generate a proof",
			//   The tree
			//            O
			//          /   \
			//         O     O
			//       /  \  /   \
			//      C    A D    B
			//
			testData: [][]byte{
				[]byte("A"),
				[]byte("B"),
				[]byte("C"),
				[]byte("D"),
			},
			expectedProof: MerkelProof{
				LeafHash: Hash128([]byte("A")),
				ProofList: [][]byte{
					Hash128([]byte("A")),
					Hash128([]byte("C")),
					GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
				},
				Directions: []bool{
					true, false, true,
				},
			},
		},
		{
			name: "Generate proof for a single node",
			// The tree
			//      O
			testData: [][]byte{
				[]byte("A"),
			},
			expectedProof: MerkelProof{
				LeafHash: Hash128([]byte("A")),
				ProofList: [][]byte{
					Hash128([]byte("A")),
				},
				Directions: []bool{
					true,
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testMerkelTree := InitMerkelTree()
			for _, data := range testCase.testData {
				testMerkelTree.Insert(data)
			}

			actualProof, err := testMerkelTree.GenerateProof(Hash128([]byte("A")))
			if err != nil {
				t.Errorf("Error: GenerateProof: %+v\n", err)
			}

			if !compareHash(testCase.expectedProof.LeafHash, actualProof.LeafHash) {
				t.Errorf("Error: GenerateProof: leaf hash mismatch: Expected: %+v, Actual: %+v\n",
					testCase.expectedProof.LeafHash, actualProof.LeafHash)
			}
			if len(testCase.expectedProof.ProofList) != len(actualProof.ProofList) {
				t.Errorf("Error: GenerateProof: ProofList length mismatch: Expected: %+v, Actual %+v\n",
					len(testCase.expectedProof.ProofList), len(actualProof.ProofList))
			}
			if len(testCase.expectedProof.Directions) != len(actualProof.Directions) {
				t.Errorf("Error: GenerateProof: ProofList length mismatch: Expected: %+v, Actual %+v\n",
					len(testCase.expectedProof.Directions), len(actualProof.Directions))
			}

			for index, direction := range actualProof.Directions {
				if testCase.expectedProof.Directions[index] != direction {
					t.Errorf("Error: GenerateProof: Direction mismatch: step: %d, Expected: %+v, Actual %+v\n",
						index, testCase.expectedProof.Directions[index], direction)
				}
			}
			for index, hash := range actualProof.ProofList {
				if !compareHash(testCase.expectedProof.ProofList[index], hash) {
					t.Errorf("Error: GenerateProof: ProofList mismatch: step: %d, Expected: %+v, Actual %+v\n",
						index, testCase.expectedProof.ProofList[index], hash)
				}
			}
		})
	}
}

func Test_VerifyProof(t *testing.T) {
	tests := []struct {
		name          string
		testData      [][]byte
		testProof     *MerkelProof
		expectedProof []byte
		expected      bool
	}{
		{
			name: "Verify a proof",
			//   The tree
			//            O
			//          /   \
			//         O     O
			//       /  \  /   \
			//      C    A D    B
			//
			testData: [][]byte{
				[]byte("A"),
				[]byte("B"),
				[]byte("C"),
				[]byte("D"),
			},
			testProof: &MerkelProof{
				LeafHash: Hash128([]byte("A")),
				ProofList: [][]byte{
					Hash128([]byte("A")),
					Hash128([]byte("C")),
					GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
				},
				Directions: []bool{
					true, false, true,
				},
			},
			expectedProof: GenerateHash(
				GenerateHash(Hash128([]byte("C")), Hash128([]byte("C"))),
				GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
			),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testMerkelTree := InitMerkelTree()
			for _, data := range testCase.testData {
				testMerkelTree.Insert(data)
			}

			verify := VerifyProof(testCase.testProof,
				testMerkelTree.root.hash)

			if !verify {
				t.Error("Error: VerifyProof: verify error: verification failed")
			}
		})
	}
}

func Test_MerkelTree(t *testing.T) {
	t.Run("Run the merkel tree; all functions in one session", func(t *testing.T) {
		testMerkelTree := InitMerkelTree()
		testMerkelTree.Insert([]byte("A"))
		testMerkelTree.Insert([]byte("B"))
		testMerkelTree.Insert([]byte("C"))
		testMerkelTree.Insert([]byte("D"))
		//   The tree
		//            O
		//          /   \
		//         O     O
		//       /  \  /   \
		//      C    A D    B
		//
		testInsert := &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(Hash128([]byte("C")), Hash128([]byte("A"))),
				GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
			),
			left: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("C")), Hash128([]byte("A"))),
				left: &Node{
					data: []byte("C"),
					hash: Hash128([]byte("C")),
				},
				right: &Node{
					data: []byte("A"),
					hash: Hash128([]byte("A")),
				},
			},
			right: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
				left: &Node{
					data: []byte("D"),
					hash: Hash128([]byte("D")),
				},
				right: &Node{
					data: []byte("B"),
					hash: Hash128([]byte("B")),
				},
			},
		}
		testInsert.left.left.prev = testInsert.left
		testInsert.left.right.prev = testInsert.left
		testInsert.right.left.prev = testInsert.right
		testInsert.right.right.prev = testInsert.right
		if !compareHash(testInsert.left.data, testMerkelTree.root.left.data) {
			t.Errorf("Error: Insert: data mismatch in left.data, Expected: %+v, Actual: %+v\n", string(testInsert.left.data), string(testMerkelTree.root.left.data))
		}
		if !compareHash(testInsert.left.hash, testMerkelTree.root.left.hash) {
			t.Errorf("Error: Insert: hash mismatch in left.hash, Expected: %+v, Actual: %+v\n", string(testInsert.left.hash), string(testMerkelTree.root.left.hash))
		}
		if !compareHash(testInsert.right.data, testMerkelTree.root.right.data) {
			t.Errorf("Error: Insert: data mismatch in right.data, Expected: %+v, Actual: %+v\n", string(testInsert.right.data), string(testMerkelTree.root.right.data))
		}
		if !compareHash(testInsert.right.hash, testMerkelTree.root.right.hash) {
			t.Errorf("Error: Insert: hash mismatch in right.hash, Expected: %+v, Actual: %+v\n", string(testInsert.right.hash), string(testMerkelTree.root.right.hash))
		}

		if !compareHash(testInsert.left.left.data, testMerkelTree.root.left.left.data) {
			t.Errorf("Error: Insert: data mismatch in left.left.data, Expected: %+v, Actual: %+v\n", string(testInsert.left.left.data), string(testMerkelTree.root.left.left.data))
		}
		if !compareHash(testInsert.left.left.hash, testMerkelTree.root.left.left.hash) {
			t.Errorf("Error: Insert: hash mismatch in left.left.hash, Expected: %+v, Actual: %+v\n", string(testInsert.left.left.hash), string(testMerkelTree.root.left.left.hash))
		}
		if !compareHash(testInsert.left.right.data, testMerkelTree.root.left.right.data) {
			t.Errorf("Error: Insert: data mismatch in left.right.data, Expected: %+v, Actual: %+v\n", string(testInsert.left.right.data), string(testMerkelTree.root.left.right.data))
		}
		if !compareHash(testInsert.left.right.hash, testMerkelTree.root.left.right.hash) {
			t.Errorf("Error: Insert: hash mismatch in left.right.hash, Expected: %+v, Actual: %+v\n", string(testInsert.left.right.hash), string(testMerkelTree.root.left.right.hash))
		}

		if !compareHash(testInsert.right.left.data, testMerkelTree.root.right.left.data) {
			t.Errorf("Error: Insert: data mismatch in left.left.data, Expected: %+v, Actual: %+v\n", string(testInsert.right.left.data), string(testMerkelTree.root.right.left.data))
		}
		if !compareHash(testInsert.right.left.hash, testMerkelTree.root.right.left.hash) {
			t.Errorf("Error: Insert: hash mismatch in left.left.hash, Expected: %+v, Actual: %+v\n", string(testInsert.right.left.hash), string(testMerkelTree.root.right.left.hash))
		}
		if !compareHash(testInsert.right.right.data, testMerkelTree.root.right.right.data) {
			t.Errorf("Error: Insert: data mismatch in left.right.data, Expected: %+v, Actual: %+v\n", string(testInsert.right.right.data), string(testMerkelTree.root.right.right.data))
		}
		if !compareHash(testInsert.right.right.hash, testMerkelTree.root.right.right.hash) {
			t.Errorf("Error: Insert: hash mismatch in left.right.hash, Expected: %+v, Actual: %+v\n", string(testInsert.right.right.hash), string(testMerkelTree.root.right.right.hash))
		}

		testInsert = &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(Hash128([]byte("C")), Hash128([]byte("E"))),
				GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
			),
			left: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("C")), Hash128([]byte("E"))),
				left: &Node{
					data: []byte("C"),
					hash: Hash128([]byte("C")),
				},
				right: &Node{
					data: []byte("E"),
					hash: Hash128([]byte("E")),
				},
			},
			right: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
				left: &Node{
					data: []byte("D"),
					hash: Hash128([]byte("D")),
				},
				right: &Node{
					data: []byte("B"),
					hash: Hash128([]byte("B")),
				},
			},
		}
		hashE, err := testMerkelTree.Update([]byte("E"), Hash128([]byte("A")))
		if err != nil {
			t.Errorf("Error: All: Update erro %+v\n", err)
		}

		if !compareHash(testInsert.hash, testMerkelTree.root.hash) {
			t.Errorf("Error: All: Update: root hash mismatch: Expected %+v, Actual: %+v\n",
				string(testInsert.hash), string(testMerkelTree.root.hash))
		}
		testInsert = &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(Hash128([]byte("C")), Hash128([]byte("F"))),
				GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
			),
			left: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("C")), Hash128([]byte("F"))),
				left: &Node{
					data: []byte("C"),
					hash: Hash128([]byte("C")),
				},
				right: &Node{
					data: []byte("F"),
					hash: Hash128([]byte("F")),
				},
			},
			right: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
				left: &Node{
					data: []byte("D"),
					hash: Hash128([]byte("D")),
				},
				right: &Node{
					data: []byte("B"),
					hash: Hash128([]byte("B")),
				},
			},
		}
		hashF, err := testMerkelTree.Update([]byte("F"), Hash128([]byte("A")))
		if err != nil {
			t.Errorf("Error: All: Update erro %+v\n", err)
		}

		if !compareHash(testInsert.hash, testMerkelTree.root.hash) {
			t.Errorf("Error: All: Update: root hash mismatch: Expected %+v, Actual: %+v\n",
				string(testInsert.hash), string(testMerkelTree.root.hash))
		}

		testInsert = &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(Hash128([]byte("C")), Hash128([]byte("Q"))),
				GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
			),
			left: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("C")), Hash128([]byte("Q"))),
				left: &Node{
					data: []byte("C"),
					hash: Hash128([]byte("C")),
				},
				right: &Node{
					data: []byte("Q"),
					hash: Hash128([]byte("Q")),
				},
			},
			right: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("D")), Hash128([]byte("B"))),
				left: &Node{
					data: []byte("D"),
					hash: Hash128([]byte("D")),
				},
				right: &Node{
					data: []byte("B"),
					hash: Hash128([]byte("B")),
				},
			},
		}
		hashQ, err := testMerkelTree.Update([]byte("Q"), hashE)
		if err != nil {
			t.Errorf("Error: All: Update erro %+v\n", err)
		}
		if !compareHash(testInsert.hash, testMerkelTree.root.hash) {
			t.Errorf("Error: All: Update: root hash mismatch: Expected %+v, Actual: %+v\n",
				string(testInsert.hash), string(testMerkelTree.root.hash))
		}

		testInsert = &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(Hash128([]byte("C")), Hash128([]byte("Q"))),
				GenerateHash(Hash128([]byte("W")), Hash128([]byte("B"))),
			),
			left: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("C")), Hash128([]byte("Q"))),
				left: &Node{
					data: []byte("C"),
					hash: Hash128([]byte("C")),
				},
				right: &Node{
					data: []byte("Q"),
					hash: Hash128([]byte("Q")),
				},
			},
			right: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("W")), Hash128([]byte("B"))),
				left: &Node{
					data: []byte("W"),
					hash: Hash128([]byte("W")),
				},
				right: &Node{
					data: []byte("B"),
					hash: Hash128([]byte("B")),
				},
			},
		}
		_, err = testMerkelTree.Update([]byte("W"), Hash128([]byte("D")))
		if err != nil {
			t.Errorf("Error: All: Update erro %+v\n", err)
		}
		if !compareHash(testInsert.hash, testMerkelTree.root.hash) {
			t.Errorf("Error: All: Update: root hash mismatch: Expected %+v, Actual: %+v\n",
				string(testInsert.hash), string(testMerkelTree.root.hash))
		}

		node, err := testMerkelTree.Lookup(Hash128([]byte("B")))
		if !compareHash(Hash128([]byte("B")),
			node.hash) {
			t.Errorf("Error: All: Lookup: hash chain error: Expected: %+v, Actual: %+v\n", string(Hash128([]byte("B"))), string(node.hash))
		}

		// Original hash
		node, err = testMerkelTree.Lookup(Hash128([]byte("A")))
		if !compareHash(Hash128([]byte("Q")), node.hash) {
			t.Errorf("Error: All: Lookup: hash chain error: Expected: %+v, Actual: %+v\n", string(Hash128([]byte("Q"))), string(node.hash))
		}

		node, err = testMerkelTree.Lookup(hashE)
		if !compareHash(Hash128([]byte("Q")), node.hash) {
			t.Errorf("Error: All: Lookup: hash chain error: Expected: %+v, Actual: %+v\n", string(Hash128([]byte("Q"))), string(node.hash))
		}

		node, err = testMerkelTree.Lookup(hashF)
		if !compareHash(Hash128([]byte("Q")), node.hash) {
			t.Errorf("Error: All: Lookup: hash chain error: Expected: %+v, Actual: %+v\n", string(Hash128([]byte("Q"))), string(node.hash))
		}

		node, err = testMerkelTree.Lookup(hashQ)
		if !compareHash(Hash128([]byte("Q")), node.hash) {
			t.Errorf("Error: All: Lookup: hash chain error: Expected: %+v, Actual: %+v\n", string(Hash128([]byte("Q"))), string(node.hash))
		}

		hash4444, _ := testMerkelTree.Insert([]byte("4444"))

		testInsert = &Node{
			data: []byte("Y"),
			hash: GenerateHash(
				GenerateHash(
					GenerateHash(
						Hash128([]byte("4444")), Hash128([]byte("C")),
					),
					Hash128([]byte("Q")),
				),
				GenerateHash(
					Hash128([]byte("W")), Hash128([]byte("B")),
				),
			),
			right: &Node{
				data: []byte("X"),
				hash: GenerateHash(Hash128([]byte("W")), Hash128([]byte("W"))),
				right: &Node{
					data: []byte("B"),
					hash: Hash128([]byte("W")),
				},
				left: &Node{
					data: []byte("W"),
					hash: Hash128([]byte("W")),
				},
			},
			left: &Node{
				data: []byte("X"),
				hash: GenerateHash(
					GenerateHash(Hash128([]byte("4444")), Hash128([]byte("C"))),
					Hash128([]byte("Q")),
				),
				right: &Node{
					data: []byte("Q"),
					hash: Hash128([]byte("Q")),
				},
				left: &Node{
					data: []byte("X"),
					hash: GenerateHash(Hash128([]byte("4444")), Hash128([]byte("C"))),
					right: &Node{
						data: []byte("C"),
						hash: Hash128([]byte("C")),
					},
					left: &Node{
						data: []byte("4444"),
						hash: Hash128([]byte("4444")),
					},
				},
			},
		}

		if !compareHash(testMerkelTree.root.hash, testInsert.hash) {
			t.Errorf("Error: All: Insert (2): Hash mismatch: Expected: %+v, Actual: %+v\n", string(testInsert.hash), string(testMerkelTree.root.hash))
		}

		hash, err := testMerkelTree.Update([]byte("1111"), hash4444)
		if !compareHash(hash, testMerkelTree.root.left.left.left.hash) {
			t.Errorf("Error: All: Update (2): Hash mismatch: Expected: %+v, Actual: %+v\n", string(hash), string(testMerkelTree.root.left.left.left.hash))
		}

		nodeNewHash, err := testMerkelTree.Lookup(hash4444)
		nodeOldHash, err := testMerkelTree.Lookup(Hash128([]byte("4444")))
		if nodeNewHash != nodeOldHash {
			t.Errorf("Error: All: Update (2): Data mismatch: old hash: %+v, new hash: %+v\n", string(nodeOldHash.data), string(nodeNewHash.data))
		}

	})
}
