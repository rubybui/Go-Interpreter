package parser 

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"fmt"
)

// Parser represents a parser for the Monkey programming language.
// It holds the lexer and current/peek tokens for parsing.
type Parser struct {
	lexer *lexer.Lexer
	currentToken token.Token
	peekToken token.Token
	errors []string
}

// New creates a new Parser instance with the given lexer.
// It initializes the parser by reading the first two tokens.
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, errors: []string{}}

	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances the parser by setting the current token to the peek token
// and reading the next token from the lexer.
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}
// ParseProgram parses the entire program and returns an AST.
// It iterates through all tokens until EOF, parsing each statement
// and adding it to the program's statement list.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement determines the type of statement to parse based on the current token.
// Currently only handles LET statements, returns nil for other token types.
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// parseLetStatement parses a let statement in the format: let <identifier> = <expression>;
// It constructs an AST node for the let statement, validates the identifier,
// and expects an equals sign. Currently skips the expression parsing.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO: parse expression

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatment parses a return statement in the format: return <expression>;
// It constructs an AST  node for the return statement
// Currently skips the expression parsing.
func  (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()


	//TODO: parse expression

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt

}

// curTokenIs checks if the current token is of the specified type.
// Helper method for token type checking.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// peekTokenIs checks if the next token is of the specified type.
// Helper method for looking ahead at the next token.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek is an assertion function that checks if the next token is of the expected type.
// If it is, it advances the tokens and returns true.
// If not, it returns false, allowing the parser to handle the error.
// This is a crucial method for enforcing the correct order of tokens in the input.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}