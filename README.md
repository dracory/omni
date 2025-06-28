# Omni Package

A universal Go module providing core interfaces and implementations for composable primitives with a functional options pattern for clean and flexible instantiation.

<img src="https://opengraph.githubassets.com/5b92c81c05d64a82c3fb4ba95739403a2d38cbad61f260a0701b3366b3d10327/dracory/omni" />

[![Tests Status](https://github.com/dracory/omni/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/omni/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/omni)](https://goreportcard.com/report/github.com/dracory/omni)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/omni)](https://pkg.go.dev/github.com/dracory/omni)


## Introduction

Omni is a powerful Go package that provides a flexible, composable architecture for building structured data models. It's designed around two core interfaces: `Property` for storing attribute-value pairs, and `Atom` for creating hierarchical, composable data structures.

## Core Concepts

### 1. Property
A `Property` is a fundamental building block that represents a single attribute-value pair. It's used to store metadata or configuration for atoms.

### 2. Atom
An `Atom` is a composable primitive that can have:
- A unique identifier (auto-generated if not provided)
- A type (required)
- Multiple properties
- Child atoms (enabling hierarchical structures)

### 3. Functional Options Pattern
Atoms are created using a functional options pattern for clean and flexible instantiation:

```go
// Basic creation with just type (auto-generated ID)
atom := NewAtom("my-type")

// With custom ID
atom := NewAtom("my-type", WithID("custom-id"))

// With properties
atom := NewAtom("my-type", 
    WithID("my-id"),
    WithProperties([]PropertyInterface{
        NewProperty("name", "value"),
    }),
)

// With children
atom := NewAtom("parent",
    WithChildren([]AtomInterface{
        NewAtom("child", WithID("child1")),
        NewAtom("child", WithID("child2")),
    }),
)
```

### 4. Thread Safety
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

## Installation

### Using Go Modules (recommended)

```bash
# Add to your project
go get github.com/dracory/omni

# Update to the latest version
go get -u github.com/dracory/omni
```

### In your Go code

```go
import "github.com/dracory/omni"
```

### Requirements

- Go 1.18 or higher (for generics support)
- No external dependencies beyond the Go standard library

## Getting Started

```go
import "github.com/dracory/omni"

// Create a new property
prop := omni.NewProperty("color", "blue")


// Create a new atom with functional options
atom := omni.NewAtom("my-type",
    omni.WithID("my-atom"),
    omni.WithProperties([]omni.PropertyInterface{prop}),
)

// Or build it up
atom := omni.NewAtom("my-type")
atom.SetProperty(prop)

// Add children
child := omni.NewAtom("child", omni.WithID("child-1"))
atom.AddChild(child)
```

## Examples

Check out the `examples` directory for complete working examples:

1. `basic/` - Basic usage of atoms and properties
2. `book/` - A more complex example modeling a book with chapters
3. `advanced/` - Concurrent operations and advanced patterns
4. `website/` - A website structure with pages and components

## Thread Safety

All operations on atoms are thread-safe, using mutexes to protect concurrent access. The implementation is designed for high concurrency with minimal lock contention.

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