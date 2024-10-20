package engine_test

import (
	"strings"
	"testing"

	"github.com/codeis4fun/data-treatment-interpreter/internal/engine"
	"github.com/codeis4fun/data-treatment-interpreter/internal/lexer"
	"github.com/codeis4fun/data-treatment-interpreter/internal/parser"
)

func TestEngineExecute(t *testing.T) {
	// Sample JSON data
	jsonData := []byte(`{"name": "john", "surname": "doe"}`)

	// Input transformation: SET completeName = concatenate(name, surname)
	program := &parser.Program{
		Variables:   []string{"completeName"},
		Transformer: "concatenate",
		Args:        []string{"' '", "name", "surname"},
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	modifiedJSON, err := e.Execute(program, jsonData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `{"name": "john", "surname": "doe","completeName":"john doe"}`
	if string(modifiedJSON) != expected {
		t.Errorf("Expected %s, got %s", expected, string(modifiedJSON))
	}
}

func TestEngineExecuteBMI(t *testing.T) {
	// Sample JSON data
	jsonData := []byte(`{"height": 1.72, "weight": 60}`)

	// Input transformation: SET bmi = bmi(weight, height)
	program := &parser.Program{
		Variables:   []string{"bmi", "isHealthy"},
		Transformer: "bmi",
		Args:        []string{"weight", "height"},
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	modifiedJSON, err := e.Execute(program, jsonData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `{"height": 1.72, "weight": 60,"bmi":20.3,"isHealthy":true}`
	if string(modifiedJSON) != expected {
		t.Errorf("Expected %s, got %s", expected, string(modifiedJSON))
	}
}

func TestEngineExecuteUppercase(t *testing.T) {
	// Sample JSON data
	jsonData := []byte(`{"name": "john", "surname": "doe"}`)

	// Input transformation: SET completeName = uppercase(name)
	program := &parser.Program{
		Variables:   []string{"completeName"},
		Transformer: "uppercase",
		Args:        []string{"name"},
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	modifiedJSON, err := e.Execute(program, jsonData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `{"name": "john", "surname": "doe","completeName":"JOHN"}`
	if string(modifiedJSON) != expected {
		t.Errorf("Expected %s, got %s", expected, string(modifiedJSON))
	}
}

func TestEngineExecuteSplit(t *testing.T) {
	// Sample JSON data
	jsonData := []byte(`{"location": "São Paulo/sp"}`)

	// Input transformation: SET city, state = split(location, '/')
	program := &parser.Program{
		Variables:   []string{"city", "state"},
		Transformer: "split",
		Args:        []string{"location", "'/'"},
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	modifiedJSON, err := e.Execute(program, jsonData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `{"location": "São Paulo/sp","city":"São Paulo","state":"sp"}`
	if string(modifiedJSON) != expected {
		t.Errorf("Expected %s, got %s", expected, string(modifiedJSON))
	}
}

func TestEngineExecuteAll(t *testing.T) {
	// Sample JSON data
	jsonData := []byte(`{"name": "john", "surname": "doe", "height": 1.72, "weight": 60, "location": "São Paulo/sp"}`)

	// Input transformation program
	programs := []*parser.Program{
		{
			Variables:   []string{"completeName"},
			Transformer: "concatenate",
			Args:        []string{"' '", "name", "surname"},
		},
		{
			Variables:   []string{"completeName"},
			Transformer: "uppercase",
			Args:        []string{"completeName"},
		},
		{
			Variables:   []string{"bmi", "isHealthy"},
			Transformer: "bmi",
			Args:        []string{"weight", "height"},
		},
		{
			Variables:   []string{"description"},
			Transformer: "concatenate",
			Args:        []string{"' BMI is '", "completeName", "bmi"},
		},
		{
			Variables:   []string{"city", "state"},
			Transformer: "split",
			Args:        []string{"location", "'/'"},
		},
		{
			Variables:   []string{"address.city"},
			Transformer: "uppercase",
			Args:        []string{"city"},
		},
		{
			Variables:   []string{"address.state"},
			Transformer: "uppercase",
			Args:        []string{"state"},
		},
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	modifiedJSON, err := e.ExecuteAll(programs, jsonData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `{"name": "john", "surname": "doe", "height": 1.72, "weight": 60, "location": "São Paulo/sp","completeName":"JOHN DOE","bmi":20.3,"isHealthy":true,"description":"JOHN DOE BMI is 20.3","city":"São Paulo","state":"sp","address":{"city":"SÃO PAULO","state":"SP"}}`
	if string(modifiedJSON) != expected {
		t.Errorf("Expected %s, got %s", expected, string(modifiedJSON))
	}
}

func TestEngineExecuteAllWithErrors(t *testing.T) {
	// Sample JSON data
	jsonData := []byte(`{"name": "john", "surname": "doe", "height": 1.72, "weight": 60, "location": "São Paulo/sp"}`)

	// Input transformation program
	programs := []*parser.Program{
		{
			Variables:   []string{"address.state"},
			Transformer: "uppercase",
			Args:        []string{"state"},
		},
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	_, err := e.ExecuteAll(programs, jsonData)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestEngineWithIterations(t *testing.T) {
	// Sample JSON data
	jsonData := []byte(`{"friends":[{"first":"Dale","last":"Murphy"},{"first":"Roger","last":"Craig"},{"first":"Jane","last":"Murphy"}]}`)

	// Input transformation: SET friends.#.first = uppercase(friends.#.first)
	program := &parser.Program{
		Variables:   []string{"friends.#.first"},
		Transformer: "uppercase",
		Args:        []string{"friends.#.first"},
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	modifiedJSON, err := e.Execute(program, jsonData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `{"friends":[{"first":"DALE","last":"Murphy"},{"first":"ROGER","last":"Craig"},{"first":"JANE","last":"Murphy"}]}`
	if string(modifiedJSON) != expected {
		t.Errorf("Expected %s, got %s", expected, string(modifiedJSON))
	}
}

func TestEngineWithIterationsAndError(t *testing.T) {
	// Sample JSON data
	jsonData := []byte(`{"friends":[{"first":"Dale","last":"Murphy"},{"first":"Roger","last":"Craig"},{"first":"Jane","last":"Murphy"}]}`)

	// Input transformation: SET friends.#.first = uppercase(friends.first)
	program := &parser.Program{
		Variables:   []string{"friends.first"},
		Transformer: "uppercase",
		Args:        []string{"friends.first"},
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	_, err := e.Execute(program, jsonData)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestEngineWithMultiplePrograms(t *testing.T) {
	// Sample JSON data
	jsonData := []byte(`{"name": "john", "surname": "doe"}`)

	input := `SET completeName = concatenate(' ',name, surname)
SET completeName = uppercase(completeName)`

	// Initialize lexer and parser
	l := lexer.NewLexer(strings.NewReader(input))
	p := parser.NewParser(l, input)

	// Parse the input
	programs, err := p.RunAll()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Input transformation program
	ExpectedPrograms := []*parser.Program{
		{
			Variables:   []string{"completeName"},
			Transformer: "concatenate",
			Args:        []string{"' '", "name", "surname"},
		},
		{
			Variables:   []string{"completeName"},
			Transformer: "uppercase",
			Args:        []string{"completeName"},
		},
	}

	// Use reflect.DeepEqual to compare program and expectedProgram
	if len(programs) != len(ExpectedPrograms) {
		t.Fatalf("Expected %d programs, got %d", len(ExpectedPrograms), len(programs))
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	modifiedJSON, err := e.ExecuteAll(programs, jsonData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `{"name": "john", "surname": "doe","completeName":"JOHN DOE"}`
	if string(modifiedJSON) != expected {
		t.Errorf("Expected %s, got %s", expected, string(modifiedJSON))
	}
}
