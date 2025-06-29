// Website example demonstrating a simple website structure with pages and content.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dracory/omni"
)

// createPage creates a new page with a header and paragraph
func createPage(site *omni.Atom, pageID, uri, title, headerText, paragraphText string) {
	// Create the page
	page := omni.NewAtom("page",
		omni.WithID(pageID),
		omni.WithProperties(map[string]string{
			"title": title,
			"uri":   uri,
		}),
	)

	// Add header
	header := omni.NewAtom("header",
		omni.WithID(pageID+"_header"),
		omni.WithProperties(map[string]string{
			"text": headerText,
		}),
	)
	page = page.ChildAdd(header).(*omni.Atom)

	// Add paragraph
	paragraph := omni.NewAtom("paragraph",
		omni.WithID(pageID+"paragraph"),
		omni.WithProperties(map[string]string{
			"content": paragraphText,
		}),
	)
	page = page.ChildAdd(paragraph).(*omni.Atom)

	// Add page to site
	site = site.ChildAdd(page).(*omni.Atom)
}

// printWebsite recursively prints the website structure
func printWebsite(atom omni.AtomInterface, indent string) {
	// Print current atom
	fmt.Printf("%s- %s (%s)", indent, atom.GetID(), atom.GetType())

	// Print properties if any
	props := atom.GetAll()
	if len(props) > 0 {
		fmt.Print(" - Properties: ")
		first := true
		for key, value := range props {
			if !first {
				fmt.Print(", ")
			}
			fmt.Printf("%s: %q", key, value)
			first = false
		}
	}
	fmt.Println()

	// Recursively print children
	for _, child := range atom.ChildrenGet() {
		printWebsite(child, indent+"  ")
	}
}

// renderPage renders a single page as HTML
func renderPage(page omni.AtomInterface) string {
	if page.GetType() != "page" {
		return ""
	}

	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	if title := page.Get("title"); title != "" {
		html += fmt.Sprintf("  <title>%s</title>\n", title)
	}
	html += "</head>\n<body>\n"

	// Add header and content
	for _, child := range page.ChildrenGet() {
		switch child.GetType() {
		case "header":
			if text := child.Get("text"); text != "" {
				html += fmt.Sprintf("  <h1>%s</h1>\n", text)
			}
		case "paragraph":
			if content := child.Get("content"); content != "" {
				html += fmt.Sprintf("  <p>%s</p>\n", content)
			}
		}
	}

	html += "</body>\n</html>"
	return html
}

// findPageByURI finds a page by its URI in the website
func findPageByURI(site *omni.Atom, uri string) omni.AtomInterface {
	for _, page := range site.ChildrenGet() {
		if pageURI := page.Get("uri"); pageURI == uri {
			return page
		}
	}
	return nil
}

// listPages returns a list of all available pages with their URIs
func listPages(site *omni.Atom) []omni.AtomInterface {
	var pages []omni.AtomInterface
	for _, page := range site.ChildrenGet() {
		if page.GetType() == "page" {
			pages = append(pages, page)
		}
	}
	return pages
}

func main() {
	// Parse command line flags
	port := flag.Int("port", 8080, "Port to run the server on")
	flag.Parse()

	// Create a new website
	site := omni.NewAtom("website",
		omni.WithID("my_website"),
		omni.WithProperties(map[string]string{
			"title": "My Awesome Website",
		}),
	)

	// Add some pages
	createPage(site, "home", "/", "Home", "Welcome to My Website", "This is the home page of my awesome website.")
	createPage(site, "about", "/about", "About", "About Us", "We are a company that builds amazing things with Go!")
	createPage(site, "contact", "/contact", "Contact", "Get in Touch", "Email us at contact@example.com")

	// Print the website structure
	fmt.Println("Website structure:")
	printWebsite(site, "  ")

	// Set up HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := findPageByURI(site, r.URL.Path)
		if page == nil {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, renderPage(page))
	})

	// Start the web server in a goroutine so it doesn't block
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: nil, // Use default mux
	}

	// Channel to listen for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Printf("\nStarting web server on http://localhost:%d\n", *port)
		fmt.Println("Available pages:")
		for _, page := range listPages(site) {
			uri := page.Get("uri")
			title := page.Get("title")
			fmt.Printf("  - %s: %s\n", uri, title)
		}
		fmt.Println("\nPress Ctrl+C to stop the server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	fmt.Println("\nShutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server gracefully stopped")
}
