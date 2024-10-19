package transformers

import (
	"fmt"
	"math"

	"github.com/tidwall/gjson"
)

type BMI struct {
	Config
}

// Transform calculates the Body Mass Index (BMI) based on the weight and height
// Transform calculates BMI and checks if it's in a healthy range
func (t *BMI) Transform() (Results, error) {
	// BMI requires two arguments: weight and height
	if len(t.Config.Args) != 2 {
		return nil, fmt.Errorf("bmi requires exactly two arguments")
	}

	weightArg := t.Config.Args[0]
	heightArg := t.Config.Args[1]

	// Get weight and height from the JSON
	weightVar := gjson.GetBytes(t.Config.Json, weightArg)
	heightVar := gjson.GetBytes(t.Config.Json, heightArg)

	if !weightVar.Exists() || !heightVar.Exists() {
		return nil, fmt.Errorf("weight or height not found in JSON")
	}

	// check if the weight and height are numbers
	if weightVar.Type != gjson.Number || heightVar.Type != gjson.Number {
		return nil, fmt.Errorf("weight and height must be numbers")
	}

	// Parse weight and height as float64
	weight := weightVar.Float()
	height := heightVar.Float()

	// Calculate BMI with the formula: weight (kg) / height (m)^2 with height in cm and weight in kg with 2 decimal places
	bmi := math.Round((weight/(height*height))*10) / 10

	// Determine if the BMI is healthy (for simplicity, we'll assume healthy BMI is between 18.5 and 24.9)
	isHealthy := bmi >= 18.5 && bmi <= 24.9

	// Return both BMI and whether it's healthy
	return Results{bmi, isHealthy}, nil
}
