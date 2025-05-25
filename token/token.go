package token

type TokenType string

type Token struct {
	Type    TokenType //  defined the TokenType type to be a string.  allows us to distinguish between different types of tokens
	Literal string
}

// we can define the possible TokenTypes as constants.
const (
	ILLEGAL TokenType = "ILLEGAL" // signifies a token/character we don’t know about
	EOF     TokenType = "EOF"     //end of file”

	// Identifiers + literals
	IDENT TokenType = "IDENT" // add, foobar, x, y, ...
	INT   TokenType = "INT"   // 1343456

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
)
