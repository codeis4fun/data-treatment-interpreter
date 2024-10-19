package lexer

import (
	"fmt"
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
	NEWLINE    TokenType = "NEWLINE"
	EOF        TokenType = "EOF"
)

// Token represents a token with a type and literal value
type Token struct {
	Type    TokenType
	Literal string
	Pos     int // Position in the input string
}

// Lexer represents the state of the lexer
type Lexer struct {
	input  string
	start  int
	pos    int
	width  int
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
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.Run() // Start the lexer in a goroutine
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
	}
	l.start = l.pos
}

// run runs the state machine for lexing
func (l *Lexer) Run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
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

// nextRune returns the next rune and a boolean indicating if it's valid
func nextRune(l *Lexer) (rune, bool) {
	if l.pos >= len(l.input) {
		return -1, false
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r, true
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
		r, ok := nextRune(l)
		if !ok {
			l.emit(EOF)
			return nil
		}

		switch {
		case r == '\n':
			l.emit(NEWLINE) // Emit NEWLINE token for line breaks
			continue
		case r == '\'':
			return lexString // Handle string literals
		case unicode.IsSpace(r):
			l.start = l.pos // Skip whitespace
			continue
		case unicode.IsLetter(r) || r == '_':
			l.backup()
			return lexIdentifierOrKeyword
		case symbols[r] != "": // symbols[r] returns a tokentype
			l.emit(symbols[r])
		default:
			// Emit an ERROR token with more context
			l.emitError(fmt.Sprintf("unexpected character '%c'", r))

		}
	}
}

// emitError emits an ERROR token with the given error message
func (l *Lexer) emitError(message string) {
	l.tokens <- Token{
		Type:    ERROR,
		Literal: message,
		Pos:     l.start, // Track the position where the error occurred
	}
	l.start = l.pos
}

func isKeyword(word string) bool {
	_, found := keywords[word]
	return found
}

func isLetterOrDigit(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_'
}

// lexIdentifierOrKeyword scans identifiers (variable names and function names) and checks for keywords
func lexIdentifierOrKeyword(l *Lexer) stateFn {
	for r := l.next(); isLetterOrDigit(r); r = l.next() {
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
