package engine

import (
	"fmt"

	"github.com/codeis4fun/data-treatment-interpreter/internal/parser"
	"github.com/codeis4fun/data-treatment-interpreter/internal/transformers"
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
	// Create the appropriate transformer based on the program
	transformerFunc, ok := e.transformers[program.Transformer]
	if !ok {
		return nil, fmt.Errorf("transformer '%s' not found", program.Transformer)
	}
	transformer := transformerFunc(transformers.Config{Args: program.Args, Json: jsonData})

	// Apply the transformation, get multiple outputs
	transformedValues, err := transformer.Transform()
	if err != nil {
		return nil, err
	}

	// Ensure the number of output values matches the number of variables
	if len(transformedValues) != len(program.Variables) {
		return nil, fmt.Errorf("number of output values does not match the number of variables")
	}

	for i, value := range transformedValues {
		jsonData, err = sjson.SetBytes(jsonData, program.Variables[i], value)
		if err != nil {
			return nil, err
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
			if variable[0] == '_' {
				jsonData, err = sjson.DeleteBytes(jsonData, variable)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return jsonData, nil
}
