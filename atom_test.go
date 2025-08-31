package omni

import (
	"fmt"
	"strings"
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

func TestChildDeleteByID_RemovesCorrectChild(t *testing.T) {
	p := NewAtom("parent")
	c1 := NewAtom("child", WithID("c1"))
	c2 := NewAtom("child", WithID("c2"))
	c3 := NewAtom("child", WithID("c3"))
	p.ChildAdd(c1).ChildAdd(c2).ChildAdd(c3)

	p.ChildDeleteByID("c2")

	kids := p.ChildrenGet()
	if len(kids) != 2 {
		t.Fatalf("want 2 children after delete, got %d", len(kids))
	}
	if kids[0].GetID() != "c1" || kids[1].GetID() != "c3" {
		t.Fatalf("unexpected children order/ids: %s, %s", kids[0].GetID(), kids[1].GetID())
	}

	// deleting non-existing should be no-op
	p.ChildDeleteByID("missing")
	if p.ChildrenLength() != 2 {
		t.Fatalf("delete of missing changed length: %d", p.ChildrenLength())
	}
}

func TestChildrenGet_ReturnsCopy(t *testing.T) {
	p := NewAtom("parent")
	c1 := NewAtom("child", WithID("c1"))
	c2 := NewAtom("child", WithID("c2"))
	p.ChildrenAdd([]AtomInterface{c1, c2})

	got := p.ChildrenGet()
	if len(got) != 2 {
		t.Fatalf("want 2 children, got %d", len(got))
	}
	// mutate returned slice; underlying parent slice must not change
	got[0] = NewAtom("child", WithID("other"))
	kids := p.ChildrenGet()
	if kids[0].GetID() != "c1" {
		t.Fatalf("parent's children modified via copy; want c1, got %s", kids[0].GetID())
	}
}

func TestChildrenSet_FiltersNil(t *testing.T) {
	p := NewAtom("parent")
	c1 := NewAtom("child", WithID("c1"))
	p.ChildrenSet([]AtomInterface{c1, nil, nil})
	if p.ChildrenLength() != 1 {
		t.Fatalf("want 1 child after filtering nils, got %d", p.ChildrenLength())
	}
}

func TestSetAllAndGetAll_CopySemantics(t *testing.T) {
	p := NewAtom("parent")
	props := map[string]string{"a": "1", "b": "2"}
	p.SetAll(props)

	// mutate original map after SetAll; atom should now reflect new map reference
	props["a"] = "x"
	if v := p.Get("a"); v != "x" {
		t.Fatalf("expected atom to reflect SetAll map reference changes; got %q", v)
	}

	// GetAll should return a copy; mutating it should not affect atom
	got := p.GetAll()
	got["a"] = "y"
	if v := p.Get("a"); v != "x" {
		t.Fatalf("GetAll must return copy; mutation leaked into atom: %q", v)
	}
}

func TestHasAndRemove_Behavior(t *testing.T) {
	p := NewAtom("parent")
	if p.Has("k") {
		t.Fatal("Has should be false when properties are nil/empty")
	}
	p.Set("k", "v")
	if !p.Has("k") {
		t.Fatal("Has should be true after Set")
	}
	p.Remove("k")
	if p.Has("k") || p.Get("k") != "" {
		t.Fatal("Remove should delete key and Get should return empty string")
	}
	// removing non-existing should not panic
	p.Remove("missing")
}

func TestWithData_SetsIDTypeAndProps(t *testing.T) {
	p := NewAtom("ignored", WithData(map[string]string{
		"id":   "ID1",
		"type": "T1",
		"x":    "y",
	}))
	if p.GetID() != "ID1" || p.GetType() != "T1" || p.Get("x") != "y" {
		t.Fatalf("WithData did not set fields correctly: id=%s type=%s x=%s", p.GetID(), p.GetType(), p.Get("x"))
	}
}

func TestToJSONPretty_AndToMap(t *testing.T) {
	p := NewAtom("t", WithID("i"))
	p.Set("a", "1")
	j, err := p.ToJSONPretty()
	if err != nil {
		t.Fatalf("ToJSONPretty error: %v", err)
	}
	if !strings.Contains(j, "\n  ") { // expects indentation
		t.Fatalf("Pretty JSON does not look indented: %q", j)
	}
	m := p.ToMap()
	if m["id"].(string) != "i" || m["type"].(string) != "t" {
		t.Fatalf("ToMap missing id/type: %+v", m)
	}
	if props, ok := m["properties"].(map[string]string); !ok || props["a"] != "1" {
		t.Fatalf("ToMap properties incorrect: %#v", m["properties"])
	}
}

func TestGobEncodeDecode_Wrappers(t *testing.T) {
	p := NewAtom("t", WithID("i")).(*Atom)
	p.Set("k", "v")
	data, err := p.GobEncode()
	if err != nil {
		t.Fatalf("GobEncode error: %v", err)
	}
	var q Atom
	if err := q.GobDecode(data); err != nil {
		t.Fatalf("GobDecode error: %v", err)
	}
	if q.GetID() != "i" || q.GetType() != "t" || q.Get("k") != "v" {
		t.Fatalf("decoded atom mismatch: id=%s type=%s k=%s", q.GetID(), q.GetType(), q.Get("k"))
	}
}

func TestMemoryUsage_PositiveAndIncreases(t *testing.T) {
	p := NewAtom("t")
	base := p.MemoryUsage()
	if base <= 0 {
		t.Fatalf("MemoryUsage should be positive, got %d", base)
	}
	p.Set("a", strings.Repeat("x", 100))
	p.ChildAdd(NewAtom("child"))
	if grew := p.MemoryUsage(); grew < base {
		t.Fatalf("MemoryUsage should increase after adding data: before=%d after=%d", base, grew)
	}
}

func TestToJSON_IncludesIDTypeAndProps(t *testing.T) {
	p := NewAtom("t", WithID("i"))
	p.Set("a", "1")
	j, err := p.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}
	if !strings.Contains(j, `"id":"i"`) || !strings.Contains(j, `"type":"t"`) {
		t.Fatalf("ToJSON missing id/type: %q", j)
	}
	if !strings.Contains(j, `"properties":{"a":"1"}`) {
		t.Fatalf("ToJSON missing properties: %q", j)
	}
}

func TestToJSON_WithoutProperties_OmitsPropertiesField(t *testing.T) {
	p := NewAtom("t", WithID("i"))
	j, err := p.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}
	if !strings.Contains(j, `"id":"i"`) || !strings.Contains(j, `"type":"t"`) {
		t.Fatalf("ToJSON missing id/type: %q", j)
	}
	if strings.Contains(j, `"properties"`) {
		t.Fatalf("ToJSON should omit empty properties field: %q", j)
	}
}
