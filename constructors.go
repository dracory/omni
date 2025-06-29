package omni

import (
	"github.com/gouniverse/uid"
)

// NewAtom creates a new Atom with the given type and applies the provided options.
// If no ID is provided via options, a human-readable UID will be generated.
func NewAtom(atomType string, opts ...AtomOption) *Atom {
	atom := &Atom{
		properties: map[string]string{
			"id":   uid.HumanUid(),
			"type": atomType,
		},
		children: make([]AtomInterface, 0),
	}

	for _, opt := range opts {
		opt(atom)
	}

	return atom
}

// NewAtomFromGob creates a new Atom from binary data encoded with the gob package.
// The data should be a gob-encoded binary that was created by ToGob().
// Returns the Atom and nil error on success, or nil and an error if the data is invalid or cannot be decoded.
//
// Deprecated: Use GobToAtom instead which provides more robust validation and error handling.
func NewAtomFromGob(data []byte) (*Atom, error) {
	return GobToAtom(data)
}

// NewAtomFromJSON creates a new Atom from a JSON string.
// The JSON should be an object with at least "id" and "type" fields.
// Properties should be in a "properties" map, and children in a "children" array.
// Returns the Atom and nil error on success, or nil and an error if the JSON is invalid or cannot be unmarshaled.
//
// Note: For parsing multiple atoms from a JSON array, use JSONToAtoms instead.
//
// Deprecated: Use JSONToAtom instead which provides more robust validation and error handling.
func NewAtomFromJSON(jsonStr string) (*Atom, error) {
	atom, err := JSONToAtom(jsonStr)
	if err != nil {
		return nil, err
	}
	return atom.(*Atom), nil
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
		a.Set("id", id)
	}
}

// WithProperties adds properties to the Atom.
func WithProperties(properties map[string]string) AtomOption {
	return func(a *Atom) {
		for k, v := range properties {
			a.Set(k, v)
		}
	}
}

// WithChildren adds child atoms to the Atom.
func WithChildren(children ...AtomInterface) AtomOption {
	return func(a *Atom) {
		a.children = append(a.children, children...)
	}
}

// WithType sets the type of the Atom.
func WithType(atomType string) AtomOption {
	return func(a *Atom) {
		a.Set("type", atomType)
	}
}
