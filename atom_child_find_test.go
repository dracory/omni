package omni

import "testing"

func TestChildFindByID_FindsImmediateChild(t *testing.T) {
	parent := NewAtom("parent", WithID("p"))
	c1 := NewAtom("child", WithID("c1"))
	c2 := NewAtom("child", WithID("c2"))
	parent.ChildAdd(c1)
	parent.ChildAdd(c2)

	got := parent.ChildFindByID("c2")
	if got == nil || got.GetID() != "c2" {
		t.Fatalf("expected to find child with ID c2, got %+v", got)
	}
}

func TestChildFindByID_NotFound(t *testing.T) {
	parent := NewAtom("parent", WithID("p"))
	parent.ChildAdd(NewAtom("child", WithID("c1")))

	got := parent.ChildFindByID("missing")
	if got != nil {
		t.Fatalf("expected nil when not found, got ID=%s", got.GetID())
	}
}

func TestChildFindByID_IgnoresGrandchildren(t *testing.T) {
	parent := NewAtom("parent", WithID("p"))
	child := NewAtom("child", WithID("c1"))
	grand := NewAtom("grand", WithID("g1"))
	child.ChildAdd(grand)
	parent.ChildAdd(child)

	got := parent.ChildFindByID("g1")
	if got != nil {
		t.Fatalf("expected nil for non-immediate child search, got ID=%s", got.GetID())
	}
}

func TestChildrenFindByType_MultipleImmediateMatches(t *testing.T) {
	parent := NewAtom("parent")
	// matching type "item"
	c1 := NewAtom("item", WithID("i1"))
	c2 := NewAtom("item", WithID("i2"))
	// non-matching type
	c3 := NewAtom("other", WithID("o1"))
	parent.ChildrenAdd([]AtomInterface{c1, c2, c3})

	got := parent.ChildrenFindByType("item")
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(got))
	}
	// verify order preserved (insertion order)
	if got[0].GetID() != "i1" || got[1].GetID() != "i2" {
		t.Fatalf("unexpected order of results: %s, %s", got[0].GetID(), got[1].GetID())
	}
}

func TestChildrenFindByType_NoMatches(t *testing.T) {
	parent := NewAtom("parent")
	parent.ChildrenAdd([]AtomInterface{
		NewAtom("x", WithID("x1")),
		NewAtom("y", WithID("y1")),
	})

	got := parent.ChildrenFindByType("z")
	if len(got) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(got))
	}
}
