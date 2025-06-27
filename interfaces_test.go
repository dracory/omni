package omni_test

import (
	"testing"

	"github.com/dracory/omni"
)

// mockProperty is a test implementation of PropertyInterface
type mockProperty struct {
	name  string
	value string
}

func (p *mockProperty) GetName() string {
	return p.name
}

func (p *mockProperty) SetName(name string) {
	p.name = name
}

func (p *mockProperty) GetValue() string {
	return p.value
}

func (p *mockProperty) SetValue(value string) {
	p.value = value
}

// mockAtom is a test implementation of AtomInterface
type mockAtom struct {
	id         string
	atomType   string
	properties []omni.PropertyInterface
	children   []omni.AtomInterface
}

func (a *mockAtom) GetID() string {
	return a.id
}

func (a *mockAtom) SetID(id string) {
	a.id = id
}

func (a *mockAtom) GetType() string {
	return a.atomType
}

func (a *mockAtom) SetType(atomType string) {
	a.atomType = atomType
}

func (a *mockAtom) GetProperties() []omni.PropertyInterface {
	return a.properties
}

func (a *mockAtom) SetProperties(properties []omni.PropertyInterface) {
	a.properties = properties
}

func (a *mockAtom) GetProperty(name string) omni.PropertyInterface {
	for _, p := range a.properties {
		if p.GetName() == name {
			return p
		}
	}
	return nil
}

func (a *mockAtom) SetProperty(property omni.PropertyInterface) {
	for i, p := range a.properties {
		if p.GetName() == property.GetName() {
			a.properties[i] = property
			return
		}
	}
	a.properties = append(a.properties, property)
}

func (a *mockAtom) RemoveProperty(name string) {
	for i, p := range a.properties {
		if p.GetName() == name {
			a.properties = append(a.properties[:i], a.properties[i+1:]...)
			return
		}
	}
}

func (a *mockAtom) SetChildren(children []omni.AtomInterface) {
	a.children = children
}

func (a *mockAtom) GetChildren() []omni.AtomInterface {
	return a.children
}

func (a *mockAtom) AddChild(child omni.AtomInterface) {
	a.children = append(a.children, child)
}

func (a *mockAtom) AddChildren(children []omni.AtomInterface) {
	for _, child := range children {
		if child == nil {
			continue
		}
		a.children = append(a.children, child)
	}
}

func TestPropertyInterface(t *testing.T) {
	t.Run("Property Get/Set Name", func(t *testing.T) {
		p := &mockProperty{}
		p.SetName("test")
		if got := p.GetName(); got != "test" {
			t.Errorf("GetName() = %v, want %v", got, "test")
		}
	})

	t.Run("Property Get/Set Value", func(t *testing.T) {
		p := &mockProperty{}
		p.SetValue("value")
		if got := p.GetValue(); got != "value" {
			t.Errorf("GetValue() = %v, want %v", got, "value")
		}
	})
}

func TestAtomInterface(t *testing.T) {
	t.Run("Atom ID Operations", func(t *testing.T) {
		a := &mockAtom{}
		a.SetID("123")
		if got := a.GetID(); got != "123" {
			t.Errorf("GetID() = %v, want %v", got, "123")
		}
	})

	t.Run("Atom Type Operations", func(t *testing.T) {
		a := &mockAtom{}
		a.SetType("test-type")
		if got := a.GetType(); got != "test-type" {
			t.Errorf("GetType() = %v, want %v", got, "test-type")
		}
	})

	t.Run("Atom Property Operations", func(t *testing.T) {
		a := &mockAtom{}
		prop1 := &mockProperty{name: "prop1", value: "value1"}
		prop2 := &mockProperty{name: "prop2", value: "value2"}

		// Test setting and getting single property
		a.SetProperty(prop1)
		if got := a.GetProperty("prop1"); got != prop1 {
			t.Errorf("GetProperty() = %v, want %v", got, prop1)
		}

		// Test getting non-existent property
		nilProp := a.GetProperty("nonexistent")
		if nilProp != nil {
			t.Errorf("Expected nil for non-existent property, got %v", nilProp)
		}

		// Test setting multiple properties
		a.SetProperties([]omni.PropertyInterface{prop1, prop2})
		if got := len(a.GetProperties()); got != 2 {
			t.Errorf("Expected 2 properties, got %d", got)
		}

		// Test updating existing property
		prop1Updated := &mockProperty{name: "prop1", value: "updated"}
		a.SetProperty(prop1Updated)
		if got := a.GetProperty("prop1").GetValue(); got != "updated" {
			t.Errorf("GetProperty().GetValue() = %v, want %v", got, "updated")
		}

		// Test removing property
		a.RemoveProperty("prop1")
		if got := a.GetProperty("prop1"); got != nil {
			t.Errorf("Expected property to be removed, but got %v", got)
		}
		if got := len(a.GetProperties()); got != 1 {
			t.Errorf("Expected 1 property after removal, got %d", got)
		}
	})

	t.Run("Atom Child Operations", func(t *testing.T) {
		parent := &mockAtom{}
		child1 := &mockAtom{id: "child1"}
		child2 := &mockAtom{id: "child2"}

		// Test adding children
		parent.AddChild(child1)
		parent.AddChild(child2)

		// Test getting children
		children := parent.GetChildren()
		if got := len(children); got != 2 {
			t.Fatalf("Expected 2 children, got %d", got)
		}
		if got := children[0].GetID(); got != "child1" {
			t.Errorf("First child ID = %v, want %v", got, "child1")
		}
		if got := children[1].GetID(); got != "child2" {
			t.Errorf("Second child ID = %v, want %v", got, "child2")
		}
	})
}
