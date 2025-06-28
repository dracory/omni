package omni

import (
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

		if atomObj := NewAtomFromMap(atom); atomObj != nil {
			result = append(result, atomObj)
		} else {
			result = append(result, nil)
		}
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

	// Check required fields
	if _, ok := atomMap["id"]; !ok {
		return nil, errors.New("atom map is missing required field 'id'")
	}

	if _, ok := atomMap["type"]; !ok {
		return nil, errors.New("atom map is missing required field 'type'")
	}

	// Create a deep copy to avoid modifying the input
	atomMapCopy := make(map[string]any, len(atomMap))
	for k, v := range atomMap {
		atomMapCopy[k] = v
	}

	// Process children if they exist
	if children, ok := atomMapCopy["children"].([]any); ok {
		processedChildren := make([]any, 0, len(children))
		for _, child := range children {
			if childMap, ok := child.(map[string]any); ok {
				if childMap != nil {
					processedChildren = append(processedChildren, childMap)
				}
			}
		}
		atomMapCopy["children"] = processedChildren
	}

	return NewAtomFromMap(atomMapCopy), nil
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
func mapToAtomMap(atomMap map[string]any) (map[string]any, error) {
	if atomMap == nil {
		return nil, errors.New("atom map cannot be nil")
	}

	// Check required fields
	id, ok := atomMap["id"].(string)
	if !ok || id == "" {
		return nil, errors.New("atom map must have a non-empty string 'id' field")
	}

	typeStr, ok := atomMap["type"].(string)
	if !ok || typeStr == "" {
		return nil, errors.New("atom map must have a non-empty string 'type' field")
	}

	// Process parameters
	params := make(map[string]string)
	if paramsAny, ok := atomMap["parameters"]; ok {
		if paramsMap, ok := paramsAny.(map[string]any); ok {
			for k, v := range paramsMap {
				if strVal, ok := v.(string); ok {
					params[k] = strVal
				} else {
					// Convert non-string values to string
					params[k] = fmt.Sprintf("%v", v)
				}
			}
		}
	}

	// Process children
	var children []map[string]any
	if childrenAny, ok := atomMap["children"]; ok {
		if childrenSlice, ok := childrenAny.([]any); ok {
			children = make([]map[string]any, 0, len(childrenSlice))
			for _, childAny := range childrenSlice {
				if childMap, ok := childAny.(map[string]any); ok && childMap != nil {
					// Recursively process child atoms
					if child, err := mapToAtomMap(childMap); err == nil {
						children = append(children, child)
					}
				}
			}
		}
	}

	// Create a new map with normalized structure
	result := make(map[string]any)
	result["id"] = id
	result["type"] = typeStr
	result["parameters"] = params
	result["children"] = children

	return result, nil
}
