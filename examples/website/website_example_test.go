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
		omni.WithProperties(map[string]string{
			"title": "Test Website",
		}),
	)

	// Create a test page
	site = createPage(site, "home", "/", "Home", "Welcome", "This is the home page")

	// Verify the page was added correctly
	pages := site.ChildrenGet()
	if len(pages) != 1 {
		t.Fatalf("Expected 1 page, got %d", len(pages))
	}

	page := pages[0]
	if page.GetID() != "home" {
		t.Errorf("Expected page ID to be 'home', got '%s'", page.GetID())
	}

	// Verify page properties
	props := page.GetAll()
	if title, exists := props["title"]; !exists || title != "Home" {
		t.Error("Expected page title to be 'Home'")
	}

	if uri, exists := props["uri"]; !exists || uri != "/" {
		t.Error("Expected page URI to be '/'")
	}

	// Verify page children (header and paragraph)
	children := page.ChildrenGet()
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
	site = createPage(site, "home", "/", "Home", "Welcome", "Home page content")
	site = createPage(site, "about", "/about", "About", "About Us", "About page content")
	site = createPage(site, "contact", "/contact", "Contact", "Contact Us", "Contact page content")

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

		title := page.Get("title")
		if title != tt.title {
			t.Errorf("Expected title '%s' for URI '%s', got '%s'",
				tt.title, tt.uri, title)
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
	site = createPage(site, "home", "/", "Home", "Welcome", "Home page content")
	site = createPage(site, "about", "/about", "About", "About Us", "About page content")

	// Get list of pages
	pages := listPages(site)

	// Verify the number of pages
	if len(pages) != 2 {
		t.Fatalf("Expected 2 pages, got %d", len(pages))
	}

	// Verify page titles
	expectedTitles := map[string]bool{"Home": true, "About": true}
	for _, page := range pages {
		title := page.Get("title")
		if !expectedTitles[title] {
			t.Errorf("Unexpected page title: %v", title)
		}
	}
}

// TestRenderPage tests the renderPage function
func TestRenderPage(t *testing.T) {
	// Create a test page
	site := omni.NewAtom("website")
	site = createPage(site, "test", "/test", "Test Page", "Test Header", "Test paragraph content")

	// Get the test page
	page := findPageByURI(site, "/test")
	if page == nil {
		t.Fatal("Failed to find test page")
	}

	// Render the page
	html := renderPage(page)

	// Basic verification of the rendered HTML
	if !strings.Contains(html, "<h1>Test Header</h1>") {
		t.Error("Expected HTML to contain header")
	}

	if !strings.Contains(html, "<p>Test paragraph content</p>") {
		t.Error("Expected HTML to contain paragraph content")
	}

	if !strings.Contains(html, "<title>Test Page</title>") {
		t.Error("Expected HTML to contain page title")
	}
}
