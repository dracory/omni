package omni

// FindAtomsByType recursively finds all atoms of a specific type in a tree.
// It performs a pre-order traversal and returns matches in that order.
func FindAtomsByType(root AtomInterface, atomType string) []AtomInterface {
	result := []AtomInterface{}
	if root == nil {
		return result
	}

	// Check current atom first (pre-order)
	if root.GetType() == atomType {
		result = append(result, root)
	}

	// Recursively collect from children
	for _, child := range root.ChildrenGet() {
		result = append(result, FindAtomsByType(child, atomType)...)
	}

	return result
}
