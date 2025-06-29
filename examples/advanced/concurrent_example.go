// Advanced example demonstrating concurrent access and property management
package main

import (
	"fmt"
	"sync"

	"github.com/dracory/omni"
)

func main() {
	// Create a document atom with initial properties
	doc := omni.NewAtom("document",
		omni.WithID("doc1"),
		omni.WithProperties(map[string]string{
			"title": "Concurrent Document",
		}),
	)

	// Number of concurrent workers
	const numWorkers = 5
	const updatesPerWorker = 10

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Start concurrent workers to update the document
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Local slice to collect this worker's children
			var children []omni.AtomInterface

			// Add properties and create children
			for j := 0; j < updatesPerWorker; j++ {
				// Each worker adds a new property with a unique name
				propName := fmt.Sprintf("worker%d_update%d", workerID, j)
				
				// Set property atomically
				mu.Lock()
				doc.Set(propName, "value")
				mu.Unlock()

				// Create child atom with type based on iteration
				childType := fmt.Sprintf("type_%d", j%3)
				child := omni.NewAtom(childType,
					omni.WithID(fmt.Sprintf("child_w%d_u%d", workerID, j)),
				)
				children = append(children, child)
			}

			// Add all children atomically one by one
			if len(children) > 0 {
				mu.Lock()
				for _, child := range children {
					doc = doc.ChildAdd(child).(*omni.Atom)
				}
				mu.Unlock()
			}
		}(i)
	}

	// Wait for all workers to finish
	wg.Wait()

	// Get all properties and children
	allProps := doc.GetAll()
	children := doc.ChildrenGet()

	// Print the results
	fmt.Printf("Document '%s' has %d properties and %d children\n",
		doc.GetID(),
		len(allProps),
		len(children),
	)

	// Print a sample of properties
	fmt.Println("\nSample of properties:")
	i := 0
	for key, value := range allProps {
		if i >= 5 { // Limit to first 5 properties
			break
		}
		fmt.Printf("  %s = %s\n", key, value)
		i++
	}
	if len(allProps) > 5 {
		fmt.Printf("  ... and %d more\n", len(allProps)-5)
	}

	// Print a sample of children
	fmt.Printf("\nSample of children (%d total):\n", len(children))
	for i, child := range children {
		if i >= 3 { // Limit to first 3 children
			break
		}
		fmt.Printf("  - %s (%s)\n", child.GetID(), child.GetType())
	}
	if len(children) > 3 {
		fmt.Printf("  ... and %d more\n", len(children)-3)
	}

	// Count children by type
	typeCount := make(map[string]int)
	for _, child := range children {
		typeCount[child.GetType()]++
	}
	fmt.Println("\nChildren by type:")
	for typ, count := range typeCount {
		fmt.Printf("  %s: %d\n", typ, count)
	}
}
