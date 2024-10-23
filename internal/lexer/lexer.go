package lexer

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// TokenType represents different types of tokens in the language
type TokenType string

const (
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	OPERATOR   TokenType = "OPERATOR"
	LPAREN     TokenType = "LPAREN"
	RPAREN     TokenType = "RPAREN"
	COMMA      TokenType = "COMMA"
	KEYWORD    TokenType = "KEYWORD"
	ERROR      TokenType = "ERROR"
	EOL        TokenType = "EOL" // End of line token
	EOF        TokenType = "EOF" // End of file token
)

// Token represents a token with a type and literal value
type Token struct {
	Type    TokenType
	Literal string
	Line    int // Line number where the token was found
	Pos     int // Position in the input string
}

// Lexer represents the state of the lexer
type Lexer struct {
	sc     *bufio.Scanner
	input  string
	start  int
	pos    int
	width  int
	line   int
	tokens chan Token
}

type stateFn func(*Lexer) stateFn

var keywords = map[string]TokenType{
	"SET": KEYWORD,
}

// Symbols table to handle operators and punctuation
var symbols = map[rune]TokenType{
	'=': OPERATOR,
	'(': LPAREN,
	')': RPAREN,
	',': COMMA,
}

// NewLexer initializes a new lexer
func NewLexer(r io.Reader) *Lexer {
	l := &Lexer{
		sc:     bufio.NewScanner(r),
		line:   1, // Start on the first line
		tokens: make(chan Token),
	}
	go l.run() // Start the lexer in a goroutine
	return l
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	return <-l.tokens
}

// Emit sends a token to the tokens channel
func (l *Lexer) emit(t TokenType) {
	l.tokens <- Token{
		Type:    t,
		Literal: l.input[l.start:l.pos],
		Pos:     l.start, // Save the position of the token
		Line:    l.line,  // Track the line number
	}
	l.start = l.pos
}

// run runs the state machine for lexing
func (l *Lexer) run() {
	for l.sc.Scan() {
		// Set the current line as input
		l.input = l.sc.Text() + "\n"
		l.pos = 0
		l.start = 0

		// Process the line by running the state machine
		for state := lexText; state != nil; {
			state = state(l)
		}
	}

	// Send EOF token when the input is completely done
	l.emit(EOF)
	close(l.tokens)

	if err := l.sc.Err(); err != nil {
		l.emitError(fmt.Sprintf("error reading input: %v", err))
	}
}

// next returns the next rune in the input and advances the position
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return -1
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// backup steps back one rune
func (l *Lexer) backup() {
	l.pos -= l.width
}

// lexString scans string literals (enclosed in single quotes)
func lexString(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == '\'' {
			l.emit(STRING)
			return lexText
		}
		if r == -1 {
			l.emit(ERROR) // Unterminated string
			return nil
		}
	}
}

// lexText is the main lexing state function for parsing identifiers and operators
func lexText(l *Lexer) stateFn {
	for {
		r := l.next()

		switch {
		case r == '\n':
			l.emit(EOL) // Emit EOL token for line breaks
			l.line++    // Increment the line number
			return nil  // Stop lexing the current line and wait for the next line
		case r == '\'':
			return lexString // Handle string literals
		case unicode.IsSpace(r):
			l.start = l.pos // Skip whitespace
			continue
		case unicode.IsLetter(r) || r == '_' || r == '#': // Allow '#' and '_' as part of identifiers
			l.backup()
			return lexIdentifierOrKeyword
		case symbols[r] != "": // symbols[r] returns the token type for the rune
			l.emit(symbols[r])
		case r == -1:
			l.emit(EOF)
		default:
			// Emit an ERROR token with more context
			l.emitError(fmt.Sprintf("unexpected character '%c'", r))
			return nil
		}
	}
}

// emitError emits an ERROR token with the given error message
func (l *Lexer) emitError(message string) {
	l.tokens <- Token{
		Type:    ERROR,
		Literal: message,
		Pos:     l.start, // Track the position where the error occurred
		Line:    l.line,  // Track the line number
	}
	l.start = l.pos
}

func isKeyword(word string) bool {
	_, found := keywords[word]
	return found
}

func isAllowedCharacter(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '#'
}

// lexIdentifierOrKeyword scans identifiers (variable names and function names) and checks for keywords
func lexIdentifierOrKeyword(l *Lexer) stateFn {
	for r := l.next(); isAllowedCharacter(r); r = l.next() {
		// Continue scanning if the character is a letter or digit
	}
	l.backup() // We've scanned one character too far; back up

	// Extract the scanned word
	word := l.input[l.start:l.pos]

	// Check if it's a keyword or identifier
	if isKeyword(word) {
		l.emit(KEYWORD)
		return lexText
	}

	// It's an identifier if it's not a keyword
	l.emit(IDENTIFIER)
	return lexText
}
