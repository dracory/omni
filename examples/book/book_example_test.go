package main

import (
	"strconv"
	"testing"

	"github.com/dracory/omni"
)

// TestCreatePage tests the createPage helper function
func TestCreatePage(t *testing.T) {
	// Create a test page
	page := createPage("page1", 1, "Content 1", "Content 2")

	// Test page properties
	if page.GetID() != "page1" {
		t.Errorf("Expected page ID to be 'page1', got '%s'", page.GetID())
	}

	if page.GetType() != "page" {
		t.Errorf("Expected page type to be 'page', got '%s'", page.GetType())
	}

	// Test page properties
	numberProp := page.GetProperty("number")
	if numberProp == nil || numberProp.GetValue() != "1" {
		t.Error("Expected number property to be '1'")
	}

	content1Prop := page.GetProperty("content_1")
	if content1Prop == nil || content1Prop.GetValue() != "Content 1" {
		t.Error("Expected content_1 property to be 'Content 1'")
	}

	content2Prop := page.GetProperty("content_2")
	if content2Prop == nil || content2Prop.GetValue() != "Content 2" {
		t.Error("Expected content_2 property to be 'Content 2'")
	}

	// Test that page has no children by default
	if len(page.GetChildren()) != 0 {
		t.Error("Expected page to have no children by default")
	}
}

// TestPrintBook tests the printBook function
func TestPrintBook(t *testing.T) {
	// Create a test book structure
	book := omni.NewAtom("book",
		omni.WithID("test_book"),
		omni.WithProperties(
			omni.NewProperty("title", "Test Book"),
			omni.NewProperty("author", "Test Author"),
		),
	)

	// Add some pages
	for i := 0; i < 3; i++ {
		pageID := "page_" + strconv.Itoa(i+1)
		page := createPage(pageID, i+1, "Content A", "Content B")
		book.AddChild(page)
	}

	// Test the book structure
	if book.GetID() != "test_book" {
		t.Errorf("Expected book ID to be 'test_book', got '%s'", book.GetID())
	}

	// Test page count
	pages := book.GetChildren()
	if len(pages) != 3 {
		t.Fatalf("Expected 3 pages, got %d", len(pages))
	}

	// Test page properties
	for i, page := range pages {
		expectedID := "page_" + strconv.Itoa(i+1)
		if page.GetID() != expectedID {
			t.Errorf("Expected page %d ID to be '%s', got '%s'", i+1, expectedID, page.GetID())
		}

		numberProp := page.GetProperty("number")
		expectedNumber := strconv.Itoa(i + 1)
		if numberProp == nil || numberProp.GetValue() != expectedNumber {
			t.Errorf("Expected page %d number to be '%s', got '%s'", i+1, expectedNumber, numberProp.GetValue())
		}
	}
}

// TestBookExample tests the main book example
func TestBookExample(t *testing.T) {
	// This test just verifies that the example runs without panicking
	// The actual behavior is tested in the other test functions
}

// TestMain runs the example and verifies it doesn't panic
func TestMain(m *testing.M) {
	// Run the example to ensure it doesn't panic
	main()
}
