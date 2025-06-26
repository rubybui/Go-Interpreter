//take source code as input and output the tokens that represent the source code.
// go through its input and output the next token it recognizes
// initialize the lexer with our source code
// repeatedly call NextToken() on it to go through the source code, token by token, character by character

// lexer/lexer.go
package lexer

import "monkey/token"

type Lexer struct {
    input        string
    position     int
    readPosition int
    ch           byte
}

func New(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token
    l.skipWhitespace()
    switch l.ch {
    case '=':
        if l.peekChar() == '=' {
            ch:= l.ch
            l.readChar()
            literal := string(ch) + string(l.ch)
            tok = token.Token{Type: token.EQ, Literal: literal}
        } else {
        tok = newToken(token.ASSIGN, l.ch)
        }
    case '"':
        tok.Type = token.STRING
        tok.Literal = l.readString()
    case ';':
        tok = newToken(token.SEMICOLON, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
    case '<':
        tok = newToken(token.LT, l.ch)
    case '>':
        tok = newToken(token.GT, l.ch)
    case '!':
        if l.peekChar() == '=' {
            ch := l.ch
            l.readChar()
            literal := string(ch) + string(l.ch)
            tok = token.Token{Type: token.NOT_EQ, Literal: literal}
        } else {
            tok = newToken(token.BANG, l.ch)
        }
    case '{':
        tok = newToken(token.LBRACE, l.ch)
    case '}':
        tok = newToken(token.RBRACE, l.ch)
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    default:    //recognize whether the current character is a letter and if so, it needs to read the rest of the identifier/keyword until it encounters a non-letter-character
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = token.LookupIdent(tok.Literal)
            return tok
        // The early exit here, our return tok statement, is necessary because when calling
        // readIdentifier(), we call readChar() repeatedly and advance our readPosition and position
        // fields past the last character of the current identifier. So we don't need the call to readChar()
        // after the switch statement again.
        } else if isDigit(l.ch) {
            tok.Literal = l.readNumber()
            tok.Type = token.INT
            return tok
        } else {
            tok = newToken(token.ILLEGAL, l.ch)
        }
    }
    l.readChar()
    return tok
}

func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition += 1
}

func isLetter(ch byte) bool { 
    // check if the character is an ascii letter, underscore, 
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
    }

func (l *Lexer) readNumber() string {
    position := l.position
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
    // peek at the immediately next character without advancing the lexer's position
    if l.readPosition >= len(l.input) {
        return 0
    }
    return l.input[l.readPosition]
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}
func newToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
    // skip over any whitespace characters in the input stream, skip newlines, tabs, and spaces
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) readString() string {
    position := l.position + 1
    for {
        l.readChar()
        if l.ch == '"' || l.ch == 0 {
            break
        }
    }
    return l.input[position:l.position]
}
