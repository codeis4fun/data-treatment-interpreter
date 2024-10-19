package transformers_test

import (
	"testing"

	"github.com/codeis4fun/data-treatment-interpreter/internal/transformers"
)

func TestBMITransform(t *testing.T) {
	transformer := &transformers.BMI{
		Config: transformers.Config{
			Args: []string{"weight", "height"},
			Json: []byte(`{"height": 1.72, "weight": 60}`),
		},
	}

	results, err := transformer.Transform()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	var expected float64 = 20.3
	if results[0] != expected {
		t.Errorf("Expected %f, got %f", expected, results[0])
	}
}

func TestBMITransformWithNonExistentField(t *testing.T) {
	transformer := &transformers.BMI{
		Config: transformers.Config{
			Args: []string{"nonexistent", "height"},
			Json: []byte(`{"height": 1.72, "weight": 60}`),
		},
	}

	_, err := transformer.Transform()
	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}
}
