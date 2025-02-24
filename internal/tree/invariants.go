package tree

// checkInvariants recursively asserts that each non-leaf node has one more child than its key count.
func (t *Tree) checkInvariants(node *Node) {
    if node == nil {
        return
    }
    if !node.IsLeaf {
        if len(node.Children) != node.Size+1 {
            t.Logger.Panicf("Invariant violation: node %v has %d keys but %d children (expected %d)",
                node.Keys, node.Size, len(node.Children), node.Size+1)
        }
        for _, child := range node.Children {
            t.checkInvariants(child)
        }
    }
}
