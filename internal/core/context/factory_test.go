package context

import "testing"

func TestNewSelectorValid(t *testing.T) {
	for _, policy := range []string{"", "ci", "k-callsite-0", "k-callsite-2", "k-callsite-10"} {
		selector, err := NewSelector(policy)
		if err != nil {
			t.Fatalf("NewSelector(%q): %v", policy, err)
		}
		if selector == nil {
			t.Fatalf("NewSelector(%q) returned nil", policy)
		}
	}
}

func TestNewSelectorInvalid(t *testing.T) {
	for _, policy := range []string{"k-callsite-", "k-callsite--1", "k-callsite-x", "k-obj", "unknown"} {
		if _, err := NewSelector(policy); err == nil {
			t.Fatalf("NewSelector(%q) succeeded, want error", policy)
		}
	}
}
