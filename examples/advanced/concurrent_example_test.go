package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/dracory/omni"
)

// TestConcurrentExample tests the concurrent example functionality
func TestConcurrentExample(t *testing.T) {
	// Create a document atom with initial properties
	doc := omni.NewAtom("document",
		omni.WithID("doc1"),
		omni.WithProperties(
			omni.NewProperty("title", "Concurrent Document"),
		),
	)

	// Test initial state
	if doc.GetID() != "doc1" {
		t.Errorf("Expected document ID to be 'doc1', got '%s'", doc.GetID())
	}

	titleProp := doc.GetProperty("title")
	if titleProp == nil || titleProp.GetValue() != "Concurrent Document" {
		t.Error("Expected title property to be 'Concurrent Document'")
	}

	// Test concurrent updates
	const numWorkers = 5
	const updatesPerWorker = 3 // Reduced for test speed

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Track expected properties and children
	expectedProps := make(map[string]bool)
	expectedChildren := make(map[string]bool)

	// Start concurrent workers to update the document
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Add properties and create children
			for j := 0; j < updatesPerWorker; j++ {
				// Each worker adds a new property with a unique name
				propName := fmt.Sprintf("worker%d_update%d", workerID, j)

				// Track expected properties
				mu.Lock()
				expectedProps[propName] = true
				mu.Unlock()

				// Create and add property atomically
				prop := omni.NewProperty(propName, "value")
				mu.Lock()
				doc.SetProperty(prop)
				mu.Unlock()

				// Create child atom with type based on iteration
				childType := fmt.Sprintf("type_%d", j%3)
				childID := fmt.Sprintf("worker%d_%d", workerID, j)
				child := omni.NewAtom(childType, omni.WithID(childID))

				// Track expected children
				mu.Lock()
				expectedChildren[childID] = true
				mu.Unlock()

				// Add child atomically
				mu.Lock()
				doc.AddChild(child)
				mu.Unlock()
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Verify properties
	props := doc.GetProperties()
	propsMap := make(map[string]bool)
	for _, prop := range props {
		propsMap[prop.GetName()] = true
	}

	// Check that all expected properties exist
	for propName := range expectedProps {
		if !propsMap[propName] {
			t.Errorf("Expected property %s not found", propName)
		}
	}

	// Verify children
	children := doc.GetChildren()
	if len(children) != numWorkers*updatesPerWorker {
		t.Errorf("Expected %d children, got %d", numWorkers*updatesPerWorker, len(children))
	}

	// Check that all expected children exist
	childrenMap := make(map[string]bool)
	for _, child := range children {
		childrenMap[child.GetID()] = true
		// Verify child type is one of the expected types
		childType := child.GetType()
		if !strings.HasPrefix(childType, "type_") ||
			childType < "type_0" || childType > "type_2" {
			t.Errorf("Unexpected child type: %s", childType)
		}
	}

	for childID := range expectedChildren {
		if !childrenMap[childID] {
			t.Errorf("Expected child %s not found", childID)
		}
	}
}

// TestConcurrentSafety tests thread safety of concurrent operations
func TestConcurrentSafety(t *testing.T) {
	// Create a new atom
	atom := omni.NewAtom("test", omni.WithID("test1"))

	// Number of concurrent operations
	const numOps = 100

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Track expected values
	expectedProps := make(map[string]string)

	// Start concurrent property writers
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i%10) // 10 unique keys
			value := fmt.Sprintf("value%d", i)

			// Update expected value
			mu.Lock()
			expectedProps[key] = value
			mu.Unlock()

			// Set property
			prop := omni.NewProperty(key, value)
			atom.SetProperty(prop)
		}(i)
	}

	// Start concurrent property readers
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i%10)

			// Get property (should be thread-safe)
			prop := atom.GetProperty(key)
			if prop != nil {
				// Just reading the value should be safe
				_ = prop.GetValue()
			}
		}(i)
	}

	// Start concurrent children operations
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			childID := fmt.Sprintf("child%d", i)
			child := omni.NewAtom("child", omni.WithID(childID))

			// Add child
			atom.AddChild(child)

			// Get children (should be thread-safe)
			_ = atom.GetChildren()

			// Remove property (just testing concurrency, not correctness here)
			if i%2 == 0 {
				// Add a property first so we can test removing it
				propName := fmt.Sprintf("temp_%d", i)
				atom.SetProperty(omni.NewProperty(propName, "temp"))
				// Then remove it
				atom.RemoveProperty(propName)
			}
		}(i)
	}

	// Wait for all operations to complete
	wg.Wait()

	// Verify properties
	for key, expectedValue := range expectedProps {
		prop := atom.GetProperty(key)
		if prop == nil {
			t.Errorf("Property %s not found", key)
			continue
		}
		if prop.GetValue() != expectedValue {
			t.Errorf("Expected %s=%s, got %s", key, expectedValue, prop.GetValue())
		}
	}
}

// TestMain runs the example and verifies it doesn't panic
func TestMain(m *testing.M) {
	// Run the example to ensure it doesn't panic
	main()
}
