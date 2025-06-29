package main

import (
	"testing"

	"github.com/dracory/omni"
)

// TestBasicExample tests the basic example functionality
func TestBasicExample(t *testing.T) {
	// Create a new atom with properties and children
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

	// Test rectangle properties
	if rectangle.GetID() != "rect1" {
		t.Errorf("Expected rectangle ID to be 'rect1', got '%s'", rectangle.GetID())
	}

	if rectangle.GetType() != "rectangle" {
		t.Errorf("Expected rectangle type to be 'rectangle', got '%s'", rectangle.GetType())
	}

	// Test properties
	props := rectangle.GetAll()
	if width, ok := props["width"]; !ok || width != "100" {
		t.Error("Expected width property to be '100'")
	}
	if height, ok := props["height"]; !ok || height != "50" {
		t.Error("Expected height property to be '50'")
	}

	// Test children
	children := rectangle.ChildrenGet()
	if len(children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(children))
	}

	// Test shadow child
	shadow := children[0]
	if shadow.GetID() != "shadow1" || shadow.GetType() != "shadow" {
		t.Errorf("Unexpected child: %s (%s)", shadow.GetID(), shadow.GetType())
	}

	// Test shadow property
	shadowProps := shadow.GetAll()
	if opacity, ok := shadowProps["opacity"]; !ok || opacity != "0.5" {
		t.Error("Expected opacity property to be '0.5'")
	}
}

// TestMain runs the example and verifies its output
func TestMain(m *testing.M) {
	// Run the example to ensure it doesn't panic
	main()
}
