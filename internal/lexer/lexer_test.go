package lexer_test

import (
	"strings"
	"testing"

	"github.com/codeis4fun/data-treatment-interpreter/internal/lexer"
)

func TestLexer(t *testing.T) {
	inputCmd := `SET bmi, isHealthy = bmi(weight, height)`
	r := strings.NewReader(inputCmd)
	l := lexer.NewLexer(r)

	expectedTokens := []lexer.Token{
		{Type: lexer.KEYWORD, Literal: "SET", Line: 1, Pos: 0},
		{Type: lexer.IDENTIFIER, Literal: "bmi", Line: 1, Pos: 4},
		{Type: lexer.COMMA, Literal: ",", Line: 1, Pos: 7},
		{Type: lexer.IDENTIFIER, Literal: "isHealthy", Line: 1, Pos: 9},
		{Type: lexer.OPERATOR, Literal: "=", Line: 1, Pos: 19},
		{Type: lexer.IDENTIFIER, Literal: "bmi", Line: 1, Pos: 21},
		{Type: lexer.LPAREN, Literal: "(", Line: 1, Pos: 24},
		{Type: lexer.IDENTIFIER, Literal: "weight", Line: 1, Pos: 25},
		{Type: lexer.COMMA, Literal: ",", Line: 1, Pos: 31},
		{Type: lexer.IDENTIFIER, Literal: "height", Line: 1, Pos: 33},
		{Type: lexer.RPAREN, Literal: ")", Line: 1, Pos: 39},
		{Type: lexer.EOL, Literal: "\n", Line: 1, Pos: 40},
		{Type: lexer.EOF, Literal: "", Line: 2, Pos: 41},
	}

	for _, expectedToken := range expectedTokens {
		actualToken := l.NextToken()
		if actualToken != expectedToken {
			t.Errorf("Expected token %v, got %v", expectedToken, actualToken)
		}
	}
}

func TestUnexpectedToken(t *testing.T) {
	input := "SET name = @123invalid"
	expectedTokens := []lexer.Token{
		{Type: lexer.KEYWORD, Literal: "SET", Line: 1, Pos: 0},
		{Type: lexer.IDENTIFIER, Literal: "name", Line: 1, Pos: 4},
		{Type: lexer.OPERATOR, Literal: "=", Line: 1, Pos: 9},
		{Type: lexer.ERROR, Literal: "unexpected character '@'", Line: 1, Pos: 11}, // Invalid token should trigger an ERROR token
		{Type: lexer.ERROR, Literal: "unexpected character '1'", Line: 1, Pos: 12}, // Invalid token should trigger an ERROR token
		{Type: lexer.ERROR, Literal: "unexpected character '2'", Line: 1, Pos: 13}, // Invalid token should trigger an ERROR token
		{Type: lexer.ERROR, Literal: "unexpected character '3'", Line: 1, Pos: 14},
		{Type: lexer.IDENTIFIER, Literal: "invalid", Line: 1, Pos: 15},
		{Type: lexer.EOL, Literal: "\n", Line: 1, Pos: 22},
		{Type: lexer.EOF, Literal: "", Line: 2, Pos: 23},
	}

	r := strings.NewReader(input)
	l := lexer.NewLexer(r)

	for _, expected := range expectedTokens {
		token := l.NextToken()

		if token.Type != expected.Type {
			t.Fatalf("TestUnexpectedToken: Expected token type %v but got %v", expected.Type, token.Type)
		}

		if token.Literal != expected.Literal {
			t.Fatalf("TestUnexpectedToken: Expected token literal %v but got %v", expected.Literal, token.Literal)
		}

		if token.Pos != expected.Pos {
			t.Fatalf("TestUnexpectedToken: Expected token position %v but got %v", expected.Pos, token.Pos)
		}
	}
}

func TestLexerWithIterations(t *testing.T) {
	input := `SET languages.# = uppercase(languages.#)`
	r := strings.NewReader(input)
	l := lexer.NewLexer(r)

	expectedTokens := []lexer.Token{
		{Type: lexer.KEYWORD, Literal: "SET", Line: 1, Pos: 0},
		{Type: lexer.IDENTIFIER, Literal: "languages.#", Line: 1, Pos: 4},
		{Type: lexer.OPERATOR, Literal: "=", Line: 1, Pos: 16},
		{Type: lexer.IDENTIFIER, Literal: "uppercase", Line: 1, Pos: 18},
		{Type: lexer.LPAREN, Literal: "(", Line: 1, Pos: 27},
		{Type: lexer.IDENTIFIER, Literal: "languages.#", Line: 1, Pos: 28},
		{Type: lexer.RPAREN, Literal: ")", Line: 1, Pos: 39},
		{Type: lexer.EOL, Literal: "\n", Line: 1, Pos: 40},
		{Type: lexer.EOF, Literal: "", Line: 2, Pos: 41},
	}

	for _, expectedToken := range expectedTokens {
		actualToken := l.NextToken()
		if actualToken != expectedToken {
			t.Errorf("Expected token %v, got %v", expectedToken, actualToken)
		}
	}
}

func TestLexerWithTwoPrograms(t *testing.T) {
	input := `SET a = t(b, c)
SET d = t(e, f)`
	r := strings.NewReader(input)
	l := lexer.NewLexer(r)

	expectedTokens := []lexer.Token{
		{Type: lexer.KEYWORD, Literal: "SET", Line: 1, Pos: 0},
		{Type: lexer.IDENTIFIER, Literal: "a", Line: 1, Pos: 4},
		{Type: lexer.OPERATOR, Literal: "=", Line: 1, Pos: 6},
		{Type: lexer.IDENTIFIER, Literal: "t", Line: 1, Pos: 8},
		{Type: lexer.LPAREN, Literal: "(", Line: 1, Pos: 9},
		{Type: lexer.IDENTIFIER, Literal: "b", Line: 1, Pos: 10},
		{Type: lexer.COMMA, Literal: ",", Line: 1, Pos: 11},
		{Type: lexer.IDENTIFIER, Literal: "c", Line: 1, Pos: 13},
		{Type: lexer.RPAREN, Literal: ")", Line: 1, Pos: 14},
		{Type: lexer.EOL, Literal: "\n", Line: 1, Pos: 15},
		{Type: lexer.KEYWORD, Literal: "SET", Line: 2, Pos: 0},
		{Type: lexer.IDENTIFIER, Literal: "d", Line: 2, Pos: 4},
		{Type: lexer.OPERATOR, Literal: "=", Line: 2, Pos: 6},
		{Type: lexer.IDENTIFIER, Literal: "t", Line: 2, Pos: 8},
		{Type: lexer.LPAREN, Literal: "(", Line: 2, Pos: 9},
		{Type: lexer.IDENTIFIER, Literal: "e", Line: 2, Pos: 10},
		{Type: lexer.COMMA, Literal: ",", Line: 2, Pos: 11},
		{Type: lexer.IDENTIFIER, Literal: "f", Line: 2, Pos: 13},
		{Type: lexer.RPAREN, Literal: ")", Line: 2, Pos: 14},
		{Type: lexer.EOL, Literal: "\n", Line: 2, Pos: 15},
		{Type: lexer.EOF, Literal: "", Line: 3, Pos: 16},
	}

	for _, expectedToken := range expectedTokens {
		actualToken := l.NextToken()
		if actualToken != expectedToken {
			t.Errorf("Expected token %v, got %v", expectedToken, actualToken)
		}
	}
}
