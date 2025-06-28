package omni

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewAtom(t *testing.T) {
	a := NewAtom("test-type")
	if a == nil {
		t.Error("NewAtom returned nil")
	}
	if a.GetID() == "" {
		t.Error("NewAtom did not generate an ID")
	}
	if a.GetType() != "test-type" {
		t.Errorf("Expected type 'test-type', got '%s'", a.GetType())
	}
	if len(a.GetProperties()) != 0 {
		t.Error("NewAtom should not have any properties by default")
	}
	if len(a.GetChildren()) != 0 {
		t.Error("NewAtom should not have any children by default")
	}
}

func TestAtom_PropertyManagement(t *testing.T) {
	prop1 := NewProperty("prop1", "value1")
	a := NewAtom("test", WithProperties(prop1))

	// Test SetProperty and GetProperty
	if got := a.GetProperty("prop1"); got == nil || got.GetValue() != "value1" {
		t.Error("Failed to set/get property")
	}

	// Test updating property
	prop1Updated := NewProperty("prop1", "updated")
	a.SetProperty(prop1Updated)
	if got := a.GetProperty("prop1").GetValue(); got != "updated" {
		t.Errorf("GetProperty().GetValue() = %v, want %v", got, "updated")
	}

	// Test SetProperties and GetProperties
	prop2 := NewProperty("prop2", "value2")
	a.SetProperties([]PropertyInterface{prop1Updated, prop2})
	if got := len(a.GetProperties()); got != 2 {
		t.Errorf("Expected 2 properties, got %d", got)
	}

	// Test RemoveProperty
	a.RemoveProperty("prop1")
	if a.GetProperty("prop1") != nil {
		t.Error("Property was not removed")
	}
}

func TestAtom_ChildManagement(t *testing.T) {
	parent := NewAtom("container", WithID("parent"))
	child1 := NewAtom("item", WithID("child1"))
	child2 := NewAtom("item", WithID("child2"))

	// Test AddChild and GetChildren
	parent.AddChild(child1)
	parent.AddChild(child2)

	children := parent.GetChildren()
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
				prop := NewProperty("prop", "value")
				a.SetProperty(prop)
				_ = a.GetProperty("prop")
				a.RemoveProperty("prop")

				// Test child operations
				child := NewAtom("item", WithID("child"))
				a.AddChild(child)
				_ = a.GetChildren()
			}
		}(i)
	}

	wg.Wait()

	// Verify final state
	if got := a.GetProperty("prop"); got != nil {
		t.Error("Property should not exist after concurrent operations")
	}
	// At least some children should have been added
	if len(a.GetChildren()) == 0 {
		t.Error("Expected some children after concurrent operations")
	}
}

func TestAtom_EdgeCases(t *testing.T) {
	a := NewAtom("test", WithID("edge"))

	// Test nil property
	a.SetProperty(nil)
	if len(a.GetProperties()) != 0 {
		t.Error("Should not add nil property")
	}

	// Test empty property name
	emptyProp := NewProperty("", "value")
	a.SetProperty(emptyProp)
	if got := a.GetProperty("").GetValue(); got != "value" {
		t.Error("Should handle empty property name")
	}

	// Test removing non-existent property
	a.RemoveProperty("nonexistent")
	// Should not panic
}

func TestAtom_AddChildren(t *testing.T) {
	t.Run("Add multiple children", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		child1 := NewAtom("item", WithID("child1"))
		child2 := NewAtom("item", WithID("child2"))
		parent.AddChildren([]AtomInterface{child1, child2})
		children := parent.GetChildren()
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
		parent.AddChild(existingChild)
		child1 := NewAtom("item", WithID("child1"))
		child2 := NewAtom("item", WithID("child2"))
		parent.AddChildren([]AtomInterface{child1, child2})

		children := parent.GetChildren()
		if len(children) != 3 {
			t.Fatalf("Expected 3 children, got %d", len(children))
		}
		if children[0].GetID() != "existing" || children[1].GetID() != "child1" || children[2].GetID() != "child2" {
			t.Error("Children not added correctly to existing children")
		}
	})

	t.Run("Add nil and empty slice", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		parent.AddChildren(nil) // Should not panic

		// Initial add
		child := NewAtom("item", WithID("child"))
		parent.AddChildren([]AtomInterface{child})

		// Add empty slice
		parent.AddChildren([]AtomInterface{})

		children := parent.GetChildren()
		if len(children) != 1 || children[0].GetID() != "child" {
			t.Error("Adding nil or empty slice should not affect existing children")
		}
	})

	t.Run("Add with nil children", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		child := NewAtom("item", WithID("child"))

		parent.AddChildren([]AtomInterface{child, nil, nil})

		children := parent.GetChildren()
		if len(children) != 1 || children[0].GetID() != "child" {
			t.Error("Nil children should be skipped")
		}
	})
}

func TestAtom_SetChildren(t *testing.T) {
	t.Run("Set children on empty parent", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		child1 := NewAtom("item", WithID("child1"))
		child2 := NewAtom("item", WithID("child2"))
		parent.SetChildren([]AtomInterface{child1, child2})
		children := parent.GetChildren()

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
		parent.AddChild(oldChild)

		newChild1 := NewAtom("item", WithID("new1"))
		newChild2 := NewAtom("item", WithID("new2"))
		parent.SetChildren([]AtomInterface{newChild1, newChild2})

		children := parent.GetChildren()
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
		parent.AddChild(child)

		parent.SetChildren([]AtomInterface{})
		if len(parent.GetChildren()) != 0 {
			t.Error("Expected no children after setting to empty slice")
		}

		parent.SetChildren(nil)
		if len(parent.GetChildren()) != 0 {
			t.Error("Expected no children after setting to nil")
		}
	})

	t.Run("Set with nil children", func(t *testing.T) {
		parent := NewAtom("container", WithID("parent"))
		child := NewAtom("item", WithID("child"))

		parent.SetChildren([]AtomInterface{child, nil, nil})

		children := parent.GetChildren()
		if len(children) != 1 || children[0].GetID() != "child" {
			t.Error("Nil children should be skipped")
		}
	})
}

func TestAtom_ConcurrentChildOperations(t *testing.T) {
	parent := NewAtom("container", WithID("parent"))
	var wg sync.WaitGroup

	// Test concurrent AddChild
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			child := NewAtom("item", WithID(fmt.Sprintf("child-%d", id)))
			parent.AddChild(child)
		}(i)
	}

	// Test concurrent SetChildren
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			children := []AtomInterface{
				NewAtom("item", WithID(fmt.Sprintf("new-child1-%d", id))),
				NewAtom("item", WithID(fmt.Sprintf("new-child2-%d", id))),
			}
			parent.SetChildren(children)
		}(i)
	}

	// Test concurrent GetChildren
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			_ = parent.GetChildren()
		}()
	}

	wg.Wait()

	// Verify some invariants
	children := parent.GetChildren()
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
