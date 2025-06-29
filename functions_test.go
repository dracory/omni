package omni_test

import (
	"strings"
	"testing"

	"github.com/dracory/omni"
)

func TestAtomsToJSON(t *testing.T) {
	// Create test atoms
	atom1 := omni.NewAtom("type1", omni.WithID("id1"))
	atom1.Set("key1", "value1")

	atom2 := omni.NewAtom("type2", omni.WithID("id2"))
	childAtom := omni.NewAtom("childType", omni.WithID("childId"))
	childAtom.Set("childKey", "childValue")
	atom2.ChildAdd(childAtom)

	// Test with single atom
	jsonStr, err := omni.AtomsToJSON([]omni.AtomInterface{atom1})
	if err != nil {
		t.Fatalf("AtomsToJSON() error = %v", err)
	}
	if want := `"id":"id1"`; !strings.Contains(jsonStr, want) {
		t.Errorf("AtomsToJSON() = %v, want contains %v", jsonStr, want)
	}
	if want := `"type":"type1"`; !strings.Contains(jsonStr, want) {
		t.Errorf("AtomsToJSON() = %v, want contains %v", jsonStr, want)
	}
	if want := `"key1":"value1"`; !strings.Contains(jsonStr, want) {
		t.Errorf("AtomsToJSON() = %v, want contains %v", jsonStr, want)
	}

	// Test with multiple atoms
	jsonStr, err = omni.AtomsToJSON([]omni.AtomInterface{atom1, atom2})
	if err != nil {
		t.Fatalf("AtomsToJSON() error = %v", err)
	}
	for _, want := range []string{`"id1"`, `"id2"`, `"childId"`} {
		if !strings.Contains(jsonStr, want) {
			t.Errorf("AtomsToJSON() = %v, want contains %v", jsonStr, want)
		}
	}

	// Test with empty slice
	jsonStr, err = omni.AtomsToJSON([]omni.AtomInterface{})
	if err != nil {
		t.Fatalf("AtomsToJSON() error = %v", err)
	}
	if want := "[]"; jsonStr != want {
		t.Errorf("AtomsToJSON() = %v, want %v", jsonStr, want)
	}
}

func TestJSONToAtoms(t *testing.T) {
	// Test with valid JSON array
	jsonStr := `[{"data":{"id":"id1","type":"type1","key1":"value1"},"children":[]}]`

	atoms, err := omni.JSONToAtoms(jsonStr)
	if err != nil {
		t.Fatalf("JSONToAtoms() error = %v", err)
	}
	if len(atoms) != 1 {
		t.Fatalf("JSONToAtoms() len = %v, want %v", len(atoms), 1)
	}
	if id := atoms[0].GetID(); id != "id1" {
		t.Errorf("atoms[0].GetID() = %v, want %v", id, "id1")
	}
	if typ := atoms[0].GetType(); typ != "type1" {
		t.Errorf("atoms[0].GetType() = %v, want %v", typ, "type1")
	}
	if prop := atoms[0].Get("key1"); prop != "value1" {
		t.Errorf("atoms[0].GetProperty(\"key1\") = %v, want %v", prop, "value1")
	}

	// Test with invalid JSON
	_, err = omni.JSONToAtoms("{invalid-json}")
	if err == nil {
		t.Error("JSONToAtoms() error = nil, want error")
	}

	// Test with empty string (should return empty slice, not error)
	atoms, err = omni.JSONToAtoms(`""`)
	if err != nil {
		t.Fatalf("JSONToAtoms() error = %v", err)
	}
	if len(atoms) != 0 {
		t.Errorf("JSONToAtoms() len = %v, want %v", len(atoms), 0)
	}
}

func TestConvertMapToAtoms(t *testing.T) {
	// Test with valid maps
	maps := []map[string]any{
		{
			"id":   "id1",
			"type": "type1",
		},
		{
			"id":   "id2",
			"type": "type2",
		},
	}

	atoms := omni.MapToAtoms(maps)
	if len(atoms) != 2 {
		t.Fatalf("MapToAtoms() len = %v, want %v", len(atoms), 2)
	}
	if id := atoms[0].GetID(); id != "id1" {
		t.Errorf("atoms[0].GetID() = %v, want %v", id, "id1")
	}
	if typ := atoms[0].GetType(); typ != "type1" {
		t.Errorf("atoms[0].GetType() = %v, want %v", typ, "type1")
	}
	if id := atoms[1].GetID(); id != "id2" {
		t.Errorf("atoms[1].GetID() = %v, want %v", id, "id2")
	}

	// Test with empty slice
	atoms = omni.MapToAtoms([]map[string]any{})
	if len(atoms) != 0 {
		t.Errorf("MapToAtoms() len = %v, want %v", len(atoms), 0)
	}
}

func TestAtomsToMap(t *testing.T) {
	// Create test atoms
	atom1 := omni.NewAtom("type1", omni.WithID("id1"))
	atom1.Set("key1", "value1")

	atom2 := omni.NewAtom("type2", omni.WithID("id2"))
	childAtom := omni.NewAtom("childType", omni.WithID("childId"))
	// Make sure we add the child to atom2
	atom2.ChildAdd(childAtom)
	// Verify child was added
	if atom2.ChildrenLength() != 1 {
		t.Fatalf("Expected 1 child, got %d", atom2.ChildrenLength())
	}

	// Convert to maps
	maps := omni.AtomsToMap([]omni.AtomInterface{atom1, atom2})
	if len(maps) != 2 {
		t.Fatalf("AtomsToMap() len = %v, want %v", len(maps), 2)
	}

	// Verify first atom
	if id, ok := maps[0]["id"].(string); !ok || id != "id1" {
		t.Errorf("maps[0][\"id\"] = %v, want %v", id, "id1")
	}
	if typ, ok := maps[0]["type"].(string); !ok || typ != "type1" {
		t.Errorf("maps[0][\"type\"] = %v, want %v", typ, "type1")
	}
	if params, ok := maps[0]["properties"].(map[string]string); !ok || params["key1"] != "value1" {
		t.Errorf("maps[0][\"properties\"] missing key1=value1")
	}

	// Verify second atom's children
	// Print the actual type and value for debugging
	t.Logf("maps[1][\"children\"] type: %T, value: %+v", maps[1]["children"], maps[1]["children"])

	// Try different type assertions
	var children []interface{}
	var ok bool

	// First try as []interface{}
	children, ok = maps[1]["children"].([]interface{})
	if !ok {
		// Then try as []map[string]interface{}
		if childrenMaps, mapOk := maps[1]["children"].([]map[string]interface{}); mapOk {
			children = make([]interface{}, len(childrenMaps))
			for i, m := range childrenMaps {
				children[i] = m
			}
			ok = true
		}
	}

	if !ok || len(children) != 1 {
		t.Fatalf("maps[1][\"children\"] is not a slice or has wrong length: %v, want 1", maps[1]["children"])
	}

	// Get the child map
	childMap, ok := children[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Child is not a map: %T %+v", children[0], children[0])
	}

	// Verify child ID and type
	childID, idOk := childMap["id"].(string)
	if !idOk || childID != "childId" {
		t.Errorf("child id = %v (type %T), want %v", childMap["id"], childMap["id"], "childId")
	}

	childType, typeOk := childMap["type"].(string)
	if !typeOk || childType != "childType" {
		t.Errorf("child type = %v (type %T), want %v", childMap["type"], childMap["type"], "childType")
	}

	// Test with empty slice
	maps = omni.AtomsToMap([]omni.AtomInterface{})
	if len(maps) != 0 {
		t.Errorf("AtomsToMap() len = %v, want %v", len(maps), 0)
	}
}

func TestAtomsToGob(t *testing.T) {
	// Create test atoms
	atom1 := omni.NewAtom("type1", omni.WithID("id1"))
	atom1.Set("key1", "value1")

	atom2 := omni.NewAtom("type2", omni.WithID("id2"))
	childAtom := omni.NewAtom("childType", omni.WithID("childId"))
	childAtom.Set("childKey", "childValue")
	atom2.ChildAdd(childAtom)

	// Test with valid atoms
	gobData, err := omni.AtomsToGob([]omni.AtomInterface{atom1, atom2})
	if err != nil {
		t.Fatalf("AtomsToGob() error = %v", err)
	}

	// Should not be empty
	if len(gobData) == 0 {
		t.Error("AtomsToGob() returned empty data")
	}

	// Test roundtrip
	decodedAtoms, err := omni.GobToAtoms(gobData)
	if err != nil {
		t.Fatalf("GobToAtoms(AtomsToGob()) error = %v", err)
	}

	if len(decodedAtoms) != 2 {
		t.Fatalf("Expected 2 atoms after roundtrip, got %d", len(decodedAtoms))
	}

	// Test with empty input
	emptyData, err := omni.AtomsToGob([]omni.AtomInterface{})
	if err != nil {
		t.Fatalf("AtomsToGob() with empty input error = %v", err)
	}
	if len(emptyData) != 0 {
		t.Error("AtomsToGob() with empty input should return empty slice")
	}

	// Test with nil input
	nilData, err := omni.AtomsToGob(nil)
	if err != nil {
		t.Fatalf("AtomsToGob() with nil input error = %v", err)
	}
	if len(nilData) != 0 {
		t.Error("AtomsToGob() with nil input should return empty slice")
	}
}

func TestGobToAtoms(t *testing.T) {
	// Create test atoms
	atom1 := omni.NewAtom("type1", omni.WithID("id1"))
	atom1.Set("key1", "value1")

	atom2 := omni.NewAtom("type2", omni.WithID("id2"))
	childAtom := omni.NewAtom("childType", omni.WithID("childId"))
	childAtom.Set("childKey", "childValue")
	atom2.ChildAdd(childAtom)

	// Encode the atoms to gob using our helper function
	gobData, err := omni.AtomsToGob([]omni.AtomInterface{atom1, atom2})
	if err != nil {
		t.Fatalf("Failed to encode atoms to gob: %v", err)
	}

	// Test with valid gob data
	decodedAtoms, err := omni.GobToAtoms(gobData)
	if err != nil {
		t.Fatalf("GobToAtoms() error = %v", err)
	}

	if len(decodedAtoms) != 2 {
		t.Fatalf("GobToAtoms() len = %v, want %v", len(decodedAtoms), 2)
	}

	// Verify first atom
	if id := decodedAtoms[0].GetID(); id != "id1" {
		t.Errorf("decodedAtoms[0].GetID() = %v, want %v", id, "id1")
	}
	if typ := decodedAtoms[0].GetType(); typ != "type1" {
		t.Errorf("decodedAtoms[0].GetType() = %v, want %v", typ, "type1")
	}

	// Verify second atom and its child
	if id := decodedAtoms[1].GetID(); id != "id2" {
		t.Errorf("decodedAtoms[1].GetID() = %v, want %v", id, "id2")
	}
	children := decodedAtoms[1].ChildrenGet()
	if len(children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(children))
	}

	// Test with empty input
	emptyAtoms, err := omni.GobToAtoms([]byte{})
	if err != nil {
		t.Fatalf("GobToAtoms() with empty input error = %v", err)
	}
	if len(emptyAtoms) != 0 {
		t.Errorf("GobToAtoms() with empty input len = %v, want %v", len(emptyAtoms), 0)
	}

	// Test with invalid gob data
	_, err = omni.GobToAtoms([]byte("invalid-gob-data"))
	if err == nil {
		t.Error("GobToAtoms() with invalid data: expected error, got nil")
	}
}

func TestMapToAtom(t *testing.T) {
	// Test with valid map
	atomMap := map[string]any{
		"id":   "testId",
		"type": "testType",
		"properties": map[string]any{
			"key1": "value1",
		},
		"children": []any{
			map[string]any{
				"id":   "childId",
				"type": "childType",
			},
		},
	}

	atom, err := omni.MapToAtom(atomMap)
	if err != nil {
		t.Fatalf("MapToAtom() error = %v", err)
	}
	if id := atom.GetID(); id != "testId" {
		t.Errorf("atom.GetID() = %v, want %v", id, "testId")
	}
	if typ := atom.GetType(); typ != "testType" {
		t.Errorf("atom.GetType() = %v, want %v", typ, "testType")
	}
	prop := atom.Get("key1")
	if prop != "value1" {
		t.Errorf("atom.Get(\"key1\") = %v, want %v", prop, "value1")
	}
	children := atom.ChildrenGet()
	if len(children) != 1 {
		t.Fatalf("len(atom.GetChildren()) = %v, want %v", len(children), 1)
	}
	if childID := children[0].GetID(); childID != "childId" {
		t.Errorf("children[0].GetID() = %v, want %v", childID, "childId")
	}

	// Test with missing required fields
	tests := []struct {
		name    string
		m       map[string]any
		wantErr bool
	}{
		{"missing id", map[string]any{"type": "missingId"}, true},
		{"missing type", map[string]any{"id": "missingType"}, true},
		{"nil map", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := omni.MapToAtom(tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapToAtom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
