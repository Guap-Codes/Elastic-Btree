package tree

import (
	"elastic-btree/pkg/logger" // Import the custom logger
	"sync"
)

// Tree represents the Elastic B-Tree.
type Tree struct {
	Root       *Node              `json:"root"`   // Root node of the tree
	Degree     int                `json:"degree"` // Minimum degree of the tree
	Size       int                `json:"size"`   // Total number of keys in the tree
	Height     int                `json:"height"` // Height of the tree
	Lock       sync.RWMutex       `json:"-"`      // Mutex for concurrent access
	Logger     *logger.Logger     `json:"-"`      // Custom logger for debugging
	Comparator func(a, b int) int `json:"-"`      // Custom key comparator (default: ascending order)
}

// NewTree creates a new Elastic B-Tree with the given degree and logger.
func NewTree(degree int, logger *logger.Logger) *Tree {
	if degree < 2 {
		logger.Panicf("degree must be at least 2")
	}
	return &Tree{
		Degree:     degree,
		Root:       nil,
		Size:       0,
		Height:     0,
		Comparator: func(a, b int) int { return a - b }, // Default comparator
		Logger:     logger,
	}
}

// Insert inserts a key into the tree.
func (t *Tree) Insert(key int, value interface{}) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	if t.Root == nil {
		t.Root = &Node{
			Keys:     []int{key},
			Values:   []interface{}{value},
			Children: []*Node{},
			IsLeaf:   true,
			Size:     1,
			MaxKeys:  2*t.Degree - 1,
			MinKeys:  t.Degree - 1,
		}
		t.Size++
		t.Height = 1
		t.Logger.Infof("Insert: created new root with key: %d", key)
		return
	}

	// Check invariants before insertion.
	t.checkInvariants(t.Root)
	t.Logger.Infof("Insert: inserting key %d", key)

	if t.Root.Size == t.Root.MaxKeys {
		// Split the root if it's full.
		newRoot := &Node{
			Keys:     []int{},
			Values:   []interface{}{},
			Children: []*Node{t.Root},
			IsLeaf:   false,
			Size:     0,
			MaxKeys:  2*t.Degree - 1,
			MinKeys:  t.Degree - 1,
		}
		t.Root.Parent = newRoot
		t.splitChild(newRoot, 0)
		t.Root = newRoot
		t.Height++
		t.Logger.Infof("Insert: root split; new root keys: %v", newRoot.Keys)
	}

	t.insertNonFull(t.Root, key, value)
	t.Size++

	// Check invariants after insertion.
	t.checkInvariants(t.Root)
	t.Logger.Infof("Insert: finished inserting key %d", key)
}

// insertNonFull inserts a key into a non-full node.
func (t *Tree) insertNonFull(node *Node, key int, value interface{}) {
	i := node.Size - 1
	if node.IsLeaf {
		// Insert into a leaf node
		for i >= 0 && t.Comparator(node.Keys[i], key) > 0 {
			i--
		}
		node.Keys = append(node.Keys[:i+1], append([]int{key}, node.Keys[i+1:]...)...)
		node.Values = append(node.Values[:i+1], append([]interface{}{value}, node.Values[i+1:]...)...)
		node.Size++
	} else {
		// Insert into an internal node
		for i >= 0 && t.Comparator(node.Keys[i], key) > 0 {
			i--
		}
		i++
		if node.Children[i].Size == node.Children[i].MaxKeys {
			// Split the child if it's full
			t.splitChild(node, i)
			if t.Comparator(node.Keys[i], key) < 0 {
				i++
			}
		}
		t.insertNonFull(node.Children[i], key, value)
	}
}

// splitChild splits a full child of a node.
func (t *Tree) splitChild(parent *Node, index int) {
    child := parent.Children[index]
    t.Logger.Infof("splitChild: splitting child at index %d with keys: %v", index, child.Keys)
    // Check invariant before split.
    t.checkInvariants(child)
    t.normalizeChildren(parent)

    medianKey := child.Keys[t.Degree-1]
    medianValue := child.Values[t.Degree-1]

    newChild := &Node{
        Keys:     make([]int, t.Degree-1),
        Values:   make([]interface{}, t.Degree-1),
        Children: []*Node{},
        IsLeaf:   child.IsLeaf,
        Size:     t.Degree - 1,
        MaxKeys:  2*t.Degree - 1,
        MinKeys:  t.Degree - 1,
        Parent:   parent,
    }

    // Copy second half of keys/values to newChild.
    copy(newChild.Keys, child.Keys[t.Degree:])
    copy(newChild.Values, child.Values[t.Degree:])
    if !child.IsLeaf {
        // Copy the second half of children to the new child
        newChild.Children = append(newChild.Children, child.Children[t.Degree:]...)
        for _, c := range newChild.Children {
            c.Parent = newChild
        }
    }

    // Update the original child
    child.Keys = child.Keys[:t.Degree-1]
    child.Values = child.Values[:t.Degree-1]
    child.Size = t.Degree - 1
    if !child.IsLeaf {
        // Retain t.Degree children for non-leaf nodes
        child.Children = child.Children[:t.Degree]
    }

    // Insert the median key/value into the parent and attach the new child.
    parent.Keys = append(parent.Keys[:index], append([]int{medianKey}, parent.Keys[index:]...)...)
    parent.Values = append(parent.Values[:index], append([]interface{}{medianValue}, parent.Values[index:]...)...)
    parent.Children = append(parent.Children[:index+1], append([]*Node{newChild}, parent.Children[index+1:]...)...)
    parent.Size = len(parent.Keys)

    t.normalizeChildren(parent)
    t.Logger.Infof("splitChild: after split, parent keys: %v, children count: %d", parent.Keys, len(parent.Children))
    // Check invariants after split.
    t.checkInvariants(child)
    t.checkInvariants(newChild)
    t.checkInvariants(parent)
}

// Search searches for a key in the tree and returns its value (if found).
func (t *Tree) Search(key int) (interface{}, bool) {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	return t.searchNode(t.Root, key)
}

// searchNode searches for a key in a subtree rooted at the given node.
func (t *Tree) searchNode(node *Node, key int) (interface{}, bool) {
	if node == nil {
		return nil, false
	}

	i := 0
	for i < node.Size && t.Comparator(node.Keys[i], key) < 0 {
		i++
	}

	if i < node.Size && t.Comparator(node.Keys[i], key) == 0 {
		// Key found
		return node.Values[i], true
	}

	if node.IsLeaf {
		// Key not found
		return nil, false
	}

	// Search in the appropriate child
	return t.searchNode(node.Children[i], key)
}

// Delete deletes a key from the tree.
func (t *Tree) Delete(key int) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	if t.Root == nil {
		return
	}
	t.Logger.Infof("Delete: deleting key %d", key)
	t.deleteNode(t.Root, key)
	if t.Root.Size == 0 && !t.Root.IsLeaf {
		t.Root = t.Root.Children[0]
		t.Height--
		t.Logger.Infof("Delete: root became empty, new root keys: %v", t.Root.Keys)
	}
	t.Size--
	t.checkInvariants(t.Root)
	t.Logger.Infof("Delete: finished deleting key %d", key)
}

// deleteNode deletes a key from a subtree rooted at the given node.
func (t *Tree) deleteNode(node *Node, key int) {
	i := 0
	for i < node.Size && t.Comparator(node.Keys[i], key) < 0 {
		i++
	}
	if i < node.Size && t.Comparator(node.Keys[i], key) == 0 {
		if node.IsLeaf {
			t.Logger.Infof("deleteNode: deleting key %d from leaf %v", key, node.Keys)
			node.Keys = append(node.Keys[:i], node.Keys[i+1:]...)
			node.Values = append(node.Values[:i], node.Values[i+1:]...)
			node.Size--
			t.checkInvariants(node)
			if node.Size < node.MinKeys {
				t.balanceAfterDeletion(node)
			}
		} else {
			t.Logger.Infof("deleteNode: deleting key %d from internal node %v", key, node.Keys)
			t.deleteInternal(node, i)
		}
	} else if !node.IsLeaf {
		t.deleteNode(node.Children[i], key)
	}
	t.checkInvariants(node)
}

// deleteInternal deletes a key from an internal node.
func (t *Tree) deleteInternal(node *Node, index int) {
	key := node.Keys[index]
	leftChild := node.Children[index]
	rightChild := node.Children[index+1]

	t.Logger.Infof("deleteInternal: deleting key %d at index %d from node %v", key, index, node.Keys)

	// Case 1: Borrow predecessor.
	if leftChild.Size > leftChild.MinKeys {
		predNode, predKey, predValue := t.getPredecessor(leftChild)
		node.Keys[index] = predKey
		node.Values[index] = predValue
		t.deleteNode(predNode, predKey)
		t.checkInvariants(node)
		return
	}

	// Case 2: Borrow successor.
	if rightChild.Size > rightChild.MinKeys {
		succNode, succKey, succValue := t.getSuccessor(rightChild)
		node.Keys[index] = succKey
		node.Values[index] = succValue
		t.deleteNode(succNode, succKey)
		t.checkInvariants(node)
		return
	}

	t.Logger.Infof("deleteInternal: merging children for key %d at index %d", key, index)
	// Case 3: Merge with sibling.
	if node.Parent == nil || node.Parent.Size == 0 {
		t.mergeChildren(node, index)
		t.checkInvariants(node)
		return
	}
	parent := node.Parent
	childIndex := -1
	for i, child := range parent.Children {
		if child == node {
			childIndex = i
			break
		}
	}
	if childIndex < 0 || childIndex >= len(parent.Children) {
		t.Logger.Panicf("deleteInternal: node not found in parent's children")
	}
	if childIndex > 0 {
		t.mergeWithLeftSibling(parent, childIndex-1)
		t.deleteNode(parent.Children[childIndex-1], key)
	} else {
		t.mergeWithRightSibling(parent, childIndex)
		t.deleteNode(parent.Children[childIndex], key)
	}
	t.checkInvariants(node)
}

// Helper functions
// getPredecessor finds the predecessor key in the subtree rooted at the given node.
func (t *Tree) getPredecessor(node *Node) (*Node, int, interface{}) {
	current := node
	for !current.IsLeaf {
		current = current.Children[current.Size] // Traverse to the last child
	}
	predKey := current.Keys[current.Size-1]
	predValue := current.Values[current.Size-1]
	return current, predKey, predValue
}

// getSuccessor finds the successor key in the subtree rooted at the given node.
func (t *Tree) getSuccessor(node *Node) (*Node, int, interface{}) {
	current := node
	for !current.IsLeaf {
		current = current.Children[0] // Traverse to the first child
	}
	succKey := current.Keys[0]
	succValue := current.Values[0]
	return current, succKey, succValue
}

// PrintTree prints the tree structure (for debugging).
func (t *Tree) PrintTree() {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	t.printNode(t.Root, 0)
}

// printNode prints a subtree rooted at the given node.
func (t *Tree) printNode(node *Node, level int) {
	if node == nil {
		return
	}

	t.Logger.Infof("Level %d: %v", level, node.Keys)
	for _, child := range node.Children {
		t.printNode(child, level+1)
	}
}

func (t *Tree) SetLogger(logger *logger.Logger) {
	t.Logger = logger
}

// RebuildParentPointers recursively sets the Parent pointers for all nodes.
func (t *Tree) RebuildParentPointers() {
    if t.Root == nil {
        return
    }
   // t.Root.Parent = nil // Root has no parent
    t.rebuildParents(t.Root)
}

func (t *Tree) rebuildParents(node *Node) {
    for _, child := range node.Children {
        child.Parent = node
        t.rebuildParents(child)
    }
}

func (t *Tree) mergeChildren(node *Node, index int) {
    if node == nil {
        t.Logger.Panicf("mergeChildren: node is nil")
    }
    t.Logger.Infof("mergeChildren: BEFORE merge at index %d, node keys: %v, children count: %d",
        index, node.Keys, len(node.Children))
    if index < 0 || index >= len(node.Children)-1 {
        t.Logger.Panicf("mergeChildren: invalid index %d (node.Children length %d)", index, len(node.Children))
    }
    leftChild := node.Children[index]
    rightChild := node.Children[index+1]

    t.Logger.Infof("mergeChildren: BEFORE merge, leftChild keys: %v, rightChild keys: %v",
        leftChild.Keys, rightChild.Keys)

    // Move the separator key from node to leftChild.
    leftChild.Keys = append(leftChild.Keys, node.Keys[index])
    leftChild.Values = append(leftChild.Values, node.Values[index])
    leftChild.Size = len(leftChild.Keys)

    // Merge rightChild's keys and values into leftChild.
    leftChild.Keys = append(leftChild.Keys, rightChild.Keys...)
    leftChild.Values = append(leftChild.Values, rightChild.Values...)
    leftChild.Size = len(leftChild.Keys)

    // Merge children pointers if not a leaf.
    if !leftChild.IsLeaf {
        leftChild.Children = append(leftChild.Children, rightChild.Children...)
        for _, child := range rightChild.Children {
            child.Parent = leftChild
        }
    }

    // Remove the separator key and the pointer for rightChild from node.
    node.Keys = append(node.Keys[:index], node.Keys[index+1:]...)
    node.Values = append(node.Values[:index], node.Values[index+1:]...)
    // Rebuild node.Children with a new allocation.
    newChildren := make([]*Node, 0, len(node.Children)-1)
    newChildren = append(newChildren, node.Children[:index+1]...)
    newChildren = append(newChildren, node.Children[index+2:]...)
    node.Children = newChildren
    node.Size = len(node.Keys)

    // Invariant: len(node.Children) should equal node.Size + 1.
    expectedChildrenCount := node.Size + 1
    if len(node.Children) != expectedChildrenCount {
        t.Logger.Panicf("mergeChildren: invariant violation: node's children count %d != keys count %d + 1",
            len(node.Children), node.Size)
    }

    t.Logger.Infof("mergeChildren: AFTER merge, node keys: %v, children count: %d",
        node.Keys, len(node.Children))
    t.checkInvariants(node)

    // If node is the root and becomes empty, update the tree’s root.
    if node.Size == 0 && node.Parent == nil {
        t.Logger.Infof("mergeChildren: node is root and empty, replacing root with leftChild keys: %v",
            leftChild.Keys)
        t.Root = leftChild
        leftChild.Parent = nil
        t.Height--
    }
}

// normalizeChildren rebuilds a node's children slice so that its length equals node.Size+1.
func (t *Tree) normalizeChildren(node *Node) {
	if node.IsLeaf {
		return
	}
	expected := node.Size + 1
	if len(node.Children) != expected {
		t.Logger.Infof("normalizeChildren: normalizing node; expected %d children, got %d", expected, len(node.Children))
		newChildren := make([]*Node, expected)
		copy(newChildren, node.Children[:expected])
		node.Children = newChildren
	}
}