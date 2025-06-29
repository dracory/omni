package omni

// AtomInterface is the universal interface that all composable primitives must satisfy.
// It defines the methods necessary for the system to understand and process any atom,
// regardless of its specific type.
type AtomInterface interface {
	// ID returns the unique identifier of the atom
	GetID() string
	SetID(id string) AtomInterface

	// Type returns the type of the atom
	GetType() string
	SetType(atomType string) AtomInterface

	// Property access
	Get(key string) string
	Has(key string) bool
	Remove(key string) AtomInterface
	Set(key, value string) AtomInterface

	GetAll() map[string]string
	SetAll(properties map[string]string) AtomInterface

	// Children management
	ChildAdd(child AtomInterface) AtomInterface
	ChildDeleteByID(id string) AtomInterface

	ChildrenAdd(children []AtomInterface) AtomInterface
	ChildrenGet() []AtomInterface
	ChildrenSet(children []AtomInterface) AtomInterface

	ChildrenLength() int

	// Serialization
	ToMap() map[string]any
	ToJSON() (string, error)
	ToJSONPretty() (string, error)
	ToGob() ([]byte, error)

	// MemoryUsage returns the estimated memory usage of the atom in bytes,
	// including all its properties and recursively all its children.
	// This is useful for memory profiling and monitoring.
	MemoryUsage() int
}
