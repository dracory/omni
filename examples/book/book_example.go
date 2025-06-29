// Book example demonstrating a hierarchical book structure with pages and content.
package main

import (
	"fmt"
	"strconv"

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
func createPage(id string, pageNum int, content1, content2 string) omni.AtomInterface {
	return omni.NewAtom("page",
		omni.WithID(id),
		omni.WithProperties(map[string]string{
			"number":    strconv.Itoa(pageNum),
			"content_1": content1,
			"content_2": content2,
		}),
	)
}

// printBook recursively prints the book structure
func printBook(atom omni.AtomInterface, indent int) {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	fmt.Printf("%s- %s (type: %s)\n", prefix, atom.GetID(), atom.GetType())
	props := atom.GetAll()
	for key, value := range props {
		fmt.Printf("%s  %s: %v\n", prefix, key, value)
	}

	for _, child := range atom.ChildrenGet() {
		printBook(child, indent+1)
	}
}

func main() {
	// Create pages with properties
	var pages []omni.AtomInterface
	pages = make([]omni.AtomInterface, 0, 5)
	for i := 1; i <= 5; i++ {
		para1 := loremIpsum[i%len(loremIpsum)]
		para2 := loremIpsum[(i+1)%len(loremIpsum)]
		page := createPage(
			fmt.Sprintf("page_%d", i),
			i,
			para1,
			para2,
		)
		pages = append(pages, page)
	}

	// Create book with pages using functional options
	book := omni.NewAtom("book",
		omni.WithID("my_book"),
		omni.WithChildren(pages...),
		omni.WithProperties(map[string]string{
			"title":  "The Art of Go",
			"author": "Gopher",
		}),
	)

	// Print the book structure
	fmt.Println("Book structure:")
	printBook(book, 0)

	// Convert to JSON and print
	jsonData, _ := book.ToJSON()
	fmt.Println("\nBook as JSON:")
	fmt.Println(string(jsonData))

	// Print a sample page
	children := book.ChildrenGet()
	if len(children) > 0 {
		firstPage := children[0]
		pageProps := firstPage.GetAll()
		fmt.Println("\nSample Page Content:")
		fmt.Printf("Page %s\n", pageProps["number"])
		fmt.Println("---")
		fmt.Println(pageProps["content_1"])
		fmt.Println()
		fmt.Println(pageProps["content_2"])
	}
}
