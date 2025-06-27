package omni

import (
	"sync"
	"testing"
)

func TestNewAtom(t *testing.T) {
	a := NewAtom("test-id", "test-type")
	if a == nil {
		t.Fatal("NewAtom() returned nil")
	}
	if got := a.GetID(); got != "test-id" {
		t.Errorf("GetID() = %v, want %v", got, "test-id")
	}
	if got := a.GetType(); got != "test-type" {
		t.Errorf("GetType() = %v, want %v", got, "test-type")
	}
}

func TestAtom_PropertyManagement(t *testing.T) {
	a := NewAtom("1", "test")
	prop1 := NewProperty("prop1", "value1")
	prop2 := NewProperty("prop2", "value2")

	// Test SetProperty and GetProperty
	a.SetProperty(prop1)
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
	a.SetProperties([]PropertyInterface{prop1, prop2})
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
	parent := NewAtom("parent", "container")
	child1 := NewAtom("child1", "item")
	child2 := NewAtom("child2", "item")

	// Test AddChild and GetChildren
	parent.AddChild(child1)
	parent.AddChild(child2)

	children := parent.GetChildren()
	if len(children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(children))
	}
	if children[0].GetID() != "child1" || children[1].GetID() != "child2" {
		t.Error("Children not added in correct order")
	}
}

func TestAtom_ConcurrentAccess(t *testing.T) {
	a := NewAtom("concurrent", "test")
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
				child := NewAtom("child", "item")
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
	a := NewAtom("edge", "test")

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
