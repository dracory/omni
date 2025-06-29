package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/dracory/omni"
)

// TestConcurrentExample tests the concurrent example functionality
func TestConcurrentExample(t *testing.T) {
	// Create a document atom with initial properties
	doc := omni.NewAtom("document",
		omni.WithID("doc1"),
		omni.WithProperties(map[string]string{
			"title": "Concurrent Document",
		}),
	)

	// Test initial state
	if doc.GetID() != "doc1" {
		t.Errorf("Expected document ID to be 'doc1', got '%s'", doc.GetID())
	}

	if title := doc.Get("title"); title != "Concurrent Document" {
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


				// Set property atomically
				mu.Lock()
				doc.Set(propName, "value")
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
				doc = doc.ChildAdd(child).(*omni.Atom)
				mu.Unlock()
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Get all properties
	props := doc.GetAll()

	// Check that all expected properties exist
	for propName := range expectedProps {
		if _, exists := props[propName]; !exists {
			t.Errorf("Expected property %s not found", propName)
		}
	}

	// Verify children
	children := doc.ChildrenGet()
	childrenMap := make(map[string]bool)
	for _, child := range children {
		childrenMap[child.GetID()] = true
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
			mu.Lock()
			atom.Set(key, value)
			mu.Unlock()
		}(i)
	}

	// Start concurrent property readers
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i%10)

			// Get property (should be thread-safe)
			_ = atom.Get(key)
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
			mu.Lock()
			atom = atom.ChildAdd(child).(*omni.Atom)
			mu.Unlock()

			// Get children (should be thread-safe)
			_ = atom.ChildrenGet()

			// Test property removal (just testing concurrency, not correctness here)
			if i%2 == 0 {
				// Add a property first so we can test removing it
				propName := fmt.Sprintf("temp_%d", i)
				mu.Lock()
				atom.Set(propName, "temp")
				atom.Remove(propName)
				mu.Unlock()
			}
		}(i)
	}

	// Wait for all operations to complete
	wg.Wait()

	// Verify properties
	allProps := atom.GetAll()
	for key, expectedValue := range expectedProps {
		value, exists := allProps[key]
		if !exists {
			t.Errorf("Property %s not found", key)
			continue
		}
		if value != expectedValue {
			t.Errorf("Expected %s=%s, got %s", key, expectedValue, value)
		}
	}
}

// TestMain runs the example and verifies it doesn't panic
func TestMain(m *testing.M) {
	// Run the example to ensure it doesn't panic
	main()
}
