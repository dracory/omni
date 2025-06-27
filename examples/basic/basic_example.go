// Basic example demonstrating the core functionality of the omni package.
package main

import (
	"fmt"
	"github.com/dracory/omni"
)

func main() {
	// Create a new property
	colorProp := omni.NewProperty("color", "blue")
	fmt.Printf("Property: %s = %s\n", colorProp.GetName(), colorProp.GetValue())

	// Update the property value
	colorProp.SetValue("red")
	fmt.Printf("Updated property: %s = %s\n", colorProp.GetName(), colorProp.GetValue())

	// Create a new atom
	rectangle := omni.NewAtom("rect1", "rectangle")
	rectangle.SetProperty(omni.NewProperty("width", "100"))
	rectangle.SetProperty(omni.NewProperty("height", "50"))

	// Add a child atom
	shadow := omni.NewAtom("shadow1", "shadow")
	shadow.SetProperty(omni.NewProperty("opacity", "0.5"))
	rectangle.AddChild(shadow)

	// Print the rectangle's properties
	fmt.Printf("\nRectangle '%s' properties:\n", rectangle.GetID())
	for _, prop := range rectangle.GetProperties() {
		fmt.Printf("  %s: %s\n", prop.GetName(), prop.GetValue())
	}

	// Print the rectangle's children
	fmt.Printf("\nRectangle '%s' children:\n", rectangle.GetID())
	for _, child := range rectangle.GetChildren() {
		fmt.Printf("  - %s (%s)\n", child.GetID(), child.GetType())
		for _, prop := range child.GetProperties() {
			fmt.Printf("    %s: %s\n", prop.GetName(), prop.GetValue())
		}
	}
}
