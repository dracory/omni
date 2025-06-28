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
		omni.WithProperties(
			omni.NewProperty("width", "100"),
			omni.NewProperty("height", "50"),
		),
		omni.WithChildren(
			omni.NewAtom("shadow",
				omni.WithID("shadow1"),
				omni.WithProperties(
					omni.NewProperty("opacity", "0.5"),
				),
			),
		),
	)

	// Print the rectangle's properties
	fmt.Printf("\nRectangle '%s' properties:\n", rectangle.GetID())
	for _, prop := range rectangle.GetProperties() {
		fmt.Printf("  %s: %s\n", prop.GetName(), prop.GetValue())
	}

	// Print the rectangle's children
	fmt.Println("\nRectangle children:")
	for _, child := range rectangle.GetChildren() {
		fmt.Printf("  - %s (%s)\n", child.GetID(), child.GetType())
		for _, prop := range child.GetProperties() {
			fmt.Printf("    %s: %s\n", prop.GetName(), prop.GetValue())
		}
	}
}
