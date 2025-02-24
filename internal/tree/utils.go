package tree

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PrintTreeStructure prints the tree in a human-readable format (level-order traversal).
func (t *Tree) PrintTreeStructure() {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	if t.Root == nil {
		t.Logger.Infof("Tree is empty")
		return
	}

	queue := []*Node{t.Root}
	level := 0
	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			// Log node keys using the custom logger
			t.Logger.Infof("Level %d: %v", level, node.Keys) // Use `level`

			// Add children to the queue
			if !node.IsLeaf {
				queue = append(queue, node.Children...)
			}

		}
		level++ // Increment after processing each level
	}
}

// ValidateTree checks if the tree adheres to B-tree properties.
func (t *Tree) ValidateTree() bool {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	if t.Root == nil {
		return true // An empty tree is valid
	}

	return t.validateNode(t.Root, true)
}

// validateNode recursively checks if a node and its children adhere to B-tree properties.
func (t *Tree) validateNode(node *Node, isRoot bool) bool {
	// Check key count
	if !isRoot && (node.Size < node.MinKeys || node.Size > node.MaxKeys) {
		t.Logger.Errorf("Invalid node: key count %d is outside range [%d, %d]\n", node.Size, node.MinKeys, node.MaxKeys)
		return false
	}

	// Check if keys are sorted
	for i := 1; i < node.Size; i++ {
		if t.Comparator(node.Keys[i-1], node.Keys[i]) >= 0 {
			t.Logger.Errorf("Invalid node: keys are not sorted (%v)\n", node.Keys)
			return false
		}
	}

	// Recursively validate children
	if !node.IsLeaf {
		for _, child := range node.Children {
			if child.Parent != node {
				t.Logger.Errorf("Invalid node: parent pointer mismatch\n")
				return false
			}
			if !t.validateNode(child, false) {
				return false
			}
		}
	}

	return true
}

// SerializeTree converts the tree to a JSON string for storage or transmission.
func (t *Tree) SerializeTree() (string, error) {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	data, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("failed to serialize tree: %v", err)
	}
	return string(data), nil
}

// DeserializeTree loads a tree from a JSON string.
func (t *Tree) DeserializeTree(data string) error {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	err := json.Unmarshal([]byte(data), t)
	if err != nil {
		return fmt.Errorf("failed to deserialize tree: %v", err)
	}
	return nil
}

// ToString returns a string representation of the tree (for debugging).
func (t *Tree) ToString() string {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Tree (degree=%d, size=%d, height=%d):\n", t.Degree, t.Size, t.Height))
	t.printNodeToString(&buffer, t.Root, 0)
	return buffer.String()
}

// printNodeToString recursively writes node information to a buffer.
func (t *Tree) printNodeToString(buffer *bytes.Buffer, node *Node, level int) {
	if node == nil {
		return
	}

	buffer.WriteString(fmt.Sprintf("Level %d: %v\n", level, node.Keys))
	for _, child := range node.Children {
		t.printNodeToString(buffer, child, level+1)
	}
}