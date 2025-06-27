# Omni Package

A universal Go module providing core interfaces for composable primitives.


<img src="https://opengraph.githubassets.com/5b92c81c05d64a82c3fb4ba95739403a2d38cbad61f260a0701b3366b3d10327/dracory/omni" />

[![Tests Status](https://github.com/dracory/omni/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/omni/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/omni)](https://goreportcard.com/report/github.com/dracory/omni)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/omni)](https://pkg.go.dev/github.com/dracory/omni)

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

## Usage

1. Implement the `PropertyInterface` for custom property types
2. Implement the `AtomInterface` for custom atom types
3. Use the interfaces to create a flexible, composable system of atoms with properties

## Examples

Check out the [examples](./examples) directory for complete, runnable examples:

- [Basic Example](./examples/basic/basic_example.go) - Demonstrates core functionality
- [Advanced Example](./examples/advanced/concurrent_example.go) - Shows concurrent usage patterns

To run the examples:

```bash
# Basic example
cd examples/basic
go run basic_example.go

# Advanced example
cd examples/advanced
go run concurrent_example.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
