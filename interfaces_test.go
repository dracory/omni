package omni

import (
	"sync"
	"testing"
)

func TestCreateNewAtom(t *testing.T) {
	a := NewAtom("test-type")
	if a == nil {
		t.Error("Expected new atom to be created")
	}
}

func TestAtomGetSetID(t *testing.T) {
	a := NewAtom("test-type", WithID("test-id"))
	if a.GetID() != "test-id" {
		t.Error("Expected ID to be 'test-id'")
	}
}

func TestAtomGetSetType(t *testing.T) {
	a := NewAtom("test-type")
	if a.GetType() != "test-type" {
		t.Error("Expected type to be 'test-type'")
	}
}

func TestAtomAddGetProperties(t *testing.T) {
	a := NewAtom("test-type")
	a.Set("prop1", "value1")
	a.Set("prop2", "value2")

	foundProp1 := a.Get("prop1")
	if foundProp1 != "value1" {
		t.Error("Failed to get property by name")
	}

	foundProp2 := a.Get("prop2")
	if foundProp2 != "value2" {
		t.Error("Failed to get property by name")
	}

	notFoundProp := a.Get("nonexistent")
	if notFoundProp != "" {
		t.Error("Expected empty string when getting nonexistent property")
	}
}

func TestAtomRemoveProperty(t *testing.T) {
	a := NewAtom("test-type")
	a.Set("test-prop", "test-value")

	// Test removing existing property
	a.Remove("test-prop")
	if a.Get("test-prop") != "" {
		t.Error("Expected property to be removed")
	}

	// Test removing non-existent property (should not panic)
	a.Remove("non-existent")
}

func TestAtomAddGetChildren(t *testing.T) {
	parent := NewAtom("parent")
	child1 := NewAtom("child", WithID("child1"))
	child2 := NewAtom("child", WithID("child2"))

	// Test ChildAdd and ChildrenGet
	parent.ChildAdd(child1)
	parent.ChildAdd(child2)
	children := parent.ChildrenGet()
	if len(children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(children))
	}
	if children[0].GetID() != "child1" || children[1].GetID() != "child2" {
		t.Fatalf("Children IDs do not match expected values. Got: %s, %s", children[0].GetID(), children[1].GetID())
	}
}

func TestAtomToMap(t *testing.T) {
	a := NewAtom("test-type", WithID("test-id"))
	child := NewAtom("child-type", WithID("child-1"))
	a.ChildAdd(child)
	a.Set("test-prop", "test-value")

	m := a.ToMap()

	if m["id"] != "test-id" || m["type"] != "test-type" {
		t.Error("Map values do not match expected")
	}

	children, ok := m["children"].([]map[string]interface{})
	if !ok || len(children) != 1 || children[0]["id"] != "child-1" {
		t.Errorf("Children in map do not match expected. Got: %+v", children)
	}

	props, ok := m["properties"].(map[string]string)
	if !ok || props["test-prop"] != "test-value" {
		t.Error("Properties in map do not match expected")
	}
}

func TestAtomChildOperations(t *testing.T) {
	// Test adding and getting children
	parent := NewAtom("parent-type")
	child1 := NewAtom("child-type", WithID("child1"))
	child2 := NewAtom("child-type", WithID("child2"))

	// Test adding single child
	parent.ChildAdd(child1)
	parent.ChildAdd(child2)

	// Test getting children
	children := parent.ChildrenGet()
	if len(children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(children))
	}

	// Test ChildrenAdd with multiple children
	child3 := NewAtom("child-type", WithID("child3"))
	child4 := NewAtom("child-type", WithID("child4"))
	parent.ChildrenAdd([]AtomInterface{child3, child4})
	if parent.ChildrenLength() != 4 {
		t.Fatalf("Expected 4 children after ChildrenAdd, got %d", parent.ChildrenLength())
	}

	// Test ChildrenSet
	parent.ChildrenSet([]AtomInterface{child1, child2})
	if parent.ChildrenLength() != 2 {
		t.Fatalf("Expected 2 children after ChildrenSet, got %d", parent.ChildrenLength())
	}

	// Test ChildrenAdd with empty slice
	initialCount := parent.ChildrenLength()
	parent.ChildrenAdd([]AtomInterface{})
	if parent.ChildrenLength() != initialCount {
		t.Fatalf("Expected %d children after adding empty slice, got %d", initialCount, parent.ChildrenLength())
	}
}

func TestAtomConcurrentOperations(t *testing.T) {
	parent := NewAtom("parent-type")
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			child := NewAtom("concurrent-type")
			parent.ChildAdd(child)

			// Also test AddChildren in goroutines
			if id%2 == 0 {
				parent.ChildrenAdd([]AtomInterface{
					NewAtom("batch-type"),
					NewAtom("batch-type"),
				})
			}

			// Exercise GetChildren in concurrent context
			_ = parent.ChildrenGet()

			// Exercise SetChildren in concurrent context
			if id%10 == 0 {
				parent.ChildrenSet([]AtomInterface{
					NewAtom("reset-type"),
				})
			}
		}(i)
	}
	wg.Wait()

	// Verify no data races or panics occurred during concurrent access
	finalChildren := parent.ChildrenGet()
	for _, child := range finalChildren {
		if child == nil {
			t.Error("Unexpected nil child in final children")
		}
	}
}

func TestEmptyAtom(t *testing.T) {
	a := NewAtom("test-type")
	if len(a.ChildrenGet()) != 0 {
		t.Errorf("Expected new atom to have 0 children, got %d", len(a.ChildrenGet()))
	}
}

func TestSetAndGetChildren(t *testing.T) {
	a := NewAtom("parent-type")
	children := []AtomInterface{
		NewAtom("child-type"),
		NewAtom("child-type"),
	}
	a.ChildrenSet(children)
	if len(a.ChildrenGet()) != 2 {
		t.Errorf("Expected 2 children, got %d", len(a.ChildrenGet()))
	}
}

func TestAddMultipleChildren(t *testing.T) {
	a := NewAtom("parent-type")
	child1 := NewAtom("child-type")
	child2 := NewAtom("child-type")
	a.ChildrenAdd([]AtomInterface{child1, child2})
	if len(a.ChildrenGet()) != 2 {
		t.Errorf("Expected 2 children after AddChildren, got %d", len(a.ChildrenGet()))
	}
}

func TestAddToExistingChildren(t *testing.T) {
	a := NewAtom("parent-type")
	a.ChildAdd(NewAtom("child-type"))
	child2 := NewAtom("child-type")
	child3 := NewAtom("child-type")
	a.ChildrenAdd([]AtomInterface{child2, child3})
	if len(a.ChildrenGet()) != 3 {
		t.Errorf("Expected 3 children after adding to existing, got %d", len(a.ChildrenGet()))
	}
}

func TestClearChildren(t *testing.T) {
	a := NewAtom("parent-type")
	child1 := NewAtom("child-type")
	child2 := NewAtom("child-type")
	a.ChildrenAdd([]AtomInterface{child1, child2})

	a.ChildrenSet(nil)
	if len(a.ChildrenGet()) != 0 {
		t.Error("Expected 0 children after setting to nil")
	}

	a.ChildrenAdd([]AtomInterface{
		NewAtom("child-type"),
		NewAtom("child-type"),
	})

	a.ChildrenSet([]AtomInterface{})
	if len(a.ChildrenGet()) != 0 {
		t.Error("Expected 0 children after setting to empty slice")
	}
}

func TestSetChildrenWithNewChildren(t *testing.T) {
	parent := NewAtom("parent-type")
	newChildren := []AtomInterface{
		NewAtom("new-type", WithID("new-child-1")),
		NewAtom("new-type", WithID("new-child-2")),
	}
	parent.ChildrenSet(newChildren)
	gotChildren := parent.ChildrenGet()
	if len(gotChildren) != 2 {
		t.Fatalf("Expected 2 children after ChildrenSet, got %d", len(gotChildren))
	}
	if gotChildren[0].GetID() != "new-child-1" || gotChildren[1].GetID() != "new-child-2" {
		t.Fatalf("Children were not set correctly by SetChildren. Got IDs: %s, %s",
			gotChildren[0].GetID(), gotChildren[1].GetID())
	}
}

func TestSetChildrenWithEmptySlice(t *testing.T) {
	parent := NewAtom("parent-type")
	parent.ChildrenSet([]AtomInterface{})
	if len(parent.ChildrenGet()) != 0 {
		t.Fatalf("Expected 0 children after SetChildren with empty slice, got %d", len(parent.ChildrenGet()))
	}
}

func TestSetChildrenWithNil(t *testing.T) {
	parent := NewAtom("parent-type")
	parent.ChildrenSet(nil)
	if len(parent.ChildrenGet()) != 0 {
		t.Fatalf("Expected 0 children after SetChildren with nil, got %d", len(parent.ChildrenGet()))
	}
}
