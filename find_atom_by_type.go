package omni

// FindAtomByType recursively finds the first atom with the given type in a tree.
// It performs a pre-order traversal: checks the current node first, then its children in order.
func FindFirstAtomByType(root AtomInterface, atomType string) AtomInterface {
	if root == nil {
		return nil
	}

	if root.GetType() == atomType {
		return root
	}

	for _, child := range root.ChildrenGet() {
		if found := FindFirstAtomByType(child, atomType); found != nil {
			return found
		}
	}

	return nil
}
