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

// FromGob decodes the atom from gob-encoded data.
// This method satisfies the AtomInterface requirement.
func (a *Atom) FromGob(data []byte) error {
	var temp struct {
		ID         string
		Type       string
		Properties map[string]string
		Children   [][]byte
	}

	// Register the type
	gob.Register(&Atom{})

	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&temp); err != nil {
		return fmt.Errorf("error decoding atom from gob: %v", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.id = temp.ID
	a.atomType = temp.Type
	a.properties = temp.Properties

	// Decode children
	a.children = make([]AtomInterface, len(temp.Children))
	for i, childData := range temp.Children {
		child := &Atom{}
		if err := child.FromGob(childData); err != nil {
			return fmt.Errorf("error decoding child %d: %v", i, err)
		}
		a.children[i] = child
	}

	return nil
}

// GobEncode implements the gob.GobEncoder interface.
// This is a wrapper around ToGob for compatibility with the gob package.
func (a *Atom) GobEncode() ([]byte, error) {
	return a.ToGob()
}

// GobDecode implements the gob.GobDecoder interface.
// This is a wrapper around FromGob for compatibility with the gob package.
func (a *Atom) GobDecode(data []byte) error {
	return a.FromGob(data)
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

// ChildFindByID returns the first immediate child with the given ID, or nil if not found.
func (a *Atom) ChildFindByID(id string) AtomInterface {
    a.mu.RLock()
    defer a.mu.RUnlock()
    for _, child := range a.children {
        if child != nil && child.GetID() == id {
            return child
        }
    }
    return nil
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

// ChildrenFindByType returns all immediate children that match the provided type.
func (a *Atom) ChildrenFindByType(atomType string) []AtomInterface {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.children) == 0 {
		return []AtomInterface{}
	}
	result := make([]AtomInterface, 0)
	for _, child := range a.children {
		if child != nil && child.GetType() == atomType {
			result = append(result, child)
		}
	}
	return result
}

// ToGob encodes the atom to a gob-encoded byte slice.
// This is the primary method for gob encoding that satisfies the AtomInterface.
func (a *Atom) ToGob() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// First, encode all children to gob
	childData := make([][]byte, len(a.children))
	for i, child := range a.children {
		if child != nil {
			childBytes, err := child.ToGob()
			if err != nil {
				return nil, fmt.Errorf("error encoding child %d: %v", i, err)
			}
			childData[i] = childBytes
		}
	}

	// Create a temporary struct for encoding with exported fields
	temp := struct {
		ID         string
		Type       string
		Properties map[string]string
		Children   [][]byte
	}{
		ID:         a.id,
		Type:       a.atomType,
		Properties: a.properties,
		Children:   childData,
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	// Register the type
	gob.Register(&Atom{})


	// Encode the data
	if err := encoder.Encode(temp); err != nil {
		return nil, fmt.Errorf("error encoding atom to gob: %v", err)
	}

	return buf.Bytes(), nil
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

// MemoryUsage returns the estimated memory usage of the atom in bytes,
// including all its properties and recursively all its children.
// This is useful for memory profiling and monitoring.
// Note: This is an approximation and doesn't account for all memory used by the Go runtime.
func (a *Atom) MemoryUsage() int {
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
