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
		omni.WithProperties(
			omni.NewProperty("title", "Concurrent Document"),
		),
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
				
				// Create and add property atomically
				prop := omni.NewProperty(propName, "value")
				mu.Lock()
				doc.SetProperty(prop)
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
					doc.AddChild(child)
				}
				mu.Unlock()
			}
		}(i)
	}

	// Wait for all workers to finish
	wg.Wait()

	// Print the results
	fmt.Printf("Document '%s' has %d properties and %d children\n",
		doc.GetID(),
		len(doc.GetProperties()),
		len(doc.GetChildren()),
	)

	// Print a sample of properties
	fmt.Println("\nSample of properties:")
	props := doc.GetProperties()
	for i, prop := range props {
		if i >= 5 { // Limit to first 5 properties
			break
		}
		fmt.Printf("  %s = %s\n", prop.GetName(), prop.GetValue())
	}
	if len(props) > 5 {
		fmt.Printf("  ... and %d more\n", len(props)-5)
	}

	// Print a sample of children
	children := doc.GetChildren()
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
