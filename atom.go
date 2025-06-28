package omni

import (
	"encoding/json"
	"sync"

	"github.com/gouniverse/uid"
)

// Atom is a basic implementation of the AtomInterface.
// It represents a composable primitive that can have properties and child atoms.
type Atom struct {
	id         string
	atomType   string
	properties []PropertyInterface
	children   []AtomInterface
	mu         sync.RWMutex // Protects concurrent access to properties and children
}

var _ AtomInterface = (*Atom)(nil)

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

// NewAtomFromMap creates a new Atom from a map representation.
// The map must contain at least "id" and "type" fields.
// Properties should be in a "parameters" map, and children in a "children" slice.
func NewAtomFromMap(atomMap map[string]any) *Atom {
	if atomMap == nil {
		return nil
	}

	id, _ := atomMap["id"].(string)
	typeStr, _ := atomMap["type"].(string)

	atom := NewAtom(typeStr, WithID(id))

	// Set properties
	if params, ok := atomMap["parameters"].(map[string]any); ok {
		for k, v := range params {
			if value, ok := v.(string); ok {
				atom.SetProperty(NewProperty(k, value))
			}
		}
	}

	// Handle children
	if children, ok := atomMap["children"].([]any); ok {
		for _, child := range children {
			if childMap, ok := child.(map[string]any); ok {
				if childAtom := NewAtomFromMap(childMap); childAtom != nil {
					atom.AddChild(childAtom)
				}
			}
		}
	}

	return atom
}

// GetID returns the unique identifier of the atom.
func (a *Atom) GetID() string {
	return a.id
}

// SetID sets the unique identifier of the atom.
func (a *Atom) SetID(id string) {
	a.id = id
}

// GetType returns the type of the atom.
func (a *Atom) GetType() string {
	return a.atomType
}

// SetType sets the type of the atom.
func (a *Atom) SetType(atomType string) {
	a.atomType = atomType
}

// GetProperties returns all properties of the atom.
func (a *Atom) GetProperties() []PropertyInterface {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Return a copy to prevent external modification of our internal slice
	props := make([]PropertyInterface, len(a.properties))
	copy(props, a.properties)
	return props
}

// SetProperties sets all properties of the atom at once.
func (a *Atom) SetProperties(properties []PropertyInterface) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Create a new slice to avoid external modification
	a.properties = make([]PropertyInterface, len(properties))
	copy(a.properties, properties)
}

// GetProperty returns a specific property by name, or nil if not found.
func (a *Atom) GetProperty(name string) PropertyInterface {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, prop := range a.properties {
		if prop.GetName() == name {
			return prop
		}
	}
	return nil
}

// SetProperty adds or updates a property.
// If the property is nil, it will be ignored.
func (a *Atom) SetProperty(property PropertyInterface) {
	if property == nil {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	for i, prop := range a.properties {
		if prop.GetName() == property.GetName() {
			a.properties[i] = property
			return
		}
	}
	a.properties = append(a.properties, property)
}

// RemoveProperty removes a property by name.
func (a *Atom) RemoveProperty(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i, prop := range a.properties {
		if prop.GetName() == name {
			a.properties = append(a.properties[:i], a.properties[i+1:]...)
			return
		}
	}
}

// AddChild adds a new child atom.
// If the child is nil, it will be ignored.
func (a *Atom) AddChild(child AtomInterface) {
	if child == nil {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.children = append(a.children, child)
}

// AddChildren adds multiple child atoms at once.
// If any of the children are nil, they will be ignored.
func (a *Atom) AddChildren(children []AtomInterface) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, child := range children {
		if child == nil {
			continue
		}
		a.children = append(a.children, child)
	}
}

// GetChildren returns all child atoms.
func (a *Atom) GetChildren() []AtomInterface {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Return a copy to prevent external modification of our internal slice
	children := make([]AtomInterface, len(a.children))
	copy(children, a.children)
	return children
}

// SetChildren sets all child atoms at once.
// If children is nil or contains nil values, they will be filtered out.
func (a *Atom) SetChildren(children []AtomInterface) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if children == nil {
		a.children = []AtomInterface{}
		return
	}

	// Filter out nil children
	a.children = make([]AtomInterface, 0, len(children))
	for _, child := range children {
		if child != nil {
			a.children = append(a.children, child)
		}
	}
}

// AtomJsonObject represents the JSON structure of an Atom
type AtomJsonObject struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Parameters map[string]string `json:"parameters"`
	Children   []AtomJsonObject  `json:"children"`
}

// ToMap converts the Atom to a map representation.
func (a *Atom) ToMap() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	childrenMap := make([]map[string]interface{}, 0, len(a.children))
	for _, child := range a.children {
		if child != nil {
			childMap := map[string]interface{}{
				"id":   child.GetID(),
				"type": child.GetType(),
			}
			childrenMap = append(childrenMap, childMap)
		}
	}

	// Convert properties to a map
	properties := make(map[string]string, len(a.properties))
	for _, prop := range a.properties {
		if prop != nil {
			properties[prop.GetName()] = prop.GetValue()
		}
	}

	return map[string]interface{}{
		"id":         a.id,
		"type":       a.atomType,
		"parameters": properties,
		"children":   childrenMap,
	}
}

// ToJson converts the Atom to a JSON string.
func (a *Atom) ToJson() (string, error) {
	jsonObject := a.toJsonObject()
	jsonBytes, err := json.Marshal(jsonObject)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ToJsonPretty converts the Atom to a nicely indented JSON string.
func (a *Atom) ToJsonPretty() (string, error) {
	jsonObject := a.toJsonObject()
	jsonBytes, err := json.MarshalIndent(jsonObject, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// toJsonObject converts the Atom to its JSON object representation.
func (a *Atom) toJsonObject() AtomJsonObject {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Convert properties to a map
	parameters := make(map[string]string, len(a.properties))
	for _, prop := range a.properties {
		if prop != nil {
			parameters[prop.GetName()] = prop.GetValue()
		}
	}

	// Convert children to JSON objects
	children := make([]AtomJsonObject, 0, len(a.children))
	for _, child := range a.children {
		if child != nil {
			if atom, ok := child.(*Atom); ok {
				children = append(children, atom.toJsonObject())
			}
		}
	}

	return AtomJsonObject{
		ID:         a.id,
		Type:       a.atomType,
		Parameters: parameters,
		Children:   children,
	}
}
