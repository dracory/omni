package omni

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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
func NewAtomFromGob(data []byte) (*Atom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data provided")
	}

	// Define the same struct used in ToGob
	type gobAtom struct {
		ID         string
		Type       string
		Properties map[string]string
		Children   [][]byte
	}

	// Decode the data into the temporary struct
	var temp gobAtom
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&temp); err != nil {
		return nil, fmt.Errorf("gob decode failed: %w", err)
	}

	// Create a new atom with the decoded data
	atom := NewAtom(temp.Type, WithID(temp.ID))

	// Set properties
	for name, value := range temp.Properties {
		atom.SetProperty(NewProperty(name, value))
	}

	// Recursively decode children
	for _, childData := range temp.Children {
		child, err := NewAtomFromGob(childData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode child: %w", err)
		}
		if child != nil {
			atom.AddChild(child)
		}
	}

	return atom, nil
}

// NewAtomFromJSON creates a new Atom from a JSON string.
// The JSON should be an object with at least "id" and "type" fields.
// Properties should be in a "parameters" map, and children in a "children" array.
// Returns the Atom and nil error on success, or nil and an error if the JSON is invalid or cannot be unmarshaled.
func NewAtomFromJSON(jsonStr string) (*Atom, error) {
	if jsonStr == "" {
		return nil, errors.New("empty JSON string provided")
	}

	var atomMap map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &atomMap); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	atom, err := NewAtomFromMap(atomMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create atom from map: %w", err)
	}
	if atom == nil {
		return nil, errors.New("unexpected nil atom from NewAtomFromMap")
	}
	return atom, nil
}

// NewAtomFromMap creates a new Atom from a map representation.
// The map must contain at least "id" and "type" fields.
// Properties should be in a "parameters" map, and children in a "children" slice.
// Returns the Atom and nil error on success, or nil and an error if the map is invalid.
func NewAtomFromMap(atomMap map[string]any) (*Atom, error) {
	if atomMap == nil {
		return nil, errors.New("atom map cannot be nil")
	}

	id, ok := atomMap["id"].(string)
	if !ok || id == "" {
		return nil, errors.New("atom map must contain a non-empty 'id' field")
	}

	typeStr, ok := atomMap["type"].(string)
	if !ok || typeStr == "" {
		return nil, errors.New("atom map must contain a non-empty 'type' field")
	}

	atom := NewAtom(typeStr, WithID(id))

	// Set properties - check both 'properties' and 'parameters' for backward compatibility
	var props map[string]any
	if p, ok := atomMap["properties"].(map[string]any); ok {
		props = p
	} else if p, ok := atomMap["parameters"].(map[string]any); ok {
		props = p
	}

	if props != nil {
		for k, v := range props {
			if value, ok := v.(string); ok {
				atom.SetProperty(NewProperty(k, value))
			}
		}
	}

	// Handle children
	if children, ok := atomMap["children"].([]any); ok {
		for _, child := range children {
			if childMap, ok := child.(map[string]any); ok {
				if childAtom, err := NewAtomFromMap(childMap); err == nil && childAtom != nil {
					atom.AddChild(childAtom)
				}
				// Silently skip children that fail to parse
			}
		}
	}

	return atom, nil
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
