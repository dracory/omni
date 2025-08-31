package omni

import "testing"

func TestFindAtomsByType_MultipleMatches(t *testing.T) {
	// Build tree:
	// root(root) -> a(node) -> a1(target), a2(node)
	//            -> b(target) -> b1(target)
	root := NewAtom("root", WithID("root"))
	a := NewAtom("node", WithID("a"))
	a1 := NewAtom("target", WithID("a1"))
	a2 := NewAtom("node", WithID("a2"))
	b := NewAtom("target", WithID("b"))
	b1 := NewAtom("target", WithID("b1"))

	a.ChildAdd(a1).ChildAdd(a2)
	b.ChildAdd(b1)
	root.ChildrenSet([]AtomInterface{a, b})

	matches := FindAtomsByType(root, "target")
	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}
	// Pre-order expected order: a1, b, b1
	if matches[0].GetID() != "a1" || matches[1].GetID() != "b" || matches[2].GetID() != "b1" {
		t.Fatalf("unexpected order: got [%s, %s, %s]", matches[0].GetID(), matches[1].GetID(), matches[2].GetID())
	}
}

func TestFindAtomsByType_NoMatches(t *testing.T) {
	root := NewAtom("root", WithID("r"))
	root.ChildAdd(NewAtom("child", WithID("c1")))
	root.ChildAdd(NewAtom("child", WithID("c2")))

	matches := FindAtomsByType(root, "target")
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(matches))
	}
}

func TestFindAtomsByType_NilRoot(t *testing.T) {
	var root AtomInterface
	matches := FindAtomsByType(root, "any")
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches for nil root, got %d", len(matches))
	}
}

func TestFindAtomsByType_RootOnly(t *testing.T) {
	root := NewAtom("target", WithID("r"))
	matches := FindAtomsByType(root, "target")
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].GetID() != "r" {
		t.Fatalf("expected root to be the match, got %q", matches[0].GetID())
	}
}
