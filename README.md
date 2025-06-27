# Omni Package

A universal Go module providing core interfaces for composable primitives.

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

## Usage

1. Implement the `PropertyInterface` for custom property types
2. Implement the `AtomInterface` for custom atom types
3. Use the interfaces to create a flexible, composable system of atoms with properties