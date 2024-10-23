package main

import (
	"fmt"
	"strings"

	"github.com/codeis4fun/data-treatment-interpreter/internal/engine"
	"github.com/codeis4fun/data-treatment-interpreter/internal/lexer"
	"github.com/codeis4fun/data-treatment-interpreter/internal/parser"
)

func main() {
	// Sample JSON data
	// jsonData := []byte(`{"name":"vanessa","surname":"teixeira","height":1.72,"weight":60,"location":"São Paulo/sp","languages":["Portuguese","English"],"friends":[{"first":"Dale","last":"Murphy","age":44,"nets":["ig","fb","tw"]},{"first":"Roger","last":"Craig","age":68,"nets":["fb","tw"]},{"first":"Jane","last":"Murphy","age":47,"nets":["ig","tw"]}]}`)
	jsonData := []byte(`{"firstName":"john","lastName":"doe","weight":75,"height":1.75,"favoriteFoods":["pizza","pasta","sushi"],"favoriteColors":["red","blue","green"],"place":"New York/USA","friends":[{"name":"Alice"},{"name":"Bob"}]}`)
	input := `SET _tempName = concatenate(' ', firstName, lastName)
SET fullName = uppercase(_tempName)
SET bmi, isHealty = bmi(weight, height)
SET favoriteFoods.0 = uppercase(favoriteFoods.0)
SET favoriteColors.# = uppercase(favoriteColors.#)
SET _city, _country = split(place, '/')
SET address.city = uppercase(_city)
SET address.country = uppercase(_country)
SET friends.#.name = uppercase(friends.#.name)`
	r := strings.NewReader(input)
	// Initialize lexer and parser
	l := lexer.NewLexer(r)
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
