package omni

import "testing"

func TestFindAtomByID_FindsRoot(t *testing.T) {
	root := NewAtom("root", WithID("root"))
	// Sanity: root has children but function should match root first
	child := NewAtom("child", WithID("c1"))
	root.ChildAdd(child)

	got := FindAtomByID(root, "root")
	if got == nil || got.GetID() != "root" {
		t.Fatalf("expected to find root with id 'root', got %#v", got)
	}
}

func TestFindAtomByID_FindsDeepChild(t *testing.T) {
	root := NewAtom("root", WithID("root"))
	level1a := NewAtom("node", WithID("a"))
	level1b := NewAtom("node", WithID("b"))
	level2a := NewAtom("node", WithID("a1"))
	level2b := NewAtom("node", WithID("b1"))
	level3 := NewAtom("leaf", WithID("target"))

	level2a.ChildAdd(level3)
	level1a.ChildAdd(level2a)
	level1b.ChildAdd(level2b)
	root.ChildrenSet([]AtomInterface{level1a, level1b})

	got := FindAtomByID(root, "target")
	if got == nil {
		t.Fatal("expected to find deep child, got nil")
	}
	if got.GetID() != "target" {
		t.Fatalf("expected id 'target', got %q", got.GetID())
	}
}

func TestFindAtomByID_NotFound(t *testing.T) {
	root := NewAtom("root", WithID("root"))
	root.ChildAdd(NewAtom("child", WithID("c1")))
	root.ChildAdd(NewAtom("child", WithID("c2")))

	got := FindAtomByID(root, "missing")
	if got != nil {
		t.Fatalf("expected nil when ID not found, got %#v", got)
	}
}

func TestFindAtomByID_NilRoot(t *testing.T) {
	var root AtomInterface
	got := FindAtomByID(root, "any")
	if got != nil {
		t.Fatalf("expected nil when root is nil, got %#v", got)
	}
}
