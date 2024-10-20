package lexer_test

import (
	"testing"

	"github.com/codeis4fun/data-treatment-interpreter/internal/lexer"
)

func TestLexer(t *testing.T) {
	inputCmd := `SET bmi, isHealthy = bmi(weight, height)`

	l := lexer.NewLexer(inputCmd)

	expectedTokens := []lexer.Token{
		{Type: lexer.KEYWORD, Literal: "SET", Pos: 0},
		{Type: lexer.IDENTIFIER, Literal: "bmi", Pos: 4},
		{Type: lexer.COMMA, Literal: ",", Pos: 7},
		{Type: lexer.IDENTIFIER, Literal: "isHealthy", Pos: 9},
		{Type: lexer.OPERATOR, Literal: "=", Pos: 19},
		{Type: lexer.IDENTIFIER, Literal: "bmi", Pos: 21},
		{Type: lexer.LPAREN, Literal: "(", Pos: 24},
		{Type: lexer.IDENTIFIER, Literal: "weight", Pos: 25},
		{Type: lexer.COMMA, Literal: ",", Pos: 31},
		{Type: lexer.IDENTIFIER, Literal: "height", Pos: 33},
		{Type: lexer.RPAREN, Literal: ")", Pos: 39},
		{Type: lexer.EOF, Literal: "", Pos: 40},
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
		{Type: lexer.KEYWORD, Literal: "SET", Pos: 0},
		{Type: lexer.IDENTIFIER, Literal: "name", Pos: 4},
		{Type: lexer.OPERATOR, Literal: "=", Pos: 9},
		{Type: lexer.ERROR, Literal: "unexpected character '@'", Pos: 11}, // Invalid token should trigger an ERROR token
		{Type: lexer.ERROR, Literal: "unexpected character '1'", Pos: 12}, // Invalid token should trigger an ERROR token
		{Type: lexer.ERROR, Literal: "unexpected character '2'", Pos: 13}, // Invalid token should trigger an ERROR token
		{Type: lexer.ERROR, Literal: "unexpected character '3'", Pos: 14},
		{Type: lexer.IDENTIFIER, Literal: "invalid", Pos: 15},
		{Type: lexer.EOF, Literal: "", Pos: 22},
	}

	l := lexer.NewLexer(input)

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

	l := lexer.NewLexer(input)

	expectedTokens := []lexer.Token{
		{Type: lexer.KEYWORD, Literal: "SET", Pos: 0},
		{Type: lexer.IDENTIFIER, Literal: "languages.#", Pos: 4},
		{Type: lexer.OPERATOR, Literal: "=", Pos: 16},
		{Type: lexer.IDENTIFIER, Literal: "uppercase", Pos: 18},
		{Type: lexer.LPAREN, Literal: "(", Pos: 27},
		{Type: lexer.IDENTIFIER, Literal: "languages.#", Pos: 28},
		{Type: lexer.RPAREN, Literal: ")", Pos: 39},
		{Type: lexer.EOF, Literal: "", Pos: 40},
	}

	for _, expectedToken := range expectedTokens {
		actualToken := l.NextToken()
		if actualToken != expectedToken {
			t.Errorf("Expected token %v, got %v", expectedToken, actualToken)
		}
	}
}
