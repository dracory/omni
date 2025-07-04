# Omni vs gouniverse/dataobject

This document provides a detailed comparison between Omni and [gouniverse/dataobject](https://github.com/gouniverse/dataobject), highlighting their differences, use cases, and when to choose one over the other.

## Core Concepts

### Omni
- **Hierarchical Data Model**: Built around composable atoms that can have parent-child relationships
- **String-based Properties**: Simple key-value storage with string values
- **Thread Safety**: Built-in synchronization for concurrent access
- **Functional Options**: Clean API using functional options pattern
- **Tree Structures**: Optimized for representing hierarchical data
- **Serialization**: Supports JSON and Gob formats

### gouniverse/dataobject
- **Flat Data Model**: Simple key-value store
- **String-based**: All values are stored as strings
- **Change Tracking**: Built-in dirty flag and change tracking
- **Serialization**: Supports both JSON and Gob formats
- **Simple API**: Traditional getter/setter pattern

## Feature Comparison

| Feature | Omni | gouniverse/dataobject |
|---------|------|----------------------|
| **Data Model** | Hierarchical tree structure | Flat key-value store |
| **Property Storage** | Map of string key-value pairs | Map of string key-value pairs |
| **Thread Safety** | Built-in (thread-safe) | Not thread-safe by default |
| **Change Tracking** | No built-in tracking | Built-in dirty flag |
| **Serialization** | JSON and Gob formats | JSON and Gob formats |
| **Performance** | Slightly slower (mutex overhead) | Faster for simple operations |
| **Dependencies** | Standard library only | Standard library only |
| **API Style** | Method chaining, functional options | Traditional getter/setter |

## Code Examples

### Creating and Using an Object

**Omni**
```go
// Create an atom with properties
atom := omni.NewAtom("person",
    omni.WithID("123"),
    omni.WithProperties(map[string]string{
        "name": "John Doe",
        "age": "30",
    }),
)

// Add a child atom
child := omni.NewAtom("address",
    omni.WithProperties(map[string]string{
        "street": "123 Main St",
    }),
)
atom.ChildAdd(child)

// Set/Get properties
atom.Set("email", "john@example.com")
name := atom.Get("name") // "John Doe"
```

**gouniverse/dataobject**
```go
// Create a data object
do := dataobject.New()

// Set values
do.Set("name", "John Doe")
do.Set("age", "30")
do.Set("address.street", "123 Main St") // Flat structure with dot notation
```

### Serialization

**Omni**
```go
// To JSON
jsonStr, err := atom.ToJSON()

// From JSON
parsedAtom, err := omni.NewAtomFromJSON(jsonStr)

// To Gob
gobData, err := atom.ToGob()

// From Gob
newAtom := &omni.Atom{}
err = newAtom.FromGob(gobData)
```

**gouniverse/dataobject**
```go
// To JSON
jsonStr, err := do.ToJSON()

// From JSON
do, err := dataobject.NewFromJSON(jsonStr)

// To Gob
gobData, err := do.ToGob()

// From Gob
do, err := dataobject.NewFromGob(gobData)
```

## When to Use Omni

1. **Hierarchical Data**
   - When you need to represent tree-like structures
   - For component-based architectures
   - When parent-child relationships are important

2. **Thread Safety**
   - When multiple goroutines will access the data structure
   - For thread-safe operations out of the box

3. **Complex Data Models**
   - For rich domain models with nested structures
   - When you need to model complex relationships

4. **Functional API**
   - When you prefer a fluent, chainable API
   - For clean configuration with functional options

## When to Use gouniverse/dataobject

1. **Simple Key-Value Storage**
   - For configuration data
   - When you just need a simple dictionary/map with change tracking

2. **Flat Data Structures**
   - When working with tabular data
   - For simple DTOs (Data Transfer Objects)

3. **Performance-Critical Code**
   - When you need maximum performance for simple operations
   - For high-throughput scenarios where the overhead of mutexes isn't needed

4. **Change Tracking**
   - When you need built-in dirty flag functionality
   - For forms or UIs that track changes

## Migration Between the Two

### From dataobject to Omni

```go
// dataobject to Omni
func ConvertDataObjectToOmni(do *dataobject.DataObject) *omni.Atom {
    atom := omni.NewAtom("dataobject")
    
    // Copy all fields as properties
    for k, v := range do.Data() {
        atom.Set(k, v)
    }
    
    return atom
}
```

### From Omni to dataobject

```go
// Omni to dataobject
func ConvertOmniToDataObject(atom omni.AtomInterface) *dataobject.DataObject {
    do := dataobject.New()
    
    // Copy all properties
    for k, v := range atom.GetAll() {
        do.Set(k, v)
    }
    
    // Note: This flattens the hierarchy
    // For hierarchical data, you might want to use dot notation or another scheme
    
    return do
}
```

## Conclusion

Both Omni and gouniverse/dataobject serve different purposes and have their own strengths. Choose Omni when you need hierarchical data structures, thread safety, and a functional API. Opt for gouniverse/dataobject when you need a simple, flat key-value store with change tracking.

For applications that deal with complex domain models or component-based architectures, Omni provides the necessary features for structured data. For simpler use cases where you just need a dynamic key-value store with change tracking, gouniverse/dataobject might be more appropriate.
