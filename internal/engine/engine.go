package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codeis4fun/data-treatment-interpreter/internal/parser"
	"github.com/codeis4fun/data-treatment-interpreter/internal/transformers"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Transformer interface {
	Transform() (transformers.Results, error)
}

type transformerConfig func(config transformers.Config) Transformer

// Engine struct that manages transformers
type Engine struct {
	transformers map[string]transformerConfig
}

// NewEngine initializes the engine
func NewEngine() *Engine {
	transformers := map[string]transformerConfig{
		"uppercase":   func(config transformers.Config) Transformer { return &transformers.Uppercase{Config: config} },
		"concatenate": func(config transformers.Config) Transformer { return &transformers.Concatenate{Config: config} },
		"bmi":         func(config transformers.Config) Transformer { return &transformers.BMI{Config: config} },
		"split":       func(config transformers.Config) Transformer { return &transformers.Split{Config: config} },
	}
	return &Engine{
		transformers: transformers,
	}
}

// Execute applies the transformations defined in the Program struct to the input JSON
func (e *Engine) Execute(program *parser.Program, jsonData []byte) ([]byte, error) {
	// Check if the command starts with an iteration keyword
	if strings.Contains(program.Variables[0], "#") {
		return e.executeIteration(program, jsonData)
	}
	return e.executeSet(program, jsonData)
}

func (e *Engine) executeSet(program *parser.Program, jsonData []byte) ([]byte, error) {
	// Create the appropriate transformer based on the program
	transformerFunc, ok := e.transformers[program.Transformer]
	if !ok {
		return jsonData, fmt.Errorf("transformer '%s' not found", program.Transformer)
	}
	transformer := transformerFunc(transformers.Config{Args: program.Args, Json: jsonData})

	// Apply the transformation, get multiple outputs
	transformedValues, err := transformer.Transform()
	if err != nil {
		return nil, err
	}

	// Ensure the number of output values matches the number of variables in the program
	if len(transformedValues) != len(program.Variables) {
		return nil, fmt.Errorf("number of output values does not match the number of variables returned by transformer")
	}

	// Update the JSON with transformed values
	for i, value := range transformedValues {
		jsonData, err = sjson.SetBytes(jsonData, program.Variables[i], value)
		if err != nil {
			return nil, err
		}
	}

	return jsonData, nil
}

// Execute the iteration command
func (e *Engine) executeIteration(program *parser.Program, jsonData []byte) ([]byte, error) {
	// Gets the index of the placeholder in the variable path (e.g., "friends.#.first" -> 8)
	variable := program.Variables[0]
	placeholderIndex := strings.Index(variable, "#")

	// Extract the array field to iterate over (e.g., "friends.#.first" -> "friends")
	arrayField := variable[:placeholderIndex-1]
	array := gjson.GetBytes(jsonData, arrayField)

	// Check if the field is a valid array
	if !array.IsArray() {
		return nil, fmt.Errorf("field '%s' is not an array", arrayField)
	}

	// Iterate over each element in the array
	for i := range array.Array() {
		// Replace `#` in the variable path with the current index (e.g., "friends.#.first" -> "friends.0.first")
		variable := strings.Replace(variable, "#", strconv.Itoa(i), 1)

		// Apply the transformation to the current element
		transformerFunc, ok := e.transformers[program.Transformer]
		if !ok {
			return jsonData, fmt.Errorf("transformer '%s' not found", program.Transformer)
		}
		transformer := transformerFunc(transformers.Config{
			Args: []string{variable}, // Pass the current field to the transformer
			Json: jsonData,
		})

		// Get the transformed result (which should be for the current item only)
		transformedValues, err := transformer.Transform()
		if err != nil {
			return nil, err
		}

		// Update the JSON for the current array element
		for j, value := range transformedValues {
			// Use the transformed value for the current array element
			jsonData, err = sjson.SetBytes(jsonData, strings.Replace(program.Variables[j], "#", strconv.Itoa(i), 1), value)
			if err != nil {
				return nil, err
			}
		}
	}

	return jsonData, nil
}

// Execute multiple transformations in sequence
func (e *Engine) ExecuteAll(programs []*parser.Program, jsonData []byte) ([]byte, error) {
	var err error
	for _, program := range programs {
		// Execute each program (command) in sequence
		jsonData, err = e.Execute(program, jsonData)
		if err != nil {
			return nil, err
		}
	}
	// delete temporary variables from JSON which start with _prefix
	for _, program := range programs {
		for _, variable := range program.Variables {
			if strings.HasPrefix(variable, "_") {
				jsonData, err = sjson.DeleteBytes(jsonData, variable)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return jsonData, nil
}
