package omni

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// AtomsToJSON converts a slice of AtomInterface to a JSON string.
//
// Business logic:
// - Handles nil input by returning an empty array JSON string
// - Converts each non-nil atom to a JSON object using toJsonObject()
// - Returns a JSON array of atom objects
//
// Returns:
// - string: JSON-encoded array of atoms
// - error: if marshaling to JSON fails
func AtomsToJSON(atoms []AtomInterface) (string, error) {
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

// JSONToAtoms converts a JSON string to a slice of AtomInterface.
//
// Business logic:
// - Handles empty input by returning an empty slice
// - Supports both array of atoms and single atom object
// - Validates JSON structure before processing
// - Converts each JSON object to an Atom using MapToAtom
//
// Parameters:
//   - atomsJson: JSON string containing atom data
//
// Returns:
//   - []AtomInterface: slice of parsed atoms
//   - error: if JSON is invalid or missing required fields
func JSONToAtoms(atomsJson string) ([]AtomInterface, error) {
	// Early return for empty input
	if atomsJson == "" {
		return []AtomInterface{}, nil
	}

	// Handle the case where the input is a JSON string literal (e.g., """)
	if len(atomsJson) >= 2 && atomsJson[0] == '"' && atomsJson[len(atomsJson)-1] == '"' {
		// This is a JSON string literal, not a JSON object/array
		// Return empty slice for empty string literals
		return []AtomInterface{}, nil
	}

	// Validate JSON structure before processing
	if !isValidAtomJSON(atomsJson) {
		return nil, errors.New("invalid atom JSON: missing required fields or malformed JSON")
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

		atom, err := MapToAtom(atomMap)
		if err != nil {
			return nil, fmt.Errorf("failed to convert map to atom: %w", err)
		}

		atoms = append(atoms, atom)
	}

	return atoms, nil
}

// AtomsToMap converts a slice of AtomInterface to a slice of maps.
//
// Business logic:
// - Handles nil input by returning nil
// - Skips nil atoms in the input
// - Converts each atom to a map using ToMap()
// - Only includes non-nil results in the output
//
// Parameters:
//   - atoms: slice of AtomInterface to convert
//
// Returns:
//   - []map[string]any: slice of atom maps (never contains nils)
func AtomsToMap(atoms []AtomInterface) []map[string]any {
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

// MapToAtom converts a map to an AtomInterface.
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
func MapToAtom(atomMap map[string]any) (AtomInterface, error) {
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

	// Create a new atom with the required fields
	atomType, _ := atomMapCopy["type"].(string)
	atom := NewAtom(atomType, WithID(atomMapCopy["id"].(string)))

	// Set properties if they exist
	if params, ok := atomMapCopy["parameters"].(map[string]any); ok {
		for key, value := range params {
			if strVal, ok := value.(string); ok {
				atom.SetProperty(NewProperty(key, strVal))
			}
		}
	}

	// Handle children if they exist
	if children, ok := atomMapCopy["children"].([]any); ok {
		for _, child := range children {
			if childMap, ok := child.(map[string]any); ok {
				childAtom, err := MapToAtom(childMap)
				if err != nil {
					return nil, fmt.Errorf("failed to create child atom: %w", err)
				}
				atom.AddChild(childAtom)
			}
		}
	}

	return atom, nil
}

// MapToAtoms converts a slice of atom maps to a slice of AtomInterface.
//
// Business logic:
// - Handles nil input by returning nil
// - Preserves nil elements in the output for nil inputs
// - Silently ignores errors from NewAtomFromMap (for backward compatibility)
// - Maintains the order of atoms from input to output
//
// Note: This function is provided for convenience and backward compatibility.
// For better error handling, consider using NewAtomFromMap directly.
//
// Parameters:
//   - atoms: slice of maps containing atom data
//
// Returns:
//   - []AtomInterface: slice of converted atoms (may contain nils)
func MapToAtoms(atoms []map[string]any) []AtomInterface {
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

// AtomsToGob encodes a slice of AtomInterface to binary data using the gob package.
// It encodes each atom using its ToGob method and collects the results.
//
// Business logic:
// - Handles nil or empty input by returning an empty byte slice
// - Encodes each atom using its ToGob method
// - Collects the encoded data for all atoms
//
// Parameters:
//   - atoms: slice of AtomInterface to encode
//
// Returns:
//   - []byte: gob-encoded binary data
//   - error: if encoding fails
func AtomsToGob(atoms []AtomInterface) ([]byte, error) {
	if len(atoms) == 0 {
		return []byte{}, nil
	}

	// Encode each atom to gob
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	
	// First encode the number of atoms
	if err := encoder.Encode(len(atoms)); err != nil {
		return nil, fmt.Errorf("failed to encode atom count: %w", err)
	}

	// Then encode each atom's gob data
	for _, atom := range atoms {
		if atom == nil {
			// Encode a nil marker
			if err := encoder.Encode(false); err != nil {
				return nil, fmt.Errorf("failed to encode nil marker: %w", err)
			}
			continue
		}

		// Encode a non-nil marker
		if err := encoder.Encode(true); err != nil {
			return nil, fmt.Errorf("failed to encode non-nil marker: %w", err)
		}

		// Encode the atom's gob data
		atomData, err := atom.ToGob()
		if err != nil {
			return nil, fmt.Errorf("failed to encode atom to gob: %w", err)
		}

		// Encode the length and then the data
		if err := encoder.Encode(len(atomData)); err != nil {
			return nil, fmt.Errorf("failed to encode atom data length: %w", err)
		}
		if _, err := buf.Write(atomData); err != nil {
			return nil, fmt.Errorf("failed to write atom data: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// GobToAtoms decodes multiple atoms from binary data encoded with the gob package.
// It decodes the data in the format written by AtomsToGob.
//
// Business logic:
// - Handles empty or nil input by returning an empty slice
// - Decodes the count of atoms first
// - Then decodes each atom's data and converts it to an Atom using NewAtomFromGob
// - Preserves the order of atoms from the encoded data
//
// Parameters:
//   - data: binary data containing gob-encoded atoms
//
// Returns:
//   - []AtomInterface: slice of decoded atoms
//   - error: if the data cannot be decoded or is invalid
func GobToAtoms(data []byte) ([]AtomInterface, error) {
	if len(data) == 0 {
		return []AtomInterface{}, nil
	}

	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	// First decode the number of atoms
	var count int
	if err := decoder.Decode(&count); err != nil {
		return nil, fmt.Errorf("failed to decode atom count: %w", err)
	}

	result := make([]AtomInterface, 0, count)

	// Then decode each atom
	for i := 0; i < count; i++ {
		// Decode the nil marker
		var isPresent bool
		if err := decoder.Decode(&isPresent); err != nil {
			return nil, fmt.Errorf("failed to decode nil marker for atom %d: %w", i, err)
		}

		if !isPresent {
			result = append(result, nil)
			continue
		}

		// Decode the atom data length
		var dataLen int
		if err := decoder.Decode(&dataLen); err != nil {
			return nil, fmt.Errorf("failed to decode data length for atom %d: %w", i, err)
		}

		// Read the atom data
		atomData := make([]byte, dataLen)
		if _, err := io.ReadFull(buffer, atomData); err != nil {
			return nil, fmt.Errorf("failed to read atom %d data: %w", i, err)
		}

		// Create the atom from the gob data
		atom, err := NewAtomFromGob(atomData)
		if err != nil {
			return nil, fmt.Errorf("failed to create atom %d from gob: %w", i, err)
		}

		result = append(result, atom)
	}

	return result, nil
}

// GobToAtom decodes an Atom from binary data encoded with the gob package.
// This is a standalone function that creates a new Atom instance from gob-encoded data.
//
// Business logic:
// - Handles empty or nil input by returning an error
// - Decodes a single atom from the gob-encoded data
// - Validates the decoded atom
//
// Parameters:
//   - data: binary data containing a gob-encoded atom
//
// Returns:
//   - AtomInterface: the decoded atom
//   - error: if the data cannot be decoded or is invalid
func GobToAtom(data []byte) (AtomInterface, error) {
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
		child, err := GobToAtom(childData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode child: %w", err)
		}
		atom.AddChild(child)
	}

	return atom, nil
}

// isValidAtomJSON checks if the passed JSON string is a valid Atom JSON object or array of objects
//
// Business logic:
// - checks if the JSON string is empty or "null" (returns false)
// - checks basic JSON structure (starts with '{' or '[' and ends with '}' or ']')
// - for objects, checks if it contains required "id" and "type" fields
// - for arrays, checks if it's empty or contains valid objects
//
// Returns:
// - true if the JSON string is a valid Atom JSON string
// - false otherwise
func isValidAtomJSON(jsonString string) bool {
	if jsonString == "" || jsonString == "null" {
		return false
	}

	// Check basic JSON structure
	if !((strings.HasPrefix(jsonString, "{") && strings.HasSuffix(jsonString, "}")) ||
		(strings.HasPrefix(jsonString, "[") && strings.HasSuffix(jsonString, "]"))) {
		return false
	}

	// If it's an array, it's valid if it's empty or contains valid objects
	if strings.HasPrefix(jsonString, "[") {
		// Empty array is valid
		if jsonString == "[]" {
			return true
		}
		// For non-empty arrays, we'll validate each element during parsing
		return true
	}

	// For single objects, check required fields
	hasID := strings.Contains(jsonString, `"id"`)
	hasType := strings.Contains(jsonString, `"type"`)

	return hasID && hasType
}
