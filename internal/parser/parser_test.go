package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/codeis4fun/data-treatment-interpreter/internal/lexer"
	"github.com/codeis4fun/data-treatment-interpreter/internal/parser"
)

func TestParser(t *testing.T) {
	input := "SET a, b = t(c, d)"
	r := strings.NewReader(input)
	l := lexer.NewLexer(r)
	p := parser.NewParser(l, input)

	program, err := p.Run()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expectedProgram := &parser.Program{
		Variables:   []string{"a", "b"},
		Transformer: "t",
		Args:        []string{"c", "d"},
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
	r := strings.NewReader(input)
	l := lexer.NewLexer(r)
	p := parser.NewParser(l, input)

	program, err := p.Run()
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	if program != nil {
		t.Errorf("Expected program to be nil, got %v", p)
	}

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParserWithIterations(t *testing.T) {
	input := "SET friends.#.first = uppercase(friends.#.first)"
	r := strings.NewReader(input)
	l := lexer.NewLexer(r)
	p := parser.NewParser(l, input)

	programs, err := p.RunAll()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(programs) != 1 {
		t.Fatalf("Expected 1 program, got %d", len(programs))
	}

	program := programs[0]

	expectedProgram := &parser.Program{
		Variables:   []string{"friends.#.first"},
		Transformer: "uppercase",
		Args:        []string{"friends.#.first"},
	}

	if !reflect.DeepEqual(program, expectedProgram) {
		t.Errorf("Expected program to be %v, got %v", expectedProgram, program)
	}

}

func TestParserWithMultiplePrograms(t *testing.T) {
	input := `SET a = t(b, c)
SET d = t(e, f)`
	r := strings.NewReader(input)
	l := lexer.NewLexer(r)
	p := parser.NewParser(l, input)

	programs, err := p.RunAll()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(programs) != 2 {
		t.Fatalf("Expected 2 programs, got %d", len(programs))
	}

	program1 := programs[0]
	program2 := programs[1]

	expectedProgram1 := &parser.Program{
		Variables:   []string{"a"},
		Transformer: "t",
		Args:        []string{"b", "c"},
	}

	expectedProgram2 := &parser.Program{
		Variables:   []string{"d"},
		Transformer: "t",
		Args:        []string{"e", "f"},
	}

	if !reflect.DeepEqual(program1, expectedProgram1) {
		t.Errorf("Expected program to be %v, got %v", expectedProgram1, program1)
	}

	if !reflect.DeepEqual(program2, expectedProgram2) {
		t.Errorf("Expected program to be %v, got %v", expectedProgram2, program2)
	}

}
