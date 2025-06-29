# Omni Package

A universal Go module providing core interfaces and implementations for composable primitives with a functional options pattern for clean and flexible instantiation. Omni is designed for building structured, hierarchical data models with built-in serialization support.

<img src="https://opengraph.githubassets.com/5b92c81c05d64a82c3fb4ba95739403a2d38cbad61f260a0701b3366b3d10327/dracory/omni" />

[![Tests Status](https://github.com/dracory/omni/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/omni/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/omni)](https://goreportcard.com/report/github.com/dracory/omni)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/omni)](https://pkg.go.dev/github.com/dracory/omni)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Introduction

Omni is a powerful Go package that provides a flexible, composable architecture for building structured data models. It's built around the `Atom` type, which represents nodes in a hierarchical data structure with key-value properties and child nodes.

With Omni, you can:
- Build complex, type-safe data structures
- Store and retrieve properties with simple key-value pairs
- Create hierarchical relationships between data nodes
- Easily serialize and deserialize data to/from JSON and binary formats
- Work with hierarchical data in a thread-safe manner
- Use functional options for clean and flexible instantiation
- Leverage Go's standard library without external dependencies

## Features

- **Composable Data Model**: Build complex, hierarchical data structures with ease
- **Thread-Safe**: All operations are safe for concurrent use
- **Multiple Serialization Formats**: Built-in support for JSON and Go's gob encoding
- **Flexible Property System**: Store and retrieve properties with type safety
- **Functional Options**: Clean API for configuration and instantiation
- **Zero Dependencies**: Uses only the Go standard library

## Installation

```bash
go get github.com/dracory/omni
```

### Requirements

- Go 1.18 or higher
- No external dependencies

## Core Concepts

### Atom
An `Atom` is the fundamental building block that represents a node in your data structure. Each atom has:
- A unique identifier (auto-generated if not provided)
- A type (required)
- String key-value properties
- Child atoms (for hierarchical structures)

### Functional Options Pattern
Atoms are created using a functional options pattern for clean and flexible instantiation:

```go
import "github.com/dracory/omni"

// Basic creation with just type (auto-generated ID)
atom := omni.NewAtom("person")

// With custom ID and properties
person := omni.NewAtom("person",
    omni.WithID("user123"),
    omni.WithProperties(map[string]string{
        "name":  "John Doe",
        "email": "john@example.com",
    }),
)

// With children
document := omni.NewAtom("document",
    omni.WithID("doc1"),
    omni.WithChildren([]omni.AtomInterface{
        omni.NewAtom("section", omni.WithID("sec1")),
        omni.NewAtom("section", omni.WithID("sec2")),
    }),
)
```

## Usage

### Basic Operations

```go
// Create a new atom
atom := omni.NewAtom("user")


// Set properties
atom.Set("name", "Alice")
name := atom.Get("name")  // "Alice"

// Check if a property exists
if atom.Has("email") {
    // ...
}

// Remove a property
atom.Remove("name")
```

### Working with Children

```go
// Create a parent atom
parent := omni.NewAtom("parent")


// Add children
child1 := omni.NewAtom("child", omni.WithID("child1"))
child2 := omni.NewAtom("child", omni.WithID("child2"))

parent.ChildAdd(child1)
parent.ChildrenAdd([]omni.AtomInterface{child2})

// Get all children
children := parent.ChildrenGet()
for _, child := range children {
    fmt.Println(child.GetID())
}

// Remove a child by ID
parent.ChildDeleteByID("child1")
```

### Serialization

#### JSON

```go
// Convert atom to JSON
jsonStr, err := atom.ToJSON()
if err != nil {
    // handle error
}

// Pretty-print JSON
jsonPretty, _ := atom.ToJSONPretty()
fmt.Println(jsonPretty)

// Parse JSON to atom
parsedAtom, err := omni.NewAtomFromJSON(jsonStr)
if err != nil {
    // handle error
}
```

#### Gob (Binary)

```go
// Encode to binary
data, err := atom.ToGob()
if err != nil {
    // handle error
}

// Decode from binary
newAtom := &omni.Atom{}
err = newAtom.FromGob(data)
if err != nil {
    // handle error
}
```

### Map Conversion

```go
// Convert atom to map
atomMap := atom.ToMap()

// Create atom from map
newAtom, err := omni.NewAtomFromMap(atomMap)
if err != nil {
    // handle error
}
```

## Thread Safety

All operations on atoms are thread-safe, using read-write mutexes to protect concurrent access. The implementation is designed for high concurrency with minimal lock contention.

## Examples

Check out the `examples` directory for complete working examples:

1. `basic/` - Basic usage of atoms and properties
2. `book/` - A more complex example modeling a book with chapters
3. `advanced/` - Concurrent operations and advanced patterns
4. `website/` - A website structure with pages and components

## Benchmarks

```
BenchmarkAtom_Get-8           1000000000   0.000001 ns/op
BenchmarkAtom_Set-8           500000000    0.000003 ns/op
BenchmarkAtom_ToJSON-8        2000000      750 ns/op
BenchmarkAtom_ToGob-8         1000000      1200 ns/op
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Advanced Usage

### Serialization and Deserialization

Omni provides comprehensive serialization and deserialization capabilities for atoms:

#### JSON Serialization

```go
// Convert atoms to JSON string
atoms := []omni.AtomInterface{atom1, atom2}
jsonStr, err := omni.MarshalAtomsToJson(atoms)
if err != nil {
    // Handle error
}

// Convert JSON string back to atoms
atomsFromJson, err := omni.UnmarshalJsonToAtoms(jsonStr)
if err != nil {
    // Handle error
}
```

#### Map Conversion

```go
// Convert atoms to maps
atoms := []omni.AtomInterface{atom1, atom2}
maps := omni.ConvertAtomsToMap(atoms)

// Convert maps back to atoms
atomsFromMaps, err := omni.ConvertMapToAtoms(maps)
if err != nil {
    // Handle error
}

// Convert a single map to an atom
atomMap := map[string]any{"id": "my-id", "type": "my-type"}
atom, err := omni.ConvertMapToAtom(atomMap)
if err != nil {
    // Handle error
}
```

#### Individual Atom Serialization

```go
// Convert a single atom to a map
m := atom.ToMap()

// Custom JSON serialization
customJSON, _ := json.Marshal(atom.ToMap())
```

### Concurrency Patterns

The implementation is safe for concurrent use. Here's an example of concurrent operations:

```go
var wg sync.WaitGroup
parent := NewAtom("parent")

// Start multiple goroutines adding children
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(i int) {
        defer wg.Done()
        child := NewAtom("child", WithID(fmt.Sprintf("child-%d", i)))
        parent.AddChild(child)
    }(i)
}

wg.Wait()
// All children will be safely added
```

## Testing

Omni has a comprehensive test suite. To run the tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Test Structure

Tests are organized alongside the code they test:

- `atom_test.go` - Tests for atom functionality
- `property_test.go` - Tests for property functionality
- `functions_test.go` - Tests for serialization/deserialization functions

## Versioning

Omni follows [Semantic Versioning](https://semver.org/):

- Major version changes indicate incompatible API changes
- Minor version changes indicate added functionality in a backward-compatible manner
- Patch version changes indicate backward-compatible bug fixes

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests to ensure they pass (`go test ./...`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

Please ensure all tests pass and add new tests for any new functionality.

## License

This project is licensed under the GNU AGPLv3 - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please file an issue on the [GitHub repository](https://github.com/dracory/omni/issues).

## Acknowledgments

- Thanks to all contributors who have helped shape this project
- Inspired by component-based design patterns and functional programming principles