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

```bash
go get github.com/dracory/omni
```

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

### Serialization

Atoms can be easily serialized to various formats:

```go
// To map
m := atom.ToMap()

// To JSON
jsonStr := atom.ToJSON()

prettyJSON := atom.ToJSONPretty()

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

## Contributing

Contributions are welcome! Please ensure all tests pass and add new tests for any new functionality.

## License

This project is licensed under the GNU AGPLv3 - see the [LICENSE](LICENSE) file for details.