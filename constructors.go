package omni

import (
	"github.com/gouniverse/uid"
)

// NewAtom creates a new Atom with the given type and applies the provided options.
// If no ID is provided via options, a human-readable UID will be generated.
// Returns an AtomInterface to maintain consistency with other constructors.
func NewAtom(atomType string, opts ...AtomOption) AtomInterface {
	atom := &Atom{
		// Don't generate ID here, it will be handled by the WithID option if provided
		atomType:   atomType,
		properties: make(map[string]string),
		children:   make([]AtomInterface, 0),
	}

	// Apply all options
	for _, opt := range opts {
		opt(atom)
	}

	// If no ID was set by options, generate one
	if atom.id == "" {
		atom.id = uid.HumanUid()
	}

	return atom
}

// NewAtomFromGob creates a new Atom from binary data encoded with the gob package.
// This is a convenience function that delegates to GobToAtom.
//
// Example:
//
//	data, _ := atom.ToGob()
//	newAtom, err := NewAtomFromGob(data)
//
// Parameters:
//   - data: binary data containing the gob-encoded atom
//
// Returns:
//   - AtomInterface: the decoded atom
//   - error: if the data is invalid or cannot be decoded
func NewAtomFromGob(data []byte) (AtomInterface, error) {
	return GobToAtom(data)
}

// NewAtomFromJSON creates a new Atom from a JSON string.
// This is a convenience function that delegates to JSONToAtom.
//
// The JSON should be an object with at least "id" and "type" fields.
// Properties should be in a "properties" map, and children in a "children" array.
//
// Parameters:
//   - jsonStr: JSON string containing the atom data
//
// Returns:
//   - AtomInterface: the parsed atom
//   - error: if the JSON is invalid or missing required fields
//
// Example:
//
//	jsonStr := `{"id":"atom1","type":"test","properties":{"key":"value"}}`
//	atom, err := NewAtomFromJSON(jsonStr)
//
// Note: For parsing multiple atoms from a JSON array, use JSONToAtoms instead.
func NewAtomFromJSON(jsonStr string) (AtomInterface, error) {
	return JSONToAtom(jsonStr)
}

// NewAtomFromMap creates a new Atom from a map.
// This is a convenience function that delegates to MapToAtom.
//
// The map should contain at least "id" and "type" fields.
// Properties should be in a "properties" map, and children in a "children" slice.
//
// Parameters:
//   - atomMap: map containing the atom data
//
// Returns:
//   - AtomInterface: the created atom
//   - error: if the map is missing required fields or is invalid
//
// Example:
//
//	atomMap := map[string]any{
//	  "id":   "atom1",
//	  "type": "test",
//	  "properties": map[string]string{"key": "value"},
//	}
//	atom, err := NewAtomFromMap(atomMap)
func NewAtomFromMap(atomMap map[string]any) (AtomInterface, error) {
	return MapToAtom(atomMap)
}

// AtomOption configures an Atom.
type AtomOption func(*Atom)

// WithID sets the ID of the Atom.
func WithID(id string) AtomOption {
	return func(a *Atom) {
		a.SetID(id)
	}
}

// WithProperties adds properties to the Atom.
// Note: This will not set 'id' or 'type' as they are now direct fields.
func WithProperties(properties map[string]string) AtomOption {
	return func(a *Atom) {
		a.mu.Lock()
		defer a.mu.Unlock()
		for k, v := range properties {
			if k != "id" && k != "type" {
				a.properties[k] = v
			}
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
		a.SetType(atomType)
	}
}

// FromGob decodes an Atom from gob-encoded data.
// This is a helper function that creates a new Atom and calls FromGob on it.
func FromGob(data []byte) (*Atom, error) {
	atom := &Atom{}
	if err := atom.FromGob(data); err != nil {
		return nil, err
	}
	return atom, nil
}
