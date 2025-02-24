package tree

// balanceAfterDeletion checks if a node is underfilled and rebalances the tree.
func (t *Tree) balanceAfterDeletion(node *Node) {
	// Add nil check for logger
    if t.Logger == nil {
        panic("logger not initialized")
    }
    if node == nil {
        return
    }

    // If node has enough keys or is root, no balancing needed
    if node.Size >= node.MinKeys || node.Parent == nil {
        // Handle root becoming empty
        if node.Parent == nil && node.Size == 0 && len(node.Children) > 0 {
            t.Root = node.Children[0]
            t.Root.Parent = nil
            t.Height--
        }
        return
    }

    parent := node.Parent
    childIndex := -1

    // Find the index of the node in parent's children
    for i, child := range parent.Children {
        if child == node {
            childIndex = i
            break
        }
    }

    // Validate child index
    if childIndex < 0 || childIndex >= len(parent.Children) {
        t.Logger.Panicf("invalid child index %d in parent with %d children",
            childIndex, len(parent.Children))
    }

    // Try to borrow from left sibling
    if childIndex > 0 {
        leftSibling := parent.Children[childIndex-1]
        if leftSibling.Size > leftSibling.MinKeys {
            t.borrowFromLeftSibling(parent, childIndex)
            return
        }
    }

    // Try to borrow from right sibling
    if childIndex < len(parent.Children)-1 {
        rightSibling := parent.Children[childIndex+1]
        if rightSibling.Size > rightSibling.MinKeys {
            t.borrowFromRightSibling(parent, childIndex)
            return
        }
    }

    // Merge with a sibling if borrowing failed
    if childIndex > 0 {
        t.mergeWithLeftSibling(parent, childIndex-1)
    } else {
        t.mergeWithRightSibling(parent, childIndex)
    }

    // Propagate balancing to parent if necessary
    if parent.Size < parent.MinKeys {
        t.balanceAfterDeletion(parent)
    }
}

// borrowFromLeftSibling borrows a key from the left sibling.
func (t *Tree) borrowFromLeftSibling(parent *Node, index int) {
    if parent == nil || index <= 0 || index >= len(parent.Children) {
        t.Logger.Panicf("borrowFromLeftSibling: invalid parameters, index %d, parent.Children %d", index, len(parent.Children))
    }
    node := parent.Children[index]
    leftSibling := parent.Children[index-1]

    t.Logger.Infof("borrowFromLeftSibling: before borrowing, node keys: %v, leftSibling keys: %v", node.Keys, leftSibling.Keys)

    // Move parent's key down to node.
    node.Keys = append([]int{parent.Keys[index-1]}, node.Keys...)
    node.Values = append([]interface{}{parent.Values[index-1]}, node.Values...)
    node.Size++

    // Move left sibling's last key to parent.
    parent.Keys[index-1] = leftSibling.Keys[leftSibling.Size-1]
    parent.Values[index-1] = leftSibling.Values[leftSibling.Size-1]

    // Remove borrowed key from left sibling.
    leftSibling.Keys = leftSibling.Keys[:leftSibling.Size-1]
    leftSibling.Values = leftSibling.Values[:leftSibling.Size-1]
    leftSibling.Size--

    if !node.IsLeaf {
        borrowedChild := leftSibling.Children[leftSibling.Size]
        node.Children = append([]*Node{borrowedChild}, node.Children...)
        borrowedChild.Parent = node
        leftSibling.Children = leftSibling.Children[:leftSibling.Size+1]
    }

    t.Logger.Infof("borrowFromLeftSibling: after borrowing, node keys: %v, leftSibling keys: %v", node.Keys, leftSibling.Keys)
    t.checkInvariants(parent)
}


// borrowFromRightSibling borrows a key from the right sibling.
func (t *Tree) borrowFromRightSibling(parent *Node, index int) {
	node := parent.Children[index]
	rightSibling := parent.Children[index+1]

	// Move parent's key down to node
	node.Keys = append(node.Keys, parent.Keys[index])
	node.Values = append(node.Values, parent.Values[index])
	node.Size++

	// Move right sibling's first key to parent
	parent.Keys[index] = rightSibling.Keys[0]
	parent.Values[index] = rightSibling.Values[0]

	// Remove borrowed key from right sibling
	rightSibling.Keys = rightSibling.Keys[1:]
	rightSibling.Values = rightSibling.Values[1:]
	rightSibling.Size--

	// Move child pointer if not a leaf
	if !node.IsLeaf {
		// Take right sibling's first child
		borrowedChild := rightSibling.Children[0]
		node.Children = append(node.Children, borrowedChild)
		borrowedChild.Parent = node

		// Remove child from right sibling
		rightSibling.Children = rightSibling.Children[1:]
	}
}

// mergeWithLeftSibling merges the node with its left sibling.
func (t *Tree) mergeWithLeftSibling(parent *Node, leftIndex int) {
    if parent == nil {
        t.Logger.Panicf("mergeWithLeftSibling: parent is nil")
    }
    t.Logger.Infof("mergeWithLeftSibling: BEFORE merge, leftIndex: %d, parent keys: %v, children count: %d",
        leftIndex, parent.Keys, len(parent.Children))
    originalParentSize := len(parent.Keys)
    if leftIndex < 0 || leftIndex >= originalParentSize {
        t.Logger.Panicf("mergeWithLeftSibling: invalid merge index %d (parent has %d keys)", leftIndex, originalParentSize)
    }
    if leftIndex+1 >= len(parent.Children) {
        t.Logger.Panicf("mergeWithLeftSibling: not enough children, leftIndex %d, children length %d", leftIndex, len(parent.Children))
    }

    leftSibling := parent.Children[leftIndex]
    node := parent.Children[leftIndex+1]

    t.Logger.Infof("mergeWithLeftSibling: BEFORE merge, leftSibling keys: %v, node keys: %v", leftSibling.Keys, node.Keys)

    // Move the separator key from parent to leftSibling.
    leftSibling.Keys = append(leftSibling.Keys, parent.Keys[leftIndex])
    leftSibling.Values = append(leftSibling.Values, parent.Values[leftIndex])
    leftSibling.Size = len(leftSibling.Keys)

    // Merge node's keys and values into leftSibling.
    leftSibling.Keys = append(leftSibling.Keys, node.Keys...)
    leftSibling.Values = append(leftSibling.Values, node.Values...)
    leftSibling.Size = len(leftSibling.Keys)

    // Merge children if node is not a leaf.
    if !node.IsLeaf {
        leftSibling.Children = append(leftSibling.Children, node.Children...)
        for _, child := range node.Children {
            child.Parent = leftSibling
        }
    }

    // Remove the separator key and the merged child (node) from parent.
    parent.Keys = append(parent.Keys[:leftIndex], parent.Keys[leftIndex+1:]...)
    parent.Values = append(parent.Values[:leftIndex], parent.Values[leftIndex+1:]...)
    // Instead of using append directly on parent.Children, create a new slice.
    newChildren := make([]*Node, 0, len(parent.Children)-1)
    newChildren = append(newChildren, parent.Children[:leftIndex+1]...)
    newChildren = append(newChildren, parent.Children[leftIndex+2:]...)
    parent.Children = newChildren

    // Ensure leftSibling remains at the correct position.
    parent.Children[leftIndex] = leftSibling
    parent.Size = len(parent.Keys)

    // Extra trimming (should be redundant now, but we check):
    expectedChildrenCount := parent.Size + 1
    if len(parent.Children) != expectedChildrenCount {
        t.Logger.Panicf("mergeWithLeftSibling: invariant violation: parent's children count %d != keys count %d + 1",
            len(parent.Children), parent.Size)
    }

    t.Logger.Infof("mergeWithLeftSibling: AFTER merge, parent keys: %v, children count: %d",
        parent.Keys, len(parent.Children))
    t.checkInvariants(parent)

    // If parent is the root and becomes empty, update the tree’s root.
    if parent.Parent == nil && parent.Size == 0 {
        t.Logger.Infof("mergeWithLeftSibling: parent is root and empty, replacing root with leftSibling keys: %v",
            leftSibling.Keys)
        t.Root = leftSibling
        leftSibling.Parent = nil
        t.Height--
    }
}


// mergeWithRightSibling merges the node with its right sibling.
func (t *Tree) mergeWithRightSibling(parent *Node, index int) {
    if parent == nil {
        t.Logger.Panicf("mergeWithRightSibling: parent is nil")
    }
    t.Logger.Infof("mergeWithRightSibling: BEFORE merge, index: %d, parent keys: %v, children count: %d",
        index, parent.Keys, len(parent.Children))
    parent.Size = len(parent.Keys)
    if parent.Size == 0 {
        t.Logger.Panicf("mergeWithRightSibling: parent has 0 keys, cannot merge")
    }
    if index < 0 || index >= parent.Size {
        t.Logger.Panicf("mergeWithRightSibling: invalid index %d (parent has %d keys)", index, parent.Size)
    }
    if index+1 >= len(parent.Children) {
        t.Logger.Panicf("mergeWithRightSibling: not enough children, index %d, children length %d", index, len(parent.Children))
    }

    node := parent.Children[index]
    rightSibling := parent.Children[index+1]

    t.Logger.Infof("mergeWithRightSibling: BEFORE merge, node keys: %v, rightSibling keys: %v",
        node.Keys, rightSibling.Keys)

    // Move parent's key down to node.
    node.Keys = append(node.Keys, parent.Keys[index])
    node.Values = append(node.Values, parent.Values[index])
    node.Size = len(node.Keys)

    // Merge rightSibling's keys and values into node.
    node.Keys = append(node.Keys, rightSibling.Keys...)
    node.Values = append(node.Values, rightSibling.Values...)
    node.Size = len(node.Keys)

    // Merge children if node is not a leaf.
    if !node.IsLeaf {
        node.Children = append(node.Children, rightSibling.Children...)
        for _, child := range rightSibling.Children {
            child.Parent = node
        }
    }

    // Remove the merged key and rightSibling pointer from parent.
    parent.Keys = append(parent.Keys[:index], parent.Keys[index+1:]...)
    parent.Values = append(parent.Values[:index], parent.Values[index+1:]...)
    // Rebuild parent's children slice with a new allocation.
    newChildren := make([]*Node, 0, len(parent.Children)-1)
    newChildren = append(newChildren, parent.Children[:index+1]...)
    newChildren = append(newChildren, parent.Children[index+2:]...)
    parent.Children = newChildren

    parent.Size = len(parent.Keys)

    // Invariant: len(parent.Children) should equal parent.Size + 1.
    expectedChildrenCount := parent.Size + 1
    if len(parent.Children) != expectedChildrenCount {
        t.Logger.Panicf("mergeWithRightSibling: invariant violation: parent's children count %d != keys count %d + 1",
            len(parent.Children), parent.Size)
    }

    t.Logger.Infof("mergeWithRightSibling: AFTER merge, parent keys: %v, children count: %d",
        parent.Keys, len(parent.Children))
    t.checkInvariants(parent)
}

