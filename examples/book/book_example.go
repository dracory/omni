// Book example demonstrating a hierarchical book structure with pages and content.
package main

import (
	"fmt"
	"github.com/dracory/omni"
)

// loremIpsum provides sample paragraphs for the book pages
var loremIpsum = []string{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
	"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	"Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.",
	"Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.",
	"Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem.",
}

// createPage creates a new page with the given number and content
func createPage(book *omni.Atom, pageNum int, content1, content2 string) {
	page := omni.NewAtom(fmt.Sprintf("page_%d", pageNum), "page")
	page.SetProperty(omni.NewProperty("number", fmt.Sprintf("%d", pageNum)))
	page.SetProperty(omni.NewProperty("content_1", content1))
	page.SetProperty(omni.NewProperty("content_2", content2))
	book.AddChild(page)
}

// printBook recursively prints the book structure
func printBook(atom omni.AtomInterface, indent string) {
	// Print current atom
	fmt.Printf("%s- %s (%s)\n", indent, atom.GetID(), atom.GetType())

	// Print properties with additional indentation
	for _, prop := range atom.GetProperties() {
		value := prop.GetValue()
		if len(value) > 50 { // Truncate long content
			value = value[:50] + "..."
		}
		fmt.Printf("%s  %s: %s\n", indent, prop.GetName(), value)
	}

	// Recursively print children
	for _, child := range atom.GetChildren() {
		printBook(child, indent+"  ")
	}
}

func main() {
	// Create a new book
	book := omni.NewAtom("my_book", "book")
	book.SetProperty(omni.NewProperty("title", "The Art of Go"))
	book.SetProperty(omni.NewProperty("author", "Gopher"))

	// Add some pages with lorem ipsum content
	for i := 0; i < 5; i++ {
		para1 := loremIpsum[i%len(loremIpsum)]
		para2 := loremIpsum[(i+1)%len(loremIpsum)]
		createPage(book, i+1, para1, para2)
	}

	// Print the book structure
	fmt.Println("Book Structure:")
	printBook(book, "")

	// Print a sample page
	if len(book.GetChildren()) > 0 {
		firstPage := book.GetChildren()[0]
		fmt.Println("\nSample Page Content:")
		fmt.Printf("Page %s\n", firstPage.GetProperty("number").GetValue())
		fmt.Println("---")
		fmt.Println(firstPage.GetProperty("content_1").GetValue())
		fmt.Println()
		fmt.Println(firstPage.GetProperty("content_2").GetValue())
	}
}
