package omni

// Property is a basic implementation of the PropertyInterface.
// It represents a named value that can be used to store attributes of atoms.
type Property struct {
	name  string
	value string
}

var _ PropertyInterface = (*Property)(nil)

// NewProperty creates a new Property with the given name and value.
func NewProperty(name, value string) *Property {
	return &Property{
		name:  name,
		value: value,
	}
}

// GetName returns the name of the property.
func (p *Property) GetName() string {
	return p.name
}

// SetName sets the name of the property.
func (p *Property) SetName(name string) {
	p.name = name
}

// GetValue returns the string representation of the property's value.
func (p *Property) GetValue() string {
	return p.value
}

// SetValue sets the property's value from a string.
func (p *Property) SetValue(value string) {
	p.value = value
}
