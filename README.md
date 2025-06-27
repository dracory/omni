# Omni Package

A universal Go module providing core interfaces for composable primitives.

## Introduction

Omni is a powerful Go package that provides a flexible, composable architecture for building structured data models. It's designed around two core interfaces: `Property` for storing attribute-value pairs, and `Atom` for creating hierarchical, composable data structures.

## Core Concepts

### 1. Property
A `Property` is a fundamental building block that represents a single attribute-value pair. It's used to store metadata or configuration for atoms.

### 2. Atom
An `Atom` is a composable primitive that can have:
- A unique identifier
- A type
- Multiple properties
- Child atoms (enabling hierarchical structures)

### 3. Thread Safety
All operations on atoms are thread-safe, making it safe to use in concurrent applications.

## Architecture

The package follows a simple yet powerful architecture:

```
┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │
│     Atom        │────▶│   Property      │
│                 │     │                 │
└────────┬────────┘     └─────────────────┘
         │
         │ 0..*
         ▼
┌─────────────────┐
│                 │
│    Child Atom   │
│                 │
└─────────────────┘
```

## Interfaces

### PropertyInterface

`PropertyInterface` defines the contract for any abstract attribute or characteristic. Any concrete property type must implement this interface.

```go
type PropertyInterface interface {
    // GetName returns the name of the property
    GetName() string
    
    // SetName sets the name of the property
    SetName(name string)
    
    // GetValue returns the string representation of the property's value
    GetValue() string
    
    // SetValue sets the property's value from a string
    SetValue(value string)
}
```

### AtomInterface

`AtomInterface` is the universal interface that all composable primitives must satisfy. It defines the methods necessary for the system to understand and process any atom, regardless of its specific type.

```go
type AtomInterface interface {
    // GetID returns the unique identifier of the atom
    GetID() string
    
    // SetID sets the unique identifier of the atom
    SetID(id string)

    // GetType returns the type of the atom
    GetType() string
    
    // SetType sets the type of the atom
    SetType(atomType string)

    // GetProperties returns all properties of the atom
    GetProperties() []PropertyInterface
    
    // SetProperties sets all properties of the atom at once
    SetProperties(properties []PropertyInterface)

    // GetProperty returns a specific property by name, or nil if not found
    GetProperty(name string) PropertyInterface
    
    // SetProperty adds or updates a property
    SetProperty(property PropertyInterface)
    
    // RemoveProperty removes a property by name
    RemoveProperty(name string)

    // GetChildren returns all child atoms
    GetChildren() []AtomInterface
    
    // AddChild adds a new child atom
    AddChild(a AtomInterface)
}
```

## Installation

```bash
go get github.com/dracory/omni
```

## Basic Usage

```go
import "github.com/dracory/omni"

// Create a new property
prop := omni.NewProperty("color", "blue")

// Create a new atom
atom := omni.NewAtom("atom1", "button")
atom.SetProperty(prop)

// Create and add a child atom
child := omni.NewAtom("child1", "text")
child.SetProperty(omni.NewProperty("content", "Hello, World!"))
atom.AddChild(child)
```

## Advanced Usage Patterns

### 1. Hierarchical Structures
Create complex nested structures by adding child atoms:

```go
// Create a book with pages
book := omni.NewAtom("my_book", "book")
page1 := omni.NewAtom("page1", "page")
book.AddChild(page1)
```

### 2. Concurrent Access
Atoms are safe for concurrent use:

```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(i int) {
        defer wg.Done()
        prop := omni.NewProperty(fmt.Sprintf("prop_%d", i), fmt.Sprint(i))
        atom.SetProperty(prop)
    }(i)
}
wg.Wait()
```

### 3. Traversing Structures
Easily traverse and process atom hierarchies:

```go
func processAtom(atom omni.AtomInterface, indent string) {
    fmt.Printf("%s%s (%s)\n", indent, atom.GetID(), atom.GetType())
    for _, child := range atom.GetChildren() {
        processAtom(child, indent + "  ")
    }
}
```

## Examples

The package includes several examples demonstrating different usage patterns:

### 1. Basic Example
Demonstrates core functionality with properties and simple atom hierarchies.

```bash
cd examples/basic
go run basic_example.go
```

### 2. Advanced Example
Shows concurrent usage patterns and thread safety.

```bash
cd examples/advanced
go run concurrent_example.go
```

### 3. Book Example
Illustrates creating a hierarchical book structure with pages and content.

```bash
cd examples/book
go run book_example.go
```

### 4. Website Example
Demonstrates building a simple website structure with pages, headers, and paragraphs.

```bash
cd examples/website
# Show website structure
go run website_example.go
# Show specific page
go run website_example.go --page=home
go run website_example.go --page=about
```

## Best Practices

1. **Naming Conventions**:
   - Use lowercase with underscores for property names (e.g., `user_name`)
   - Keep atom types simple and descriptive (e.g., `button`, `user_profile`)

2. **Concurrency**:
   - The package handles concurrent access, but be mindful of deadlocks in your application logic
   - Use read locks for operations that only read atom state

3. **Memory Management**:
   - Large hierarchies should be managed carefully to avoid memory leaks
   - Consider implementing cleanup logic for long-running applications

## Performance Considerations

- Property lookups are O(n) where n is the number of properties
- Child atom lookups are O(n) where n is the number of children
- For high-performance scenarios, consider caching frequently accessed properties or children

## Contributing

Contributions are welcome! Please ensure all tests pass and add new tests for any new functionality.

## License

This project is licensed under the GNU Affero General Public License v3.0.
See the [LICENSE](LICENSE) file for the full license text.