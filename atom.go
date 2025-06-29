package omni

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sync"
)

// Atom is the main implementation of AtomInterface using map[string]string for properties.
// All properties are stored as strings, and children are stored as a slice of AtomInterface.
//
// There are two special properties:
// - "id": the unique identifier of the atom
// - "type": the type of the atom
type Atom struct {
	id         string
	atomType   string
	properties map[string]string
	children   []AtomInterface
	mu         sync.RWMutex
}

// GetID returns the atom's ID.
func (a *Atom) GetID() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.id
}

// SetID sets the atom's ID.
func (a *Atom) SetID(id string) AtomInterface {
	a.mu.Lock()
	a.id = id
	a.mu.Unlock()
	return a
}

// GetType returns the atom's type.
func (a *Atom) GetType() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.atomType
}

// SetType sets the atom's type.
func (a *Atom) SetType(atomType string) AtomInterface {
	a.mu.Lock()
	a.atomType = atomType
	a.mu.Unlock()
	return a
}

// WithData adds initial data to the Atom.
// This is a convenience function for setting multiple key-value pairs at once.
func WithData(data map[string]string) AtomOption {
	return func(a *Atom) {
		a.mu.Lock()
		defer a.mu.Unlock()
		for k, v := range data {
			switch k {
			case "id":
				a.id = v
			case "type":
				a.atomType = v
			default:
				if a.properties == nil {
					a.properties = make(map[string]string)
				}
				a.properties[k] = v
			}
		}
	}
}

// Has checks if the atom has a property with the given key.
func (a *Atom) Has(key string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.properties == nil {
		return false
	}
	_, ok := a.properties[key]
	return ok
}

// Get returns the value for the given key, or "" if not found.
func (a *Atom) Get(key string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.properties == nil {
		return ""
	}
	return a.properties[key]
}

// Remove removes the value for the given key.
func (a *Atom) Remove(key string) AtomInterface {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.properties, key)
	return a
}

// Set sets the value for the given key.
func (a *Atom) Set(key, value string) AtomInterface {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.properties == nil {
		a.properties = make(map[string]string)
	}
	a.properties[key] = value
	return a
}

// GetAll returns all properties of the atom.
func (a *Atom) GetAll() map[string]string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.properties == nil {
		return nil
	}
	props := make(map[string]string, len(a.properties))
	for k, v := range a.properties {
		props[k] = v
	}
	return props
}

// SetAll sets all properties of the atom.
func (a *Atom) SetAll(properties map[string]string) AtomInterface {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.properties = properties
	return a
}

// ChildAdd adds a child atom.
// If child is nil, it's a no-op.
func (a *Atom) ChildAdd(child AtomInterface) AtomInterface {
	if child == nil {
		return a
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.children = append(a.children, child)
	return a
}

// ChildDeleteByID removes a child atom by its ID.
func (a *Atom) ChildDeleteByID(id string) AtomInterface {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, child := range a.children {
		if child.GetID() == id {
			a.children = append(a.children[:i], a.children[i+1:]...)
			break
		}
	}
	return a
}

// ChildrenAdd adds multiple child atoms.
func (a *Atom) ChildrenAdd(children []AtomInterface) AtomInterface {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.children = append(a.children, children...)
	return a
}

// ChildrenGet returns a copy of the children slice.
func (a *Atom) ChildrenGet() []AtomInterface {
	a.mu.RLock()
	defer a.mu.RUnlock()
	children := make([]AtomInterface, len(a.children))
	copy(children, a.children)
	return children
}

// ChildrenLength returns the number of children.
func (a *Atom) ChildrenLength() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.children)
}

// ChildrenSet replaces all children with the given slice.
// Nil children in the input slice will be filtered out.
func (a *Atom) ChildrenSet(children []AtomInterface) AtomInterface {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Filter out nil children
	validChildren := make([]AtomInterface, 0, len(children))
	for _, child := range children {
		if child != nil {
			validChildren = append(validChildren, child)
		}
	}

	a.children = make([]AtomInterface, len(validChildren))
	copy(a.children, validChildren)
	return a
}

// ToMap converts the atom to a map representation with the following structure:
// - id: the atom's ID
// - type: the atom's type
// - properties: a map containing all properties (excluding id and type)
// - children: an array of child atoms
func (a *Atom) ToMap() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Create a copy of properties excluding id and type
	props := make(map[string]string, len(a.properties))
	for k, v := range a.properties {
		if k != "id" && k != "type" {
			props[k] = v
		}
	}

	// Convert children to maps
	children := make([]map[string]interface{}, 0, len(a.children))
	for _, child := range a.children {
		if child != nil {
			children = append(children, child.ToMap())
		}
	}

	// Build the result map
	result := map[string]interface{}{
		"id":       a.id,
		"type":     a.atomType,
		"children": children,
	}

	// Only add properties if not empty
	if len(props) > 0 {
		result["properties"] = props
	}

	return result
}

// ToJSON converts the atom to a JSON string.
func (a *Atom) ToJSON() (string, error) {
	data := a.ToMap()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(jsonData), nil
}

// ToJSONPretty converts the atom to a nicely indented JSON string.
func (a *Atom) ToJSONPretty() (string, error) {
	data := a.ToMap()
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to pretty JSON: %w", err)
	}
	return string(jsonData), nil
}

func (a *Atom) ToGob() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Create a temporary struct for encoding
	temp := struct {
		Properties map[string]string
		Children   []AtomInterface
	}{
		Properties: a.properties,
		Children:   a.children,
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	// Register the types that might be in the interface
	gob.Register(&Atom{})

	// Encode the data
	err := encoder.Encode(temp)
	if err != nil {
		return nil, fmt.Errorf("error encoding atom to gob: %v", err)
	}

	return buf.Bytes(), nil
}

// FromGob decodes a Atom from gob-encoded data.
func FromGob(data []byte) (*Atom, error) {
	var temp struct {
		Properties map[string]string
		Children   []*Atom // Using *Atom here for proper type assertion
	}

	// Register the types that might be in the interface
	gob.Register(&Atom{})

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&temp)
	if err != nil {
		return nil, fmt.Errorf("error decoding atom from gob: %v", err)
	}

	// Create a new atom with the decoded data
	atom := &Atom{
		properties: temp.Properties,
		children:   make([]AtomInterface, len(temp.Children)),
	}

	// Convert children to AtomInterface slice
	for i, child := range temp.Children {
		atom.children[i] = child
	}

	return atom, nil
}

// NewAtomWithData creates a new Atom with the given ID, type, and initial data.
// Deprecated: Use NewAtom with WithID and WithData options instead
func NewAtomWithData(id, atomType string, data map[string]string) *Atom {
	atom := &Atom{
		properties: data,
	}
	return atom.SetID(id).SetType(atomType).(*Atom)
}

// Size returns the approximate memory usage of the atom in bytes.
// Note: This is an approximation and doesn't account for all memory used by the Go runtime.
func (a *Atom) Size() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	size := 0

	// Approximate size of the properties map structure
	size += 8                           // map header
	size += len(a.properties) * (8 + 8) // key and value pointers
	for k, v := range a.properties {
		size += len(k) + len(v)
	}

	// Approximate size of the children slice
	size += 8 + (len(a.children) * 8) // slice header + interface pointers

	// Approximate size of the mutex (3 int64 values)
	size += 24

	return size
}
