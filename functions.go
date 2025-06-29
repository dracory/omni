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
// - Converts each non-nil atom to a JSON object using ToMap()
// - Returns a JSON array of atom objects
//
// Returns:
// - string: JSON-encoded array of atoms
// - error: if marshaling to JSON fails
func AtomsToJSON(atoms []AtomInterface) (string, error) {
	if atoms == nil {
		return "[]", nil
	}

	atomsMaps := make([]map[string]any, 0, len(atoms))

	for _, atom := range atoms {
		if atom != nil {
			atomsMaps = append(atomsMaps, atom.ToMap())
		}
	}

	atomsJSON, err := json.Marshal(atomsMaps)
	if err != nil {
		return "", fmt.Errorf("failed to marshal atoms to JSON: %w", err)
	}

	return string(atomsJSON), nil
}

// JSONToAtom converts a JSON string to a single Atom.
//
// Business logic:
// - Handles empty input by returning an error
// - Validates JSON structure before processing
// - Converts JSON object to an Atom using MapToAtom
//
// Parameters:
//   - jsonStr: JSON string containing a single atom's data
//
// Returns:
//   - AtomInterface: the parsed atom
//   - error: if JSON is invalid or missing required fields
func JSONToAtom(jsonStr string) (AtomInterface, error) {
	if jsonStr == "" {
		return nil, errors.New("empty JSON string provided")
	}

	// First try to parse as a single atom
	var atomMap map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &atomMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Convert map to atom
	atom, err := MapToAtom(atomMap)
	if err != nil {
		return nil, fmt.Errorf("invalid atom data: %w", err)
	}

	return atom, nil
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
// Properties should be in a nested "properties" map, and children in a "children" slice.
// For backward compatibility, top-level properties are also supported but not recommended.
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

	// Extract type and id from top-level fields
	atomType, _ := atomMapCopy["type"].(string)
	id, _ := atomMapCopy["id"].(string)

	// Create a new atom with the required fields
	var atom *Atom
	if id != "" && atomType != "" {
		atom = NewAtom(atomType, WithID(id))
	} else if atomType != "" {
		atom = NewAtom(atomType)
	} else {
		return nil, errors.New("missing required 'type' field in atom map")
	}

	// Process properties from the nested properties map if it exists
	props := make(map[string]string)
	if propsMap, ok := atomMapCopy["properties"].(map[string]any); ok && len(propsMap) > 0 {
		for k, v := range propsMap {
			// Convert value to string
			var strVal string
			switch v := v.(type) {
			case string:
				strVal = v
			case fmt.Stringer:
				strVal = v.String()
			default:
				strVal = fmt.Sprintf("%v", v)
			}
			props[k] = strVal
		}
	}

	// For backward compatibility, also check for top-level properties
	for k, v := range atomMapCopy {
		if k != "id" && k != "type" && k != "properties" && k != "children" {
			var strVal string
			switch v := v.(type) {
			case string:
				strVal = v
			case fmt.Stringer:
				strVal = v.String()
			default:
				strVal = fmt.Sprintf("%v", v)
			}
			props[k] = strVal
		}
	}

	// Set all properties
	if len(props) > 0 {
		atom.SetAll(props)
	}

	// Handle children if they exist
	if children, ok := atomMapCopy["children"].([]any); ok && len(children) > 0 {
		for _, child := range children {
			if childMap, ok := child.(map[string]any); ok {
				childAtom, err := MapToAtom(childMap)
				if err != nil {
					return nil, fmt.Errorf("failed to create child atom: %w", err)
				}
				// Add child and update the atom reference
				atom = atom.ChildAdd(childAtom).(*Atom)
			}
		}
	}

	// Final validation
	if atom.GetType() == "" {
		return nil, errors.New("missing required 'type' field in atom map")
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
//   - []*Atom: slice of decoded atoms
//   - error: if the data cannot be decoded or is invalid
func GobToAtoms(data []byte) ([]*Atom, error) {
	if len(data) == 0 {
		return []*Atom{}, nil
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

	result := make([]*Atom, 0, count)

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
		// Create the atom from the gob data
		atom, err := GobToAtom(atomData)
		if err != nil {
			return nil, fmt.Errorf("failed to create atom %d from gob: %w", i, err)
		}

		result = append(result, atom)
	}

	return result, nil
}

// GobToAtom decodes an atom from gob-encoded data.
//
// Business logic:
// - Decodes the binary data into a temporary struct
// - Creates a new atom with the decoded type and id from the data map
// - Sets all properties from the data map
// - Recursively decodes and adds child atoms
// - Returns an error if the data is invalid
//
// Parameters:
//   - data: binary data containing the gob-encoded atom
//
// Returns:
//   - *Atom: the decoded atom
//   - error: if decoding fails
func GobToAtom(data []byte) (*Atom, error) {
	// Validate the input data first
	if valid, err := isValidAtomGob(data); !valid {
		return nil, fmt.Errorf("invalid gob data: %w", err)
	}

	// Create a temporary struct for decoding
	var temp struct {
		Data     map[string]string
		Children [][]byte
	}

	// Decode the data (we know it's valid from the validation above)
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&temp); err != nil {
		// This should theoretically never happen since we already validated the data
		return nil, fmt.Errorf("unexpected error during gob decode: %w", err)
	}

	// Get type from data map (we know it exists from validation)
	atomType := temp.Data["type"]

	// Create a new atom with the decoded type
	atom := NewAtom(atomType)

	// Set all properties from the data map (including id if present)
	for key, value := range temp.Data {
		atom.Set(key, value)
	}

	// Recursively decode children
	for _, childData := range temp.Children {
		child, err := GobToAtom(childData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode child: %w", err)
		}
		// Use ChildAdd to add the child
		atom = atom.ChildAdd(child).(*Atom)
	}

	return atom, nil
}

// isValidAtomGob validates that the given data is a valid gob-encoded atom.
//
// Business logic:
// - Checks for empty or nil input
// - Attempts to decode the data into a temporary struct
// - Validates the presence of required fields (id, type)
// - Recursively validates child atoms
//
// Parameters:
//   - data: binary data to validate
//
// Returns:
//   - bool: true if the data is valid
//   - error: description of the validation failure if invalid
func isValidAtomGob(data []byte) (bool, error) {
	if len(data) == 0 {
		return false, errors.New("cannot validate empty data")
	}

	// Create a temporary struct for validation
	var temp struct {
		Data     map[string]string
		Children [][]byte
	}

	// Try to decode the data
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&temp); err != nil {
		return false, fmt.Errorf("invalid gob data: %w", err)
	}

	// Validate data map exists and contains required fields
	if temp.Data == nil {
		return false, errors.New("data map cannot be nil")
	}

	// Check for required fields
	if temp.Data["type"] == "" {
		return false, errors.New("missing required field: type")
	}

	// ID is optional, but if present it must be a non-empty string
	if id, exists := temp.Data["id"]; exists && id == "" {
		return false, errors.New("id cannot be empty if present")
	}

	// Recursively validate children
	for i, childData := range temp.Children {
		if valid, err := isValidAtomGob(childData); !valid {
			return false, fmt.Errorf("invalid child at index %d: %v", i, err)
		}
	}

	return true, nil
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
// - Validates that all properties are strings or in a nested properties map
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

	// Check for invalid top-level keys (only id, type, properties, children are allowed)
	for key := range atomMap {
		switch key {
		case "id", "type", "properties", "children":
			// These are valid top-level keys
			continue
		default:
			return false, fmt.Errorf("invalid top-level key '%s' in atom map, only 'id', 'type', 'properties', and 'children' are allowed", key)
		}
	}

	// Validate properties map if present
	if props, ok := atomMap["properties"]; ok && props != nil {
		propsMap, ok := props.(map[string]any)
		if !ok {
			return false, errors.New("properties must be a map[string]any")
		}

		// Validate that all property values are strings or convertible to strings
		for propKey, propValue := range propsMap {
			switch propValue.(type) {
			case string, fmt.Stringer, int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64, float32, float64, bool:
				// These types can be converted to strings
				continue
			default:
				return false, fmt.Errorf("property '%s' in properties map has invalid type %T, must be string or convertible to string", propKey, propValue)
			}
		}
	}

	// Validate children if present
	if children, ok := atomMap["children"]; ok && children != nil {
		childrenSlice, ok := children.([]any)
		if !ok {
			return false, errors.New("atom children must be a slice")
		}

		for i, child := range childrenSlice {
			if child == nil {
				continue // Skip nil children
			}

			childMap, ok := child.(map[string]any)
			if !ok {
				return false, fmt.Errorf("child at index %d is not a valid atom map (got type %T)", i, child)
			}

			if valid, err := isValidAtomMap(childMap); !valid {
				return false, fmt.Errorf("invalid child at index %d: %v", i, err)
			}
		}
	}

	return true, nil
}
