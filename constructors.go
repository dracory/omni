package omni

import (
	"errors"
	"fmt"

	"github.com/gouniverse/uid"
)

// NewAtom creates a new Atom with the given type and options.
// If no ID is provided, a human-readable UID will be generated.
func NewAtom(atomType string, opts ...AtomOption) *Atom {
	a := &Atom{
		id:         uid.HumanUid(),
		atomType:   atomType,
		properties: make([]PropertyInterface, 0),
		children:   make([]AtomInterface, 0),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// NewAtomFromGob creates a new Atom from binary data encoded with the gob package.
// The data should be a gob-encoded binary that was created by ToGob().
// Returns the Atom and nil error on success, or nil and an error if the data is invalid or cannot be decoded.
//
// Deprecated: Use GobToAtom instead which provides more robust validation and error handling.
func NewAtomFromGob(data []byte) (*Atom, error) {
	atom, err := GobToAtom(data)
	if err != nil {
		return nil, err
	}
	return atom.(*Atom), nil
}

// NewAtomFromJSON creates a new Atom from a JSON string.
// The JSON should be an object with at least "id" and "type" fields.
// Properties should be in a "properties" map, and children in a "children" array.
// Returns the Atom and nil error on success, or nil and an error if the JSON is invalid or cannot be unmarshaled.
//
// Note: For parsing multiple atoms from a JSON array, use JSONToAtoms instead.
func NewAtomFromJSON(jsonStr string) (*Atom, error) {
	if jsonStr == "" {
		return nil, errors.New("empty JSON string provided")
	}

	atoms, err := JSONToAtoms(jsonStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(atoms) == 0 {
		return nil, errors.New("no valid atom found in JSON")
	}

	// Return the first atom and ignore any additional ones
	atom, ok := atoms[0].(*Atom)
	if !ok {
		return nil, errors.New("unexpected atom type")
	}

	return atom, nil
}

// NewAtomFromMap creates a new Atom from a map representation.
// The map must contain at least "id" and "type" fields.
// Properties should be in a "properties" map, and children in a "children" slice.
// Returns the Atom and nil error on success, or nil and an error if the map is invalid.
//
// Deprecated: Use MapToAtom instead which provides more robust validation and error handling.
func NewAtomFromMap(atomMap map[string]any) (*Atom, error) {
	atom, err := MapToAtom(atomMap)
	if err != nil {
		return nil, err
	}
	return atom.(*Atom), nil
}

// AtomOption configures an Atom.
type AtomOption func(*Atom)

// WithID sets the ID of the Atom.
func WithID(id string) AtomOption {
	return func(a *Atom) {
		a.id = id
	}
}

// WithProperties adds properties to the Atom.
func WithProperties(properties ...PropertyInterface) AtomOption {
	return func(a *Atom) {
		a.properties = append(a.properties, properties...)
	}
}

// WithChildren adds child atoms to the Atom.
func WithChildren(children ...AtomInterface) AtomOption {
	return func(a *Atom) {
		a.children = append(a.children, children...)
	}
}
