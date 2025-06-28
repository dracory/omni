package omni

import (
	"sync"
	"testing"
)

// Test helper functions

// Helper function to create a test atom with properties and children
func createTestAtom(id, atomType string, props map[string]string, children []*Atom) *Atom {
	opts := make([]AtomOption, 0, 2)
	if id != "" {
		opts = append(opts, WithID(id))
	}

	// Add properties if any
	if len(props) > 0 {
		propsList := make([]PropertyInterface, 0, len(props))
		for k, v := range props {
			propsList = append(propsList, NewProperty(k, v))
		}
		opts = append(opts, WithProperties(propsList...))
	}

	// Add children if any
	if len(children) > 0 {
		childrenList := make([]AtomInterface, len(children))
		for i, child := range children {
			childrenList[i] = child
		}
		opts = append(opts, WithChildren(childrenList...))
	}

	return NewAtom(atomType, opts...)
}

// Helper function to create a test property
func createTestProperty(name, value string) *Property {
	return NewProperty(name, value)
}

func TestPropertyInterface(t *testing.T) {
	// Test creating a new property
	t.Run("Create New Property", func(t *testing.T) {
		p := NewProperty("test", "value")
		if p == nil {
			t.Error("Expected new property to be created")
		}
	})
	t.Run("Property Get/Set Name", func(t *testing.T) {
		p := NewProperty("", "")
		p.SetName("test")
		if p.GetName() != "test" {
			t.Error("Expected name to be 'test'")
		}
	})

	t.Run("Property Get/Set Value", func(t *testing.T) {
		p := NewProperty("", "")
		p.SetValue("value")
		if p.GetValue() != "value" {
			t.Error("Expected value to be 'value'")
		}
	})
}

func TestAtomInterface(t *testing.T) {
	// Test creating a new atom
	t.Run("Create New Atom", func(t *testing.T) {
		a := NewAtom("test-type")
		if a == nil {
			t.Error("Expected new atom to be created")
		}
	})
	t.Run("Atom Get/Set ID", func(t *testing.T) {
		a := NewAtom("test-type", WithID("test-id"))
		if a.GetID() != "test-id" {
			t.Error("Expected ID to be 'test-id'")
		}
	})

	t.Run("Atom Get/Set Type", func(t *testing.T) {
		a := NewAtom("test-type")
		if a.GetType() != "test-type" {
			t.Error("Expected type to be 'test-type'")
		}
	})

	t.Run("Atom Add/Get Properties", func(t *testing.T) {
		a := NewAtom("test-type")
		prop1 := NewProperty("prop1", "value1")
		prop2 := NewProperty("prop2", "value2")

		a.SetProperty(prop1)
		a.SetProperty(prop2)

		foundProp1 := a.GetProperty("prop1")
		if foundProp1 == nil || foundProp1.GetValue() != "value1" {
			t.Error("Failed to get property by name")
		}

		foundProp2 := a.GetProperty("prop2")
		if foundProp2 == nil || foundProp2.GetValue() != "value2" {
			t.Error("Failed to get property by name")
		}

		notFoundProp := a.GetProperty("nonexistent")
		if notFoundProp != nil {
			t.Error("Expected nil when getting nonexistent property")
		}
	})

	t.Run("Atom Add/Get Children", func(t *testing.T) {
		parent := NewAtom("parent")
		child1 := NewAtom("child", WithID("child1"))
		child2 := NewAtom("child", WithID("child2"))

		// Test AddChild and GetChildren
		parent.AddChild(child1)
		parent.AddChild(child2)
		children := parent.GetChildren()
		if len(children) != 2 {
			t.Errorf("Expected 2 children, got %d", len(children))
		}
		if children[0].GetID() != "child1" || children[1].GetID() != "child2" {
			t.Errorf("Children IDs do not match expected values. Got: %s, %s", children[0].GetID(), children[1].GetID())
		}
	})

	t.Run("Atom ToMap", func(t *testing.T) {
		a := NewAtom("test-type", WithID("test-id"))
		child := NewAtom("child-type", WithID("child-1"))

		a.AddChild(child)

		prop := NewProperty("test-prop", "test-value")
		a.SetProperty(prop)

		m := a.ToMap()

		if m["id"] != "test-id" || m["type"] != "test-type" {
			t.Error("Map values do not match expected")
		}

		children, ok := m["children"].([]map[string]interface{})
		if !ok || len(children) != 1 || children[0]["id"] != "child-1" {
			t.Errorf("Children in map do not match expected. Got: %+v", children)
		}

		props, ok := m["parameters"].(map[string]string)
		if !ok || props["test-prop"] != "test-value" {
			t.Error("Properties in map do not match expected")
		}
	})

	t.Run("Atom Child Operations", func(t *testing.T) {
		parent := NewAtom("parent-type")
		child1 := NewAtom("child-type")
		child2 := NewAtom("child-type")

		// Test adding single child
		parent.AddChild(child1)
		parent.AddChild(child2)

		// Test getting children
		children := parent.GetChildren()
		if len(children) != 2 {
			t.Fatalf("Expected 2 children, got %d", len(children))
		}

		// Test AddChildren with multiple children
		child3 := NewAtom("child-type")
		child4 := NewAtom("child-type")
		parent.AddChildren([]AtomInterface{child3, child4})
		if len(parent.GetChildren()) != 4 {
			t.Fatalf("Expected 4 children after AddChildren, got %d", len(parent.GetChildren()))
		}

		// Test AddChildren with nil children
		initialCount := len(parent.GetChildren())
		parent.AddChildren([]AtomInterface{nil, nil})
		if len(parent.GetChildren()) != initialCount {
			t.Fatalf("Expected %d children after adding nil children, got %d", initialCount, len(parent.GetChildren()))
		}

		// Test AddChildren with empty slice
		parent.AddChildren([]AtomInterface{})
		if len(parent.GetChildren()) != initialCount {
			t.Fatalf("Expected %d children after adding empty slice, got %d", initialCount, len(parent.GetChildren()))
		}

		// Test SetChildren
		newChildren := []AtomInterface{
			NewAtom("new-type", WithID("new-child-1")),
			NewAtom("new-type", WithID("new-child-2")),
		}
		parent.SetChildren(newChildren)
		gotChildren := parent.GetChildren()
		if len(gotChildren) != 2 {
			t.Fatalf("Expected 2 children after SetChildren, got %d", len(gotChildren))
		}
		if gotChildren[0].GetID() != "new-child-1" || gotChildren[1].GetID() != "new-child-2" {
			t.Fatalf("Children were not set correctly by SetChildren. Got IDs: %s, %s", 
				gotChildren[0].GetID(), gotChildren[1].GetID())
		}

		// Test SetChildren with empty slice
		parent.SetChildren([]AtomInterface{})
		if len(parent.GetChildren()) != 0 {
			t.Fatalf("Expected 0 children after SetChildren with empty slice, got %d", len(parent.GetChildren()))
		}

		// Test SetChildren with nil
		parent.SetChildren(nil)
		if len(parent.GetChildren()) != 0 {
			t.Fatalf("Expected 0 children after SetChildren with nil, got %d", len(parent.GetChildren()))
		}

		// Test concurrent access to children methods
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				child := NewAtom("concurrent-type")
				parent.AddChild(child)

				// Also test AddChildren in goroutines
				if id%2 == 0 {
					parent.AddChildren([]AtomInterface{
						NewAtom("batch-type"),
						NewAtom("batch-type"),
					})
				}

				// Exercise GetChildren in concurrent context
				_ = parent.GetChildren()

				// Exercise SetChildren in concurrent context
				if id%10 == 0 {
					parent.SetChildren([]AtomInterface{
						NewAtom("reset-type"),
					})
				}
			}(i)
		}
		wg.Wait()

		// Verify no data races or panics occurred during concurrent access
		// We can't make strong assertions about the final state due to concurrency,
		// but we can verify the atom is still in a valid state
		finalChildren := parent.GetChildren()
		for _, child := range finalChildren {
			if child == nil {
				t.Error("Unexpected nil child in final children")
			}
		}

		// Test that we can create a new atom with no children
		t.Run("Empty Atom", func(t *testing.T) {
			a := NewAtom("test-type")
			if len(a.GetChildren()) != 0 {
				t.Errorf("Expected new atom to have 0 children, got %d", len(a.GetChildren()))
			}
		})

		// Test that we can set and get children
		t.Run("Set and Get Children", func(t *testing.T) {
			a := NewAtom("parent-type")
			children := []AtomInterface{
				NewAtom("child-type"),
				NewAtom("child-type"),
			}
			a.SetChildren(children)
			if len(a.GetChildren()) != 2 {
				t.Errorf("Expected 2 children, got %d", len(a.GetChildren()))
			}
		})

		// Test that we can add multiple children at once
		t.Run("Add Multiple Children", func(t *testing.T) {
			a := NewAtom("parent-type")
			child1 := NewAtom("child-type")
			child2 := NewAtom("child-type")
			a.AddChildren([]AtomInterface{child1, child2})
			if len(a.GetChildren()) != 2 {
				t.Errorf("Expected 2 children after AddChildren, got %d", len(a.GetChildren()))
			}
		})

		// Test that we can add children to existing children
		t.Run("Add to Existing Children", func(t *testing.T) {
			a := NewAtom("parent-type")
			a.AddChild(NewAtom("child-type"))
			child2 := NewAtom("child-type")
			child3 := NewAtom("child-type")
			a.AddChildren([]AtomInterface{child2, child3})
			if len(a.GetChildren()) != 3 {
				t.Errorf("Expected 3 children after adding to existing, got %d", len(a.GetChildren()))
			}
		})

		// Test that we can clear children by setting to nil or empty slice
		t.Run("Clear Children", func(t *testing.T) {
			a := NewAtom("parent-type")
			child1 := NewAtom("child-type")
			child2 := NewAtom("child-type")
			a.AddChildren([]AtomInterface{child1, child2})

			a.SetChildren(nil)
			if len(a.GetChildren()) != 0 {
				t.Error("Expected 0 children after setting to nil")
			}

			a.AddChildren([]AtomInterface{
				NewAtom("child-type"),
				NewAtom("child-type"),
			})

			a.SetChildren([]AtomInterface{})
			if len(a.GetChildren()) != 0 {
				t.Error("Expected 0 children after setting to empty slice")
			}
		})
	})
}
