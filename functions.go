package omni

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
)

// MarshalAtomsToJson converts a slice of AtomInterface to a JSON string.
func MarshalAtomsToJson(atoms []AtomInterface) (string, error) {
	if atoms == nil {
		return "[]", nil
	}

	atomsMap := make([]AtomJsonObject, 0, len(atoms))

	for _, atom := range atoms {
		if atom != nil {
			atomsMap = append(atomsMap, atom.toJsonObject())
		}
	}

	atomsJson, err := json.Marshal(atomsMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal atoms to JSON: %w", err)
	}

	return string(atomsJson), nil
}

// UnmarshalJsonToAtoms converts a JSON string to a slice of AtomInterface.
// The JSON should be an array of atom objects or a single atom object.
func UnmarshalJsonToAtoms(atomsJson string) ([]AtomInterface, error) {
	if atomsJson == "" {
		return []AtomInterface{}, nil
	}

	// Handle the case where the input is a JSON string literal (e.g., """)
	if len(atomsJson) >= 2 && atomsJson[0] == '"' && atomsJson[len(atomsJson)-1] == '"' {
		// This is a JSON string literal, not a JSON object/array
		// Return empty slice for empty string literals
		return []AtomInterface{}, nil
	}

	// First try to unmarshal as an array
	var atomsMap []map[string]any
	err := json.Unmarshal([]byte(atomsJson), &atomsMap)
	if err != nil {
		// If it's not an array, try as a single object
		var singleAtomMap map[string]any
		if singleErr := json.Unmarshal([]byte(atomsJson), &singleAtomMap); singleErr == nil {
			atomsMap = []map[string]any{singleAtomMap}
		} else {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	atoms := make([]AtomInterface, 0, len(atomsMap))

	for _, atomMap := range atomsMap {
		if atomMap == nil {
			continue
		}

		atom, err := ConvertMapToAtom(atomMap)
		if err != nil {
			return nil, fmt.Errorf("failed to convert map to atom: %w", err)
		}

		atoms = append(atoms, atom)
	}

	return atoms, nil
}

// ConvertMapToAtoms converts a slice of atom maps to a slice of AtomInterface.
// This is a convenience function that calls NewAtomFromMap on each element.
// Nil maps in the input will result in nil elements in the output.
// Errors from NewAtomFromMap are ignored, and nil is appended to the result in case of errors.
func ConvertMapToAtoms(atoms []map[string]any) []AtomInterface {
	if atoms == nil {
		return nil
	}

	result := make([]AtomInterface, 0, len(atoms))

	for _, atom := range atoms {
		if atom == nil {
			result = append(result, nil)
			continue
		}

		// Ignore errors from NewAtomFromMap to maintain backward compatibility
		atomObj, _ := NewAtomFromMap(atom)
		result = append(result, atomObj)
	}

	return result
}

// ConvertAtomsToMap converts a slice of AtomInterface to a slice of maps.
// This is a convenience function that calls ToMap() on each atom.
// Nil atoms in the input will be skipped in the output.
func ConvertAtomsToMap(atoms []AtomInterface) []map[string]any {
	if atoms == nil {
		return nil
	}

	result := make([]map[string]any, 0, len(atoms))

	for _, atom := range atoms {
		if atom == nil {
			continue
		}

		if atomMap := atom.ToMap(); atomMap != nil {
			result = append(result, atomMap)
		}
	}

	return result
}

// ConvertMapToAtom converts a map to an AtomInterface.
//
// The map must represent a valid atom with at least "id" and "type" fields.
// Properties should be in a "parameters" map, and children in a "children" slice.
//
// Parameters:
//   - atomMap: map containing the atom data
//
// Returns:
//   - AtomInterface: the converted atom
//   - error: if the map is not a valid atom
func ConvertMapToAtom(atomMap map[string]any) (AtomInterface, error) {
	if atomMap == nil {
		return nil, errors.New("atom map cannot be nil")
	}

	// Make a copy of the map to avoid modifying the original
	atomMapCopy := make(map[string]any, len(atomMap))
	for k, v := range atomMap {
		atomMapCopy[k] = v
	}

	// Ensure the map has the required fields
	if _, ok := atomMapCopy["id"].(string); !ok {
		return nil, errors.New("atom map must contain a string 'id' field")
	}

	if _, ok := atomMapCopy["type"].(string); !ok {
		return nil, errors.New("atom map must contain a string 'type' field")
	}

	// Ensure parameters is a map
	if _, ok := atomMapCopy["parameters"].(map[string]any); !ok {
		atomMapCopy["parameters"] = make(map[string]any)
	}

	// Ensure children is a slice
	if _, ok := atomMapCopy["children"].([]any); !ok {
		atomMapCopy["children"] = make([]any, 0)
	}

	// Create the atom from the map copy
	atom, err := NewAtomFromMap(atomMapCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to create atom from map: %w", err)
	}

	return atom, nil
}

// mapToAtomMap is an internal function that validates and normalizes an atom map.
// It ensures all required fields are present and have the correct types.
//
// Parameters:
//   - atomMap: the map to validate and normalize
//
// Returns:
//   - map[string]any: the normalized atom map
//   - error: if the map is not a valid atom
//
// FromGob decodes an Atom from binary data encoded with the gob package.
// This is a standalone function since it needs to create a new Atom instance.
func FromGob(data []byte) (AtomInterface, error) {
	if len(data) == 0 {
		return nil, errors.New("cannot decode empty data")
	}

	// Create a temporary struct for decoding
	var temp struct {
		ID         string
		Type       string
		Properties map[string]string
		Children   [][]byte
	}

	// Decode the data into the temporary struct
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&temp); err != nil {
		return nil, fmt.Errorf("gob decode failed: %w", err)
	}

	// Create a new atom
	atom := NewAtom(temp.Type, WithID(temp.ID))

	// Convert properties
	for name, value := range temp.Properties {
		atom.SetProperty(NewProperty(name, value))
	}

	// Recursively decode children
	for _, childData := range temp.Children {
		child, err := FromGob(childData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode child: %w", err)
		}
		atom.AddChild(child)
	}

	return atom, nil
}
