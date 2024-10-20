package main

import (
	"fmt"

	"github.com/codeis4fun/data-treatment-interpreter/internal/engine"
	"github.com/codeis4fun/data-treatment-interpreter/internal/lexer"
	"github.com/codeis4fun/data-treatment-interpreter/internal/parser"
)

func main() {
	// Sample JSON data
	jsonData := []byte(`{"name":"vanessa","surname":"teixeira","height":1.72,"weight":60,"location":"São Paulo/sp","languages":["Portuguese","English"],"friends":[{"first":"Dale","last":"Murphy","age":44,"nets":["ig","fb","tw"]},{"first":"Roger","last":"Craig","age":68,"nets":["fb","tw"]},{"first":"Jane","last":"Murphy","age":47,"nets":["ig","tw"]}]}`)

	input := `SET fullName = concatenate(' ',name, surname)
SET fullName = uppercase(fullName)
SET bmi, isHealty = bmi(weight, height)
SET description = concatenate(' BMI is ', fullName,  bmi)
SET _city, _state = split(location, '/')
SET address.city = uppercase(_city)
SET address.state = uppercase(_state)
SET languages.# = uppercase(languages.#)
SET friends.#.first = uppercase(friends.first)`
	// Initialize lexer and parser
	l := lexer.NewLexer(input)
	p := parser.NewParser(l, input)

	// Parse the input
	programs, err := p.RunAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Initialize engine
	e := engine.NewEngine()

	// Apply transformations to JSON
	modifiedJSON, err := e.ExecuteAll(programs, jsonData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(modifiedJSON)) // Output: {"completeName":"johndoe","name":"john","surname":"doe"}
}
