package omni

import "testing"

func TestIsValidAtomJSON_EmptyString_ReturnsError(t *testing.T) {
	ok, err := isValidAtomJSON("")
	if err == nil {
		t.Fatal("expected error for empty string, got nil")
	}
	if ok {
		t.Fatalf("expected ok=false for empty string, got %v", ok)
	}
}

func TestIsValidAtomJSON_NullString_ReturnsError(t *testing.T) {
	ok, err := isValidAtomJSON("null")
	if err == nil {
		t.Fatal("expected error for 'null', got nil")
	}
	if ok {
		t.Fatalf("expected ok=false for 'null', got %v", ok)
	}
}

func TestIsValidAtomJSON_EmptyArray_IsValid(t *testing.T) {
	ok, err := isValidAtomJSON("[]")
	if err != nil {
		t.Fatalf("unexpected error for empty array: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true for empty array")
	}
}

func TestIsValidAtomJSON_InvalidRoot_ReturnsError(t *testing.T) {
	ok, err := isValidAtomJSON("42")
	if err == nil || ok {
		t.Fatalf("expected invalid for non-object/array root, got ok=%v err=%v", ok, err)
	}
}

func TestIsValidAtomMap_InvalidTopLevelKey(t *testing.T) {
	m := map[string]any{
		"id":   "id1",
		"type": "type1",
		"bad":  true,
	}
	ok, err := isValidAtomMap(m)
	if ok || err == nil {
		t.Fatalf("expected invalid due to extra top-level key, got ok=%v err=%v", ok, err)
	}
}

func TestIsValidAtomMap_InvalidPropertiesType(t *testing.T) {
	m := map[string]any{
		"id":         "id1",
		"type":       "type1",
		"properties": 123, // not a map
	}
	ok, err := isValidAtomMap(m)
	if ok || err == nil {
		t.Fatalf("expected invalid due to properties type, got ok=%v err=%v", ok, err)
	}
}

func TestIsValidAtomMap_InvalidChildrenType(t *testing.T) {
	m := map[string]any{
		"id":       "id1",
		"type":     "type1",
		"children": "not-a-slice",
	}
	ok, err := isValidAtomMap(m)
	if ok || err == nil {
		t.Fatalf("expected invalid due to children type, got ok=%v err=%v", ok, err)
	}
}
