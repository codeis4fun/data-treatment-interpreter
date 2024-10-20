package transformers_test

import (
	"testing"

	"github.com/codeis4fun/data-treatment-interpreter/internal/transformers"
)

func TestUppercaseTransform(t *testing.T) {
	transformer := &transformers.Uppercase{
		Config: transformers.Config{
			Args: []string{"name"},
			Json: []byte(`{"name": "john"}`),
		},
	}

	results, err := transformer.Transform()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "JOHN"
	if results[0] != expected {
		t.Errorf("Expected %s, got %s", expected, results[0])
	}
}

func TestUppercaseTransformWithNonExistentField(t *testing.T) {
	transformer := &transformers.Uppercase{
		Config: transformers.Config{
			Args: []string{"nonexistent"},
			Json: []byte(`{"name": "john"}`),
		},
	}

	_, err := transformer.Transform()
	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}
}

func TestConcatenateTransform(t *testing.T) {
	transformer := &transformers.Concatenate{
		Config: transformers.Config{
			Args: []string{"' '", "name", "surname"},
			Json: []byte(`{"name": "john", "surname": "doe"}`),
		},
	}

	results, err := transformer.Transform()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "john doe"
	if results[0] != expected {
		t.Errorf("Expected %s, got %s", expected, results[0])
	}
}

func TestConcatenateTransformWithNonExistentField(t *testing.T) {
	transformer := &transformers.Concatenate{
		Config: transformers.Config{
			Args: []string{"' '", "name", "nonexistent"},
			Json: []byte(`{"name": "john", "surname": "doe"}`),
		},
	}

	_, err := transformer.Transform()
	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}
}

func TestSplit(t *testing.T) {
	transformer := &transformers.Split{
		Config: transformers.Config{
			Args: []string{"name", "/"},
			Json: []byte(`{"name": "john/doe"}`),
		},
	}

	results, err := transformer.Transform()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []string{"john", "doe"}
	if results[0] != expected[0] || results[1] != expected[1] {
		t.Errorf("Expected %v, got %v", expected, results)
	}
}

func TestSplitWithNonExistentField(t *testing.T) {
	transformer := &transformers.Split{
		Config: transformers.Config{
			Args: []string{"' '", "nonexistent"},
			Json: []byte(`{"name": "john doe"}`),
		},
	}

	_, err := transformer.Transform()
	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}
}
