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
	if valid, err := isValidAtomJSON(atomsJson); !valid {
		if err != nil {
			return nil, fmt.Errorf("invalid atom JSON: %v", err)
		}
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

	// Validate the map structure first
	if valid, err := isValidAtomMap(atomMap); !valid {
		return nil, fmt.Errorf("invalid atom map: %v", err)
	}

	// Make a copy of the map to avoid modifying the original
	atomMapCopy := make(map[string]any, len(atomMap))
	for k, v := range atomMap {
		atomMapCopy[k] = v
	}

	// Ensure parameters is a map if not present
	if _, ok := atomMapCopy["parameters"].(map[string]any); !ok {
		atomMapCopy["parameters"] = make(map[string]any)
	}

	// Ensure children is a slice if not present
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

		// Validate the atom map first
		if valid, _ := isValidAtomMap(atom); !valid {
			// For backward compatibility, we ignore invalid atoms
			// Consider changing this to return an error in a future major version
			continue
		}

		// Ignore errors from NewAtomFromMap to maintain backward compatibility
		// We already validated the map structure, so this should not fail
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
// - Validates the input data structure
// - Decodes the count of atoms first
// - Then decodes each atom's data and validates it before conversion
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

	// Validate count is reasonable
	if count < 0 {
		return nil, fmt.Errorf("invalid atom count: %d", count)
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

		// Validate data length is reasonable
		if dataLen < 0 || dataLen > 10*1024*1024 { // 10MB max per atom
			return nil, fmt.Errorf("invalid data length %d for atom %d", dataLen, i)
		}

		// Read the atom data
		atomData := make([]byte, dataLen)
		if _, err := io.ReadFull(buffer, atomData); err != nil {
			return nil, fmt.Errorf("failed to read atom %d data: %w", i, err)
		}

		// Validate the atom data before creating the atom
		if valid, err := isValidGobData(atomData); !valid {
			return nil, fmt.Errorf("invalid gob data for atom %d: %v", i, err)
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

// isValidGobData validates that the given data is a valid gob-encoded atom.
//
// Business logic:
// - Checks for empty or nil input
// - Attempts to decode the data into a temporary struct
// - Validates the presence of required fields
// - Recursively validates child atoms
//
// Parameters:
//   - data: binary data to validate
//
// Returns:
//   - bool: true if the data is valid
//   - error: description of the validation failure if invalid
func isValidGobData(data []byte) (bool, error) {
	if len(data) == 0 {
		return false, errors.New("cannot validate empty data")
	}

	// Temporary struct for validation
	var temp struct {
		ID         string
		Type       string
		Properties map[string]string
		Children   [][]byte
	}

	// Try to decode the data
	r := bytes.NewReader(data)
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(&temp); err != nil {
		return false, fmt.Errorf("invalid gob data: %w", err)
	}

	// Validate required fields
	if temp.Type == "" {
		return false, errors.New("missing required field: Type")
	}

	if temp.ID == "" {
		return false, errors.New("missing required field: ID")
	}

	// Validate properties if present
	if temp.Properties == nil {
		return false, errors.New("properties map cannot be nil")
	}

	// Recursively validate children
	for i, childData := range temp.Children {
		if valid, err := isValidGobData(childData); !valid {
			return false, fmt.Errorf("invalid child at index %d: %v", i, err)
		}
	}

	return true, nil
}

// GobToAtom decodes an Atom from binary data encoded with the gob package.
// This is a standalone function that creates a new Atom instance from gob-encoded data.
//
// Business logic:
// - Handles empty or nil input by returning an error
// - Validates the gob data before decoding
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

	// Validate the gob data first
	if valid, err := isValidGobData(data); !valid {
		return nil, fmt.Errorf("invalid gob data: %v", err)
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
// - true, nil if the JSON string is a valid Atom JSON string
// - false, error with details if the JSON is invalid
func isValidAtomJSON(jsonString string) (bool, error) {
	if jsonString == "" || jsonString == "null" {
		return false, errors.New("JSON string cannot be empty or 'null'")
	}

	// Check basic JSON structure
	isObject := strings.HasPrefix(jsonString, "{") && strings.HasSuffix(jsonString, "}")
	isArray := strings.HasPrefix(jsonString, "[") && strings.HasSuffix(jsonString, "]")

	if !isObject && !isArray {
		return false, errors.New("JSON must be an object or array")
	}

	// If it's an array, it's valid if it's empty or contains valid objects
	if isArray {
		// Empty array is valid
		if jsonString == "[]" {
			return true, nil
		}
		// For non-empty arrays, we'll validate each element during parsing
		return true, nil
	}

	// For single objects, check required fields
	hasID := strings.Contains(jsonString, `"id"`)
	hasType := strings.Contains(jsonString, `"type"`)

	if !hasID || !hasType {
		missing := []string{}
		if !hasID {
			missing = append(missing, "id")
		}
		if !hasType {
			missing = append(missing, "type")
		}
		return false, fmt.Errorf("missing required fields: %v", strings.Join(missing, ", "))
	}

	return true, nil
}

// isValidAtomMap validates that a map represents a valid atom structure.
//
// Business logic:
// - Checks for required fields (id, type)
// - Validates that parameters is a map if present
// - Validates that children is a slice if present
// - Validates that all children are valid atom maps
//
// Parameters:
//   - atomMap: the map to validate
//
// Returns:
//   - bool: true if the map is a valid atom structure
//   - error: description of the validation failure if invalid
func isValidAtomMap(atomMap map[string]any) (bool, error) {
	if atomMap == nil {
		return false, errors.New("atom map cannot be nil")
	}

	// Check required fields
	id, idOk := atomMap["id"].(string)
	if !idOk || id == "" {
		return false, errors.New("atom map must contain a non-empty string 'id' field")
	}

	typeStr, typeOk := atomMap["type"].(string)
	if !typeOk || typeStr == "" {
		return false, errors.New("atom map must contain a non-empty string 'type' field")
	}

	// Validate parameters if present
	if params, ok := atomMap["parameters"]; ok && params != nil {
		if _, ok := params.(map[string]any); !ok {
			return false, errors.New("atom parameters must be a map[string]any")
		}
	}

	// Validate children if present
	if children, ok := atomMap["children"]; ok && children != nil {
		childrenSlice, ok := children.([]any)
		if !ok {
			return false, errors.New("atom children must be a slice")
		}

		for i, child := range childrenSlice {
			childMap, ok := child.(map[string]any)
			if !ok {
				return false, fmt.Errorf("child at index %d is not a valid atom map", i)
			}

			if valid, err := isValidAtomMap(childMap); !valid {
				return false, fmt.Errorf("invalid child at index %d: %v", i, err)
			}
		}
	}

	return true, nil
}
