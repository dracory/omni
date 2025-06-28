package omni

// PropertyInterface defines the contract for any abstract attribute or characteristic.
// Any concrete property type must implement this interface.
type PropertyInterface interface {
	GetName() string
	SetName(name string)
	GetValue() string
	SetValue(value string)
}

// AtomInterface is the universal interface that all composable primitives must satisfy.
// It defines the methods necessary for the system to understand and process any atom,
// regardless of its specific type.
type AtomInterface interface {
	// ID
	GetID() string
	SetID(id string)

	// Type
	GetType() string
	SetType(atomType string)

	// Properties
	GetProperties() []PropertyInterface
	SetProperties(properties []PropertyInterface)

	GetProperty(name string) PropertyInterface
	RemoveProperty(name string)
	SetProperty(property PropertyInterface)

	// Children
	AddChild(a AtomInterface)
	AddChildren(children []AtomInterface)
	GetChildren() []AtomInterface
	SetChildren(children []AtomInterface)

	// Serialization
	ToMap() map[string]any
	ToJson() (string, error)
	ToJsonPretty() (string, error)

	// ToGob encodes the atom to a binary format using the gob package.
	// Returns the binary data and any encoding error.
	ToGob() ([]byte, error)

	// toJsonObject is an internal method for JSON serialization
	toJsonObject() AtomJsonObject
}
