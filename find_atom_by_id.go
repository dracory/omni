package omni

// FindAtomByID recursively finds an atom by ID in a tree.
// It performs a pre-order traversal: checks the current node first,
// then descends into its children in order.
func FindAtomByID(root AtomInterface, id string) AtomInterface {
	if root == nil {
		return nil
	}

	if root.GetID() == id {
		return root
	}

	for _, child := range root.ChildrenGet() {
		if found := FindAtomByID(child, id); found != nil {
			return found
		}
	}

	return nil
}
