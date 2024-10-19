package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/codeis4fun/data-treatment-interpreter/internal/lexer"
	"github.com/codeis4fun/data-treatment-interpreter/internal/parser"
)

func TestParser(t *testing.T) {
	input := "SET a, b = concatenate(name, surname)"
	l := lexer.NewLexer(input)
	p := parser.NewParser(l, input)

	program, err := p.Run()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expectedProgram := &parser.Program{
		Variables:   []string{"a", "b"},
		Transformer: "concatenate",
		Args:        []string{"name", "surname"},
	}

	if program == nil {
		t.Fatalf("Expected program to be non-nil")
	}

	// Use reflect.DeepEqual to compare program and expectedProgram
	if !reflect.DeepEqual(program, expectedProgram) {
		t.Errorf("Expected program to be %v, got %v", expectedProgram, program)
	}
}

func TestParserWithSyntaxError(t *testing.T) {
	input := "SET a, b = concatenate(name, surname"
	l := lexer.NewLexer(input)
	p := parser.NewParser(l, input)

	program, err := p.Run()
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	if program != nil {
		t.Errorf("Expected program to be nil, got %v", p)
	}

	var expectedErrorBuilder strings.Builder
	expectedErrorBuilder.WriteString("unexpected token in arguments at position 36: ")
	expectedErrorBuilder.WriteString("\n")
	expectedErrorBuilder.WriteString("SET a, b = concatenate(name, surname")
	expectedErrorBuilder.WriteString("\n")
	expectedErrorBuilder.WriteString(`                                    ^`)

	expectedError := expectedErrorBuilder.String()
	if err.Error() != expectedError {
		t.Errorf("Expected error message to be '%s', got '%s'", expectedError, err.Error())
	}
}
