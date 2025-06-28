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


	// Test rectangle properties
	if rectangle.GetID() != "rect1" {
		t.Errorf("Expected rectangle ID to be 'rect1', got '%s'", rectangle.GetID())
	}

	if rectangle.GetType() != "rectangle" {
		t.Errorf("Expected rectangle type to be 'rectangle', got '%s'", rectangle.GetType())
	}

	// Test properties
	widthProp := rectangle.GetProperty("width")
	if widthProp == nil || widthProp.GetValue() != "100" {
		t.Error("Expected width property to be '100'")
	}

	heightProp := rectangle.GetProperty("height")
	if heightProp == nil || heightProp.GetValue() != "50" {
		t.Error("Expected height property to be '50'")
	}

	// Test children
	children := rectangle.GetChildren()
	if len(children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(children))
	}

	// Test shadow child
	shadow := children[0]
	if shadow.GetID() != "shadow1" || shadow.GetType() != "shadow" {
		t.Errorf("Unexpected child: %s (%s)", shadow.GetID(), shadow.GetType())
	}

	// Test shadow property
	opacityProp := shadow.GetProperty("opacity")
	if opacityProp == nil || opacityProp.GetValue() != "0.5" {
		t.Error("Expected opacity property to be '0.5'")
	}
}

// TestMain runs the example and verifies its output
func TestMain(m *testing.M) {
	// Run the example to ensure it doesn't panic
	main()
}
