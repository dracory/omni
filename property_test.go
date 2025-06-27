package omni

import (
	"sync"
	"testing"
)

func TestNewProperty(t *testing.T) {
	p := NewProperty("test", "value")
	if p == nil {
		t.Fatal("NewProperty() returned nil")
	}
	if got := p.GetName(); got != "test" {
		t.Errorf("GetName() = %v, want %v", got, "test")
	}
	if got := p.GetValue(); got != "value" {
		t.Errorf("GetValue() = %v, want %v", got, "value")
	}
}

func TestProperty_SettersAndGetters(t *testing.T) {
	p := NewProperty("initial", "initial")

	// Test SetName and GetName
	p.SetName("newName")
	if got := p.GetName(); got != "newName" {
		t.Errorf("GetName() = %v, want %v", got, "newName")
	}

	// Test SetValue and GetValue
	p.SetValue("newValue")
	if got := p.GetValue(); got != "newValue" {
		t.Errorf("GetValue() = %v, want %v", got, "newValue")
	}

	// Test empty values
	p.SetName("")
	if got := p.GetName(); got != "" {
		t.Error("Failed to set empty name")
	}

	p.SetValue("")
	if got := p.GetValue(); got != "" {
		t.Error("Failed to set empty value")
	}
}

func TestProperty_ConcurrentAccess(t *testing.T) {
	p := NewProperty("test", "initial")
	var wg sync.WaitGroup

	// Start multiple goroutines to modify the property
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				// Test name operations
				name := string(rune('a' + (j % 26)))
				p.SetName(name)
				_ = p.GetName()

				// Test value operations
				value := string(rune('A' + (j % 26)))
				p.SetValue(value)
				_ = p.GetValue()
			}
		}(i)
	}

	wg.Wait()

	// Verify final state is consistent
	name := p.GetName()
	value := p.GetValue()
	if name == "" || value == "" {
		t.Error("Name or value should not be empty after concurrent operations")
	}
}
