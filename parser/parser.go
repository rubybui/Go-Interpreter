package parser 

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"fmt"
	"strconv"
)

// Precedence levels for operator precedence parsing
const (
	_ int = iota
	LOWEST      // Lowest precedence
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// Parser represents a parser for the Monkey programming language.
// It holds:
// - lexer: the lexical analyzer
// - currentToken: the current token being processed
// - peekToken: the next token to be processed
// - errors: list of parsing errors
// - prefixParseFns: map of prefix parsing functions
// - infixParseFns: map of infix parsing functions
type Parser struct {
	lexer          *lexer.Lexer
	currentToken   token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}	

// New creates a new Parser instance with the given lexer.
// It initializes the parser by:
// 1. Creating prefix and infix parse function maps
// 2. Registering parsing functions for different token types
// 3. Reading the first two tokens
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, errors: []string{}}

	// Initialize prefix parse functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// Initialize infix parse functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)


	// Read first two tokens
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

// Errors returns the list of parsing errors.
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
// It handles:
// - LET statements
// - RETURN statements
// - Expression statements
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses a let statement in the format: let <identifier> = <expression>;
// It:
// 1. Creates a new LetStatement node
// 2. Validates the identifier
// 3. Expects an equals sign
// 4. Parses the expression
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// Skip until semicolon for now
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement in the format: return <expression>;
// It:
// 1. Creates a new ReturnStatement node
// 2. Advances to the next token
// 3. Skips until semicolon for now
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	// Skip until semicolon for now
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// noPrefixParseFnError adds an error when no prefix parse function is found
// for the given token type.
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parseExpression parses an expression with the given precedence.
// It:
// 1. Gets the prefix parse function for the current token
// 2. Parses the left-hand side of the expression
// 3. Continues parsing while the next token has higher precedence
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// parseExpressionStatement parses a statement that consists of a single expression.
// It:
// 1. Creates a new ExpressionStatement node
// 2. Parses the expression with lowest precedence
// 3. Optionally consumes a semicolon
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// curTokenIs checks if the current token is of the specified type.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// peekTokenIs checks if the next token is of the specified type.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek is an assertion function that checks if the next token is of the expected type.
// If it is, it advances the tokens and returns true.
// If not, it adds an error and returns false.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekError adds an error when the next token is not of the expected type.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Type definitions for parse functions
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// registerPrefix registers a prefix parse function for a token type.
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers an infix parse function for a token type.
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// parseIdentifier creates an Identifier node for the current token.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

// parseIntegerLiteral creates an IntegerLiteral node for the current token.
// It parses the token's literal value as an int64.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// parsePrefixExpression creates a PrefixExpression node for the current token.
// It:
// 1. Creates the node with the current token and operator
// 2. Advances to the next token
// 3. Parses the right-hand expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// peekPrecedence returns the precedence of the next token.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence returns the precedence of the current token.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

// parseInfixExpression creates an InfixExpression node for the current token.
// It:
// 1. Creates the node with the current token and operator
// 2. Stores the left-hand expression
// 3. Gets the precedence of the current operator
// 4. Advances to the next token
// 5. Parses the right-hand expression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseBoolean creates a Boolean node for the current token.
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

// parseGroupedExpression parses an expression enclosed in parentheses.
// It:
// 1. Advances past the opening parenthesis
// 2. Parses the expression
// 3. Expects a closing parenthesis
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseIfExpression parses an if expression.
// It:
// 1. Creates an IfExpression node
// 2. Expects an opening parenthesis
// 3. Parses the condition
// 4. Expects a closing parenthesis
// 5. Expects an opening brace
// 6. Parses the consequence block
// 7. Optionally parses an else block
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// parseBlockStatement parses a block of statements enclosed in curly braces.
// It:
// 1. Creates a BlockStatement node
// 2. Advances past the opening brace
// 3. Parses statements until a closing brace or EOF
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseFunctionParameters parses the parameters of a function literal.
// It:
// 1. Handles empty parameter lists
// 2. Parses the first parameter
// 3. Parses additional parameters separated by commas
// 4. Expects a closing parenthesis
// Returns a slice of parameter identifiers or nil if parsing fails.
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseFunctionLiteral parses a function literal expression.
// It:
// 1. Creates a FunctionLiteral node
// 2. Expects an opening parenthesis
// 3. Parses the parameter list through parseFunctionParameters
// 4. Expects an opening brace through block statement
// 5. Parses the function body as a block statement
// Returns the parsed FunctionLiteral or nil if parsing fails.
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()
	return lit
}

// precedences maps token types to their precedence levels.
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}
