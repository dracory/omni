package omni

import "sync"

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

// NewAtom creates a new Atom with the given ID and type.
func NewAtom(id, atomType string) *Atom {
	return &Atom{
		id:         id,
		atomType:   atomType,
		properties: make([]PropertyInterface, 0),
		children:   make([]AtomInterface, 0),
	}
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

// GetChildren returns all child atoms.
func (a *Atom) GetChildren() []AtomInterface {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Return a copy to prevent external modification of our internal slice
	children := make([]AtomInterface, len(a.children))
	copy(children, a.children)
	return children
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
