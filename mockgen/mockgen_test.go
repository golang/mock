package main

import "testing"

func TestNewIdentifierAllocator(t *testing.T) {
	a := newIdentifierAllocator([]string{"taken1", "taken2"})
	if len(a) != 2 {
		t.Fatalf("expected 2 items, got %v", len(a))
	}

	_, ok := a["taken1"]
	if !ok {
		t.Errorf("allocator doesn't contain 'taken1': %#v", a)
	}

	_, ok = a["taken2"]
	if !ok {
		t.Errorf("allocator doesn't contain 'taken2': %#v", a)
	}
}

func allocatorContainsIdentifiers(a identifierAllocator, ids []string) bool {
	if len(a) != len(ids) {
		return false
	}

	for _, id := range ids {
		_, ok := a[id]
		if !ok {
			return false
		}
	}

	return true
}

func TestIdentifierAllocator_allocateIdentifier(t *testing.T) {
	a := newIdentifierAllocator([]string{"taken"})

	t2 := a.allocateIdentifier("taken_2")
	if t2 != "taken_2" {
		t.Fatalf("expected 'taken_2', got %q", t2)
	}
	expected := []string{"taken", "taken_2"}
	if !allocatorContainsIdentifiers(a, expected) {
		t.Fatalf("allocator doesn't contain the expected items - allocator: %#v, expected items: %#v", a, expected)
	}

	t3 := a.allocateIdentifier("taken")
	if t3 != "taken_3" {
		t.Fatalf("expected 'taken_3', got %q", t3)
	}
	expected = []string{"taken", "taken_2", "taken_3"}
	if !allocatorContainsIdentifiers(a, expected) {
		t.Fatalf("allocator doesn't contain the expected items - allocator: %#v, expected items: %#v", a, expected)
	}

	t4 := a.allocateIdentifier("taken")
	if t4 != "taken_4" {
		t.Fatalf("expected 'taken_4', got %q", t4)
	}
	expected = []string{"taken", "taken_2", "taken_3", "taken_4"}
	if !allocatorContainsIdentifiers(a, expected) {
		t.Fatalf("allocator doesn't contain the expected items - allocator: %#v, expected items: %#v", a, expected)
	}

	id := a.allocateIdentifier("id")
	if id != "id" {
		t.Fatalf("expected 'id', got %q", id)
	}
	expected = []string{"taken", "taken_2", "taken_3", "taken_4", "id"}
	if !allocatorContainsIdentifiers(a, expected) {
		t.Fatalf("allocator doesn't contain the expected items - allocator: %#v, expected items: %#v", a, expected)
	}
}
