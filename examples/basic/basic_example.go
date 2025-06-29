// Basic example demonstrating the core functionality of the omni package.
package main

import (
	"fmt"
	"github.com/dracory/omni"
)

func main() {
	// Create a new atom with properties and children using functional options
	rectangle := omni.NewAtom("rectangle",
		omni.WithID("rect1"),
		omni.WithProperties(map[string]string{
			"width":  "100",
			"height": "50",
		}),
		omni.WithChildren(
			omni.NewAtom("shadow",
				omni.WithID("shadow1"),
				omni.WithProperties(map[string]string{
					"opacity": "0.5",
				}),
			),
		),
	)

	// Print the rectangle's properties
	fmt.Printf("\nRectangle '%s' properties:\n", rectangle.GetID())
	for key, value := range rectangle.GetAll() {
		fmt.Printf("  %s: %s\n", key, value)
	}

	// Print the rectangle's children
	fmt.Println("\nRectangle children:")
	for _, child := range rectangle.ChildrenGet() {
		fmt.Printf("  - %s (%s)\n", child.GetID(), child.GetType())
		for key, value := range child.GetAll() {
			fmt.Printf("    %s: %s\n", key, value)
		}
	}
}
