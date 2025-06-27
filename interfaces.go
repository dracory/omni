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
	GetID() string
	SetID(id string)

	GetType() string
	SetType(atomType string)

	GetProperties() []PropertyInterface
	SetProperties(properties []PropertyInterface)

	GetProperty(name string) PropertyInterface
	SetProperty(property PropertyInterface)
	RemoveProperty(name string)

	GetChildren() []AtomInterface
	AddChild(a AtomInterface)
}
