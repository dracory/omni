# Omni Examples

This directory contains example code demonstrating how to use the `omni` package.

## Basic Example

The [basic example](basic/basic_example.go) demonstrates the core functionality of the `omni` package, including:
- Creating and updating properties
- Creating atoms
- Adding properties to atoms
- Creating a hierarchy of atoms

To run the basic example:
```bash
cd basic
go run basic_example.go
```

## Advanced Example

The [advanced example](advanced/concurrent_example.go) shows more complex usage, including:
- Concurrent access to atoms and properties
- Working with multiple goroutines
- Aggregating and analyzing atom data

To run the advanced example:
```bash
cd advanced
go run concurrent_example.go
```

## Building and Running

1. Ensure you have Go installed on your system
2. Navigate to the example directory you want to run
3. Run the example with `go run <filename>.go`

## Dependencies

All examples depend on the `omni` package. Make sure it's available in your Go module or GOPATH.
