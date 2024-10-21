package transformers

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// Uppercase struct holds the arguments and JSON data for transformation
type Uppercase struct {
	Config
}

// Transform converts the input to uppercase
func (t *Uppercase) Transform() (Results, error) {
	if len(t.Args) != 1 {
		return nil, fmt.Errorf("uppercase requires exactly one argument")
	}

	argument := t.Args[0]
	// Fetch the argument value from the JSON
	value := gjson.GetBytes(t.Json, argument)
	if !value.Exists() {
		return nil, fmt.Errorf("argument '%s' not found in JSON", argument)
	}

	// Convert the value to uppercase
	return Results{strings.ToUpper(value.String())}, nil
}

// Concatenate struct holds the arguments and JSON data for transformation
type Concatenate struct {
	Config
}

// Transform concatenates two or more strings
func (t *Concatenate) Transform() (Results, error) {
	if len(t.Args) < 3 {
		return nil, fmt.Errorf("concatenate requires at least three arguments")
	}

	// Fetch the argument values from the JSON
	var values []string
	separator := strings.Trim(t.Args[0], "'")
	for _, arg := range t.Args[1:] {
		value := gjson.GetBytes(t.Json, arg)
		if !value.Exists() {
			return nil, fmt.Errorf("argument '%s' not found in JSON", arg)
		}
		values = append(values, value.String())
	}

	// Concatenate the values
	return Results{strings.Join(values, separator)}, nil
}

type Split struct {
	Config
}

// split takes one argument which is a field from json and a literal string
func (t *Split) Transform() (Results, error) {
	if len(t.Args) != 2 {
		return nil, fmt.Errorf("split requires exactly two arguments")
	}

	field := t.Args[0]
	separator := t.Args[1]
	// Fetch the argument value from the JSON
	value := gjson.GetBytes(t.Json, field)
	if !value.Exists() {
		return nil, fmt.Errorf("argument '%s' not found in JSON", field)
	}

	// Split the value
	splitValue := strings.Split(value.String(), strings.Trim(separator, "'"))
	results := Results{}
	for _, v := range splitValue {
		results = append(results, v)
	}
	return results, nil

}
