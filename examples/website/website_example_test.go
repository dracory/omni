package main

import (
	"strings"
	"testing"

	"github.com/dracory/omni"
)

// TestCreatePage tests the createPage helper function
func TestCreatePage(t *testing.T) {
	// Create a test site
	site := omni.NewAtom("website",
		omni.WithID("test_site"),
		omni.WithProperties(
			omni.NewProperty("title", "Test Website"),
		),
	)

	// Create a test page
	createPage(site, "home", "/", "Home", "Welcome", "This is the home page")

	// Verify the page was added correctly
	pages := site.GetChildren()
	if len(pages) != 1 {
		t.Fatalf("Expected 1 page, got %d", len(pages))
	}

	page := pages[0]
	if page.GetID() != "home" {
		t.Errorf("Expected page ID to be 'home', got '%s'", page.GetID())
	}

	// Verify page properties
	titleProp := page.GetProperty("title")
	if titleProp == nil || titleProp.GetValue() != "Home" {
		t.Error("Expected page title to be 'Home'")
	}

	uriProp := page.GetProperty("uri")
	if uriProp == nil || uriProp.GetValue() != "/" {
		t.Error("Expected page URI to be '/'")
	}

	// Verify page children (header and paragraph)
	children := page.GetChildren()
	if len(children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(children)) 
	}

	header := children[0]
	if header.GetType() != "header" || header.GetID() != "home_header" {
		t.Errorf("Unexpected header: %s (%s)", header.GetID(), header.GetType())
	}

	paragraph := children[1]
	if paragraph.GetType() != "paragraph" || paragraph.GetID() != "homeparagraph" {
		t.Errorf("Unexpected paragraph: %s (%s)", paragraph.GetID(), paragraph.GetType())
	}
}

// TestFindPageByURI tests the findPageByURI function
func TestFindPageByURI(t *testing.T) {
	// Create a test site with pages
	site := omni.NewAtom("website")
	createPage(site, "home", "/", "Home", "Welcome", "Home page content")
	createPage(site, "about", "/about", "About", "About Us", "About page content")
	createPage(site, "contact", "/contact", "Contact", "Contact Us", "Contact page content")

	// Test finding existing pages
	tests := []struct {
		uri   string
		title string
	}{
		{"/", "Home"},
		{"/about", "About"},
		{"/contact", "Contact"},
	}

	for _, tt := range tests {
		page := findPageByURI(site, tt.uri)
		if page == nil {
			t.Errorf("Expected to find page with URI '%s', but got nil", tt.uri)
			continue
		}

		titleProp := page.GetProperty("title")
		if titleProp == nil || titleProp.GetValue() != tt.title {
			t.Errorf("Expected title '%s' for URI '%s', got '%v'", 
				tt.title, tt.uri, titleProp)
		}
	}

	// Test non-existent page
	page := findPageByURI(site, "/nonexistent")
	if page != nil {
		t.Errorf("Expected nil for non-existent page, got %v", page)
	}
}

// TestListPages tests the listPages function
func TestListPages(t *testing.T) {
	// Create a test site with pages
	site := omni.NewAtom("website")
	createPage(site, "home", "/", "Home", "Welcome", "Home page content")
	createPage(site, "about", "/about", "About", "About Us", "About page content")

	// Get list of pages
	pages := listPages(site)

	// Verify the number of pages
	if len(pages) != 2 {
		t.Fatalf("Expected 2 pages, got %d", len(pages))
	}

	// Verify page titles
	expectedTitles := map[string]bool{"Home": true, "About": true}
	for _, page := range pages {
		titleProp := page.GetProperty("title")
		if titleProp == nil || !expectedTitles[titleProp.GetValue()] {
			t.Errorf("Unexpected page title: %v", titleProp)
		}
	}
}

// TestRenderPage tests the renderPage function
func TestRenderPage(t *testing.T) {
	// Create a test page
	page := omni.NewAtom("page",
		omni.WithID("test_page"),
		omni.WithProperties(
			omni.NewProperty("title", "Test Page"),
		),
	)

	// Add header and paragraph
	header := omni.NewAtom("header",
		omni.WithID("test_header"),
		omni.WithProperties(
			omni.NewProperty("text", "Test Header"),
		),
	)
	page.AddChild(header)

	paragraph := omni.NewAtom("paragraph",
		omni.WithID("test_paragraph"),
		omni.WithProperties(
			omni.NewProperty("content", "Test content"),
		),
	)
	page.AddChild(paragraph)

	// Render the page
	html := renderPage(page)

	// Basic verification of the rendered HTML
	if !strings.Contains(html, "<h1>Test Header</h1>") {
		t.Error("Expected HTML to contain header")
	}

	if !strings.Contains(html, "<p>Test content</p>") {
		t.Error("Expected HTML to contain paragraph")
	}

	if !strings.Contains(html, "<title>Test Page</title>") {
		t.Error("Expected HTML to contain page title")
	}
}

// TestMain runs the example and verifies it doesn't panic
func TestMain(m *testing.M) {
	// Run the example to ensure it doesn't panic
	main()
}
