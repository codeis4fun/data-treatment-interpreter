package parser

import (
	"fmt"
	"strings"

	"github.com/codeis4fun/data-treatment-interpreter/internal/lexer" // Replace with the actual import path of your lexer package
)

// Program struct holds the parsed program information
type Program struct {
	Variables   []string // Variables being assigned
	Transformer string   // The transformation function
	Args        []string // Arguments to the transformation
}

// Parser struct, which wraps the lexer and consumes tokens
type Parser struct {
	lexer  *lexer.Lexer
	buffer []lexer.Token // Buffer to allow peeking tokens
	input  string        // Store the input string for error reporting
}

// NewParser initializes a new parser with the given lexer and input
func NewParser(l *lexer.Lexer, input string) *Parser {
	return &Parser{
		lexer:  l,
		buffer: []lexer.Token{},
		input:  input,
	}
}

// nextToken fetches the next token, considering the buffer
func (p *Parser) nextToken() lexer.Token {
	if len(p.buffer) > 0 {
		token := p.buffer[0]
		p.buffer = p.buffer[:0] // Clear the buffer after consuming
		return token
	}
	return p.lexer.NextToken()
}

// peekToken looks at the next token without consuming it
func (p *Parser) peekToken() lexer.Token {
	if len(p.buffer) > 0 {
		return p.buffer[0]
	}
	token := p.lexer.NextToken()
	p.buffer = append(p.buffer, token)
	return token
}

// Run starts the parser and processes tokens from the lexer
func (p *Parser) Run() (*Program, error) {
	// Process tokens as they are emitted by the lexer
	return p.parseProgram()
}

// Parse multiple commands
func (p *Parser) RunAll() ([]*Program, error) {
	var programs []*Program

	// Process multiple commands
	// label to break out of the loop when EOF is reached
parseLoop:
	for {
		program, err := p.parseProgram()
		if err != nil {
			return nil, err
		}

		if program != nil {
			programs = append(programs, program)
		}

		// Expect a NEWLINE or EOF to indicate the end of a command
		token := p.peekToken()
		switch token.Type {
		case lexer.NEWLINE:
			p.nextToken() // Consume the newline
		case lexer.EOF:
			break parseLoop // Stop if the end of the file is reached
		default:
			return nil, p.errorWithContext(token, "expected newline or end of file")
		}
	}

	return programs, nil
}

// parseProgram parses the input and returns a Program struct
func (p *Parser) parseProgram() (*Program, error) {
	// Expect 'SET' keyword
	token := p.nextToken()
	if token.Type != lexer.KEYWORD || token.Literal != "SET" {
		return nil, p.errorWithContext(token, "expected 'SET' keyword")
	}

	// Parse the rest of the program (variables, transformer, args)
	return p.parseAssignment()
}

// parseAssignment parses an assignment like: SET var1, var2 = transformer(arg1, arg2)
func (p *Parser) parseAssignment() (*Program, error) {
	// Parse variables
	variables, err := p.parseVariables()
	if err != nil {
		return nil, err
	}
	// Expect '='
	token := p.nextToken()
	if token.Type != lexer.OPERATOR || token.Literal != "=" {
		return nil, p.errorWithContext(token, fmt.Sprintf("expected operator '%s'", "="))
	}

	// Parse transformer and arguments
	transformer, args, err := p.parseTransformer()
	if err != nil {
		return nil, err
	}

	return &Program{
		Variables:   variables,
		Transformer: transformer,
		Args:        args,
	}, nil
}

// parseVariables parses the list of variables being assigned, including placeholders
func (p *Parser) parseVariables() ([]string, error) {
	var variables []string

	// Expect at least one identifier (variable name, which may include '#')
	firstVar := p.nextToken()
	if firstVar.Type != lexer.IDENTIFIER {
		return nil, p.errorWithContext(firstVar, "expected variable name")
	}
	variables = append(variables, firstVar.Literal)

	// Continue parsing more variables separated by commas, stop when encountering '='
parseLoop:
	for {
		switch nextToken := p.peekToken(); {
		case nextToken.Type == lexer.OPERATOR && nextToken.Literal == "=":
			break parseLoop

		case nextToken.Type == lexer.COMMA:
			p.nextToken() // Consume the comma
			nextVar := p.nextToken()
			if nextVar.Type != lexer.IDENTIFIER {
				return nil, p.errorWithContext(nextVar, "expected variable name")
			}
			variables = append(variables, nextVar.Literal)

		default:
			return nil, p.errorWithContext(nextToken, "unexpected token in variables")
		}
	}

	return variables, nil
}

// parseTransformer parses the transformer function and its arguments
func (p *Parser) parseTransformer() (string, []string, error) {
	// Expect the transformer name (an identifier)
	transformer := p.nextToken()
	if transformer.Type != lexer.IDENTIFIER {
		return "", nil, p.errorWithContext(transformer, "expected transformer name")
	}

	// Expect '(' to start argument list
	if err := p.expectSymbol(lexer.LPAREN); err != nil {
		return "", nil, err
	}

	// Parse the arguments
	var args []string
parseLoop:
	for {
		nextToken := p.nextToken()
		if nextToken.Type == lexer.RPAREN {
			break
		}
		if nextToken.Type != lexer.IDENTIFIER && nextToken.Type != lexer.STRING {
			return "", nil, p.errorWithContext(nextToken, "expected argument (field or string)")
		}
		args = append(args, nextToken.Literal)

		// Handle commas between arguments
		nextToken = p.peekToken()
		switch nextToken.Type {
		case lexer.COMMA:
			p.nextToken() // Consume the comma
		case lexer.RPAREN:
			p.nextToken() // Consume the closing parenthesis
			break parseLoop
		default:
			return "", nil, p.errorWithContext(nextToken, "unexpected token in arguments")
		}
	}

	return transformer.Literal, args, nil
}

// expectSymbol checks if the next token is the expected symbol type
func (p *Parser) expectSymbol(expectedType lexer.TokenType) error {
	token := p.nextToken()
	if token.Type != expectedType {
		return p.errorWithContext(token, fmt.Sprintf("expected symbol '%s'", expectedType))
	}
	return nil
}

// errorWithContext provides an error message with context and highlights where the error occurred
func (p *Parser) errorWithContext(tok lexer.Token, message string) error {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s at position %d: %s", message, tok.Pos, tok.Literal)) // Include position in error
	builder.WriteString("\n")
	builder.WriteString(p.input)
	builder.WriteString("\n")
	builder.WriteString(p.makePointer(tok.Pos))
	return fmt.Errorf(builder.String())
}

// makePointer creates a pointer string (e.g., "   ^") to show where the error occurred
func (p *Parser) makePointer(pos int) string {
	pointer := make([]rune, pos)
	for i := range pointer {
		pointer[i] = ' '
	}
	return string(pointer) + "^"
}
