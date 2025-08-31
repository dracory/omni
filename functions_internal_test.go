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
