package omni

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewAtom(t *testing.T) {
	a := NewAtom("test-type", WithID("test-id"))
	if a == nil {
		t.Error("NewAtom returned nil")
	}
	if got := a.GetID(); got != "test-id" {
		t.Errorf("NewAtom() ID = %v, want %v", got, "test-id")
	}
	if got := a.GetType(); got != "test-type" {
		t.Errorf("NewAtom() type = %v, want %v", got, "test-type")
	}
	if got := a.Get("test"); got != "" {
		t.Errorf("NewAtom() Get(\"test\") = %v, want empty string", got)
	}
	if children := a.ChildrenGet(); len(children) != 0 {
		t.Errorf("NewAtom() should have no children, got %d", len(children))
	}
}

func TestAtom_PropertyManagement(t *testing.T) {
	a := NewAtom("test", WithID("test-id"))

	// Test Set and Get
	a.Set("prop1", "value1")
	if got := a.Get("prop1"); got != "value1" {
		t.Errorf("SetProperty(\"prop1\", \"value1\") = %v, want %v", got, "value1")
	}

	// Test updating property
	a.Set("prop1", "updated")
	if got := a.Get("prop1"); got != "updated" {
		t.Errorf("SetProperty(\"prop1\", \"updated\") = %v, want %v", got, "updated")
	}

	// Test getting non-existent property
	if got := a.Get("nonexistent"); got != "" {
		t.Error("Expected empty string for non-existent property, got:", got)
	}
}

func TestAtom_ChildManagement(t *testing.T) {
	parent := NewAtom("container", WithID("parent"))
	child1 := NewAtom("item", WithID("child1"))
	child2 := NewAtom("item", WithID("child2"))

	// Test ChildAdd and ChildrenGet
	parent.ChildAdd(child1).ChildAdd(child2)

	children := parent.ChildrenGet()
	if len(children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(children))
	}
	if children[0].GetID() != "child1" || children[1].GetID() != "child2" {
		t.Errorf("Children not added in correct order. Got: %v, %v", children[0].GetID(), children[1].GetID())
	}
}

func TestAtom_ConcurrentAccess(t *testing.T) {
	a := NewAtom("test", WithID("concurrent"))
	var wg sync.WaitGroup

	// Start multiple goroutines to modify the atom
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				// Test property operations
				a.Set("prop", "value")
				_ = a.Get("prop")
				a.Remove("prop") // Clear property

				// Test child operations
				child := NewAtom("item", WithID(fmt.Sprintf("child-%d-%d", id, j)))
				a.ChildAdd(child)
				_ = a.ChildrenGet()
			}
		}(i)
	}

	wg.Wait()

	// Verify final state
	if got := a.Get("prop"); got != "" {
		t.Error("Property should not exist after concurrent operations")
	}
	// At least some children should have been added
	if len(a.ChildrenGet()) == 0 {
		t.Error("Expected some children after concurrent operations")
	}
}

func TestAtom_EdgeCases(t *testing.T) {
	a := NewAtom("test", WithID("edge"))

	// Test empty key
	a.Set("", "empty-key")
	if got := a.Get(""); got != "empty-key" {
		t.Error("Should handle empty key")
	}

	// Test setting empty value
	a.Set("empty-value", "")
	if got := a.Get("empty-value"); got != "" {
		t.Error("Should handle empty value")
	}

	// Test getting non-existent key
	if got := a.Get("nonexistent"); got != "" {
		t.Error("Should return empty string for non-existent key")
	}
}

func TestAtom_AddChildren(t *testing.T) {
	t.Run("Add multiple children", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		child1 := NewAtom("item", WithID("child1"))
		child2 := NewAtom("item", WithID("child2"))
		parent.ChildAdd(child1).ChildAdd(child2)
		children := parent.ChildrenGet()
		if len(children) != 2 {
			t.Fatalf("Expected 2 children, got %d", len(children))
		}
		if children[0].GetID() != "child1" || children[1].GetID() != "child2" {
			t.Error("Children not added in correct order")
		}
	})

	t.Run("Add to existing children", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		existingChild := NewAtom("item", WithID("existing"))
		parent.ChildAdd(existingChild)
		child1 := NewAtom("item", WithID("child1"))
		child2 := NewAtom("item", WithID("child2"))
		parent.ChildAdd(child1).ChildAdd(child2)

		children := parent.ChildrenGet()
		if len(children) != 3 {
			t.Fatalf("Expected 3 children, got %d", len(children))
		}
		if children[0].GetID() != "existing" || children[1].GetID() != "child1" || children[2].GetID() != "child2" {
			t.Error("Children not added correctly to existing children")
		}
	})

	t.Run("Add nil child", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		// Adding nil should be a no-op
		parent.ChildAdd(nil)
		if parent.ChildrenLength() != 0 {
			t.Error("Adding nil child should be a no-op")
		}

		// Add a real child
		child := NewAtom("item", WithID("child"))
		parent.ChildAdd(child)
		children := parent.ChildrenGet()
		if parent.ChildrenLength() != 1 || children[0].GetID() != "child" {
			t.Error("Failed to add valid child")
		}
	})

	t.Run("Set children", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		child1 := NewAtom("item", WithID("child1"))
		child2 := NewAtom("item", WithID("child2"))

		parent.ChildrenSet([]AtomInterface{child1, child2})

		children := parent.ChildrenGet()
		if len(children) != 2 {
			t.Fatalf("Expected 2 children, got %d", len(children))
		}
		if children[0].GetID() != "child1" || children[1].GetID() != "child2" {
			t.Error("Children not set correctly")
		}

		// Test setting to nil
		parent.ChildrenSet(nil)
		if parent.ChildrenLength() != 0 {
			t.Error("Setting children to nil should clear children")
		}
	})
}

func TestAtom_SetChildren(t *testing.T) {
	t.Run("Set children on empty parent", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		child1 := NewAtom("item", WithID("child1"))
		child2 := NewAtom("item", WithID("child2"))
		parent.ChildrenSet([]AtomInterface{child1, child2})
		children := parent.ChildrenGet()

		if len(children) != 2 {
			t.Fatalf("Expected 2 children, got %d", len(children))
		}
		if children[0].GetID() != "child1" || children[1].GetID() != "child2" {
			t.Error("Children not set correctly")
		}
	})

	t.Run("Replace existing children", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		oldChild := NewAtom("item", WithID("old"))
		parent.ChildAdd(oldChild)

		newChild1 := NewAtom("item", WithID("new1"))
		newChild2 := NewAtom("item", WithID("new2"))
		parent.ChildrenSet([]AtomInterface{newChild1, newChild2})

		children := parent.ChildrenGet()
		if len(children) != 2 {
			t.Fatalf("Expected 2 children, got %d", len(children))
		}
		if children[0].GetID() != "new1" || children[1].GetID() != "new2" {
			t.Error("Children not replaced correctly")
		}
	})

	t.Run("Set empty children", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		child := NewAtom("item", WithID("child"))
		parent.ChildAdd(child)

		parent.ChildrenSet([]AtomInterface{})
		if parent.ChildrenLength() != 0 {
			t.Error("Expected no children after setting to empty slice")
		}

		parent.ChildrenSet(nil)
		if parent.ChildrenLength() != 0 {
			t.Error("Expected no children after setting to nil")
		}
	})

	t.Run("Set with nil children", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		child := NewAtom("item", WithID("child"))

		parent.ChildrenSet([]AtomInterface{child, nil, nil})

		children := parent.ChildrenGet()
		if len(children) != 1 {
			t.Fatalf("Expected 1 child, got %d", len(children))
		}
		if children[0].GetID() != "child" {
			t.Error("Nil children should be filtered out")
		}
	})
}

func TestAtom_ConcurrentChildOperations(t *testing.T) {
	parent := NewAtom("container", WithID("parent"))
	var wg sync.WaitGroup

	// Test concurrent ChildAdd
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			child := NewAtom("item", WithID(fmt.Sprintf("child-%d", id)))
			parent.ChildAdd(child)
		}(i)
	}

	// Test concurrent ChildrenSet
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			children := []AtomInterface{
				NewAtom("item", WithID(fmt.Sprintf("new-child1-%d", id))),
				NewAtom("item", WithID(fmt.Sprintf("new-child2-%d", id))),
			}
			parent.ChildrenSet(children)
		}(i)
	}

	// Test concurrent ChildrenGet
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			_ = parent.ChildrenGet()
		}()
	}

	wg.Wait()

	// Verify some invariants
	children := parent.ChildrenGet()
	if len(children) == 0 {
		t.Error("Expected some children after concurrent operations")
	}

	// Check for duplicates (shouldn't happen with atomic operations)
	seen := make(map[string]bool)
	for _, child := range children {
		if seen[child.GetID()] {
			t.Errorf("Duplicate child ID found: %s", child.GetID())
		}
		seen[child.GetID()] = true
	}
}
