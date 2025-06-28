package omni

// import (
// 	"encoding/json"
// 	"errors"
// )

// func MarshalAtomsToJson(atoms []AtomInterface) (string, error) {
// 	atomsMap := []atomJsonObject{}

// 	for _, atom := range atoms {
// 		atomsMap = append(atomsMap, atom.ToJsonObject())
// 	}

// 	atomsJson, err := json.Marshal(atomsMap)

// 	return string(atomsJson), err
// }

// func UnmarshalJsonToAtoms(atomsJson string) ([]AtomInterface, error) {
// 	atomsMap := []map[string]any{}

// 	err := json.Unmarshal([]byte(atomsJson), &atomsMap)

// 	if err != nil {
// 		return nil, err
// 	}

// 	atoms := []AtomInterface{}

// 	for _, atomMap := range atomsMap {
// 		atom, err := ConvertMapToAtom(atomMap)

// 		if err != nil {
// 			return nil, err
// 		}

// 		atoms = append(atoms, atom)
// 	}

// 	return atoms, nil
// }

// func ConvertMapToAtoms(atoms []map[string]any) []AtomInterface {
// 	atomsMap := []AtomInterface{}

// 	for _, atom := range atoms {
// 		atomsMap = append(atomsMap, NewAtomFromMap(atom))
// 	}

// 	return atomsMap
// }

// func ConvertAtomsToMap(atoms []AtomInterface) []map[string]any {
// 	atomsMap := []map[string]any{}

// 	for _, atom := range atoms {
// 		atomsMap = append(atomsMap, atom.ToMap())
// 	}

// 	return atomsMap
// }

// // ConvertMapToAtom converts a map to a atom
// //
// // The map must represent a valid atom (have parameters like id, and type),
// // otherwise an error will be returned
// //
// // Parameters:
// // - atomMap - a map[string]any to convert to a atom
// //
// // Returns:
// // - AtomInterface - a atom
// // - error - if the map[string]any is not a valid atom
// func ConvertMapToAtom(atomMap map[string]any) (AtomInterface, error) {
// 	atomMap, err := mapToAtomMap(atomMap)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return NewAtomFromMap(atomMap), nil
// }

// // mapToAtomMap converts a map[string]any to a map[string]any
// // the map[string]any must be a valid atom, otherwise an error
// // will be returned
// //
// // Parameters:
// // - atomMap - a map[string]any to convert to a atom
// //
// // Returns:
// // - map[string]any - a atom
// // - error - if the map[string]any is not a valid atom
// func mapToAtomMap(atomMap map[string]any) (map[string]any, error) {
// 	idAny, ok := atomMap["id"]

// 	if !ok {
// 		return nil, errors.New("id not found")
// 	}

// 	typeAny, ok := atomMap["type"]

// 	if !ok {
// 		return nil, errors.New("type not found")
// 	}

// 	parametersAny, ok := atomMap["parameters"]

// 	if !ok {
// 		parametersAny = map[string]any{}
// 	}

// 	childrenAny, ok := atomMap["children"]

// 	if !ok {
// 		childrenAny = []any{}
// 	}

// 	childrenArrayAny := childrenAny.([]any)

// 	childrenMap := []map[string]any{}
// 	for _, childAny := range childrenArrayAny {
// 		childAny := childAny.(map[string]any)
// 		child, err := mapToAtomMap(childAny)

// 		if err != nil {
// 			return nil, err
// 		}

// 		childrenMap = append(childrenMap, child)
// 	}

// 	parametersMapAny := parametersAny.(map[string]any)
// 	parametersMap := map[string]string{}

// 	for k, v := range parametersMapAny {
// 		parametersMap[k] = v.(string)
// 	}

// 	atomMap["id"] = idAny.(string)
// 	atomMap["type"] = typeAny.(string)
// 	atomMap["parameters"] = parametersMap
// 	atomMap["children"] = childrenMap

// 	return atomMap, nil
// }
