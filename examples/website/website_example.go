// Website example demonstrating a simple website structure with pages and content.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dracory/omni"
)

// createPage creates a new page with a header and paragraph
func createPage(site *omni.Atom, pageID, uri, title, headerText, paragraphText string) {
	// Create the page
	page := omni.NewAtom(pageID, "page")
	page.SetProperty(omni.NewProperty("title", title))
	page.SetProperty(omni.NewProperty("uri", uri))

	// Add header
	header := omni.NewAtom(pageID+"_header", "header")
	header.SetProperty(omni.NewProperty("text", headerText))
	page.AddChild(header)

	// Add paragraph
	paragraph := omni.NewAtom(pageID+"_paragraph", "paragraph")
	paragraph.SetProperty(omni.NewProperty("content", paragraphText))
	page.AddChild(paragraph)

	// Add page to site
	site.AddChild(page)
}

// printWebsite recursively prints the website structure
func printWebsite(atom omni.AtomInterface, indent string) {
	// Print current atom
	fmt.Printf("%s- %s (%s)", indent, atom.GetID(), atom.GetType())

	// Print properties if any
	if props := atom.GetProperties(); len(props) > 0 {
		fmt.Print(" - Properties: ")
		for i, prop := range props {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s: %q", prop.GetName(), prop.GetValue())
		}
	}
	fmt.Println()

	// Recursively print children
	for _, child := range atom.GetChildren() {
		printWebsite(child, indent+"  ")
	}
}

// renderPage renders a single page as HTML
func renderPage(page omni.AtomInterface) string {
	if page.GetType() != "page" {
		return ""
	}

	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	if title := page.GetProperty("title"); title != nil {
		html += fmt.Sprintf("  <title>%s</title>\n", title.GetValue())
	}
	html += "</head>\n<body>\n"

	// Add header and content
	for _, child := range page.GetChildren() {
		switch child.GetType() {
		case "header":
			if text := child.GetProperty("text"); text != nil {
				html += fmt.Sprintf("  <h1>%s</h1>\n", text.GetValue())
			}
		case "paragraph":
			if content := child.GetProperty("content"); content != nil {
				html += fmt.Sprintf("  <p>%s</p>\n", content.GetValue())
			}
		}
	}

	html += "</body>\n</html>"
	return html
}

// findPageByURI finds a page by its URI in the website
func findPageByURI(site *omni.Atom, uri string) omni.AtomInterface {
	// Normalize the URI to ensure it starts with /
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}

	for _, page := range site.GetChildren() {
		if page.GetType() == "page" {
			if uriProp := page.GetProperty("uri"); uriProp != nil && uriProp.GetValue() == uri {
				return page
			}
		}
	}
	return nil
}

// listPages returns a list of all available pages with their URIs
func listPages(site *omni.Atom) []omni.AtomInterface {
	var pages []omni.AtomInterface
	for _, child := range site.GetChildren() {
		if child.GetType() == "page" {
			pages = append(pages, child)
		}
	}
	return pages
}

func main() {
	// Parse command line flags
	pageFlag := flag.String("page", "", "Page to display (e.g., 'home' or 'about')")
	flag.Parse()

	// Create the website root
	site := omni.NewAtom("my_website", "website")
	site.SetProperty(omni.NewProperty("name", "My Awesome Site"))

	// Add home page
	createPage(site, "home", "/", "Home", "Welcome to My Website",
		"This is the home page of my awesome website. Here you can find all the latest news and updates.")

	// Add about page
	createPage(site, "about", "/about", "About Us", "About Our Company",
		"We are a small team of passionate developers creating amazing web experiences. "+
			"Our mission is to make the web a better place, one website at a time.")

	// Handle page display based on command line flag
	if *pageFlag != "" {
		// Try to find the requested page
		var pageURI string
		if *pageFlag == "home" || *pageFlag == "" {
			pageURI = "/"
		} else if !strings.HasPrefix(*pageFlag, "/") {
			pageURI = "/" + *pageFlag
		} else {
			pageURI = *pageFlag
		}

		page := findPageByURI(site, pageURI)
		if page != nil {
			fmt.Println(renderPage(page))
		} else {
			fmt.Printf("Page '%s' not found. Available pages:\n", pageURI)
			for _, p := range listPages(site) {
				uri := p.GetProperty("uri").GetValue()
				title := p.GetProperty("title").GetValue()
				fmt.Printf("  %s - %s\n", uri, title)
			}
			os.Exit(1)
		}
	} else {
		// No page specified, show the website structure
		fmt.Println("Website Structure:")
		printWebsite(site, "")

		// Also show usage
		fmt.Println("\nUsage:")
		fmt.Println("  go run website_example.go --page=home")
		fmt.Println("  go run website_example.go --page=about")
		fmt.Println("  go run website_example.go --page=/")
		fmt.Println("  go run website_example.go --page=/about")
	}
}
