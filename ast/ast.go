package ast 

import (
	"monkey/token"
	"bytes"
	"strings"
)

// Node represents a node in the Abstract Syntax Tree.
// Every node in AST must implement the Node interface by providing:
// - TokenLiteral(): returns the literal value of the token it's associated with
// - String(): returns a string representation of the node for debugging
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement node in the AST.
// It must implement both Node interface and provide a statementNode() method
// to distinguish it from expressions.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node in the AST.
// It must implement both Node interface and provide an expressionNode() method
// to distinguish it from statements.
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of the AST.
// It contains a slice of statements that make up the program.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the literal value of the first statement's token.
// If there are no statements, it returns an empty string.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns a string representation of the program by concatenating
// the string representation of all its statements.
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// LetStatement represents a variable declaration statement.
// It contains:
// - Token: the 'let' token
// - Name: the identifier being declared
// - Value: the expression being assigned to the identifier
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// String returns a string representation of the let statement in the format:
// "let <identifier> = <expression>;"
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents a return statement.
// It contains:
// - Token: the 'return' token
// - ReturnValue: the expression being returned
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns a string representation of the return statement in the format:
// "return <expression>;"
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents a statement that consists of a single expression.
// It contains:
// - Token: the first token of the expression
// - Expression: the expression itself
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns a string representation of the expression statement.
// If the expression is nil, returns an empty string.
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Identifier represents an identifier expression.
// It contains:
// - Token: the identifier token
// - Value: the name of the identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String returns the identifier's value as a string.
func (i *Identifier) String() string {
	return i.Value
}

// IntegerLiteral represents an integer literal expression.
// It contains:
// - Token: the integer token
// - Value: the integer value
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns the integer literal's value as a string.
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// PrefixExpression represents a prefix operator expression (e.g., !true, -5).
// It contains:
// - Token: the prefix operator token
// - Operator: the operator string
// - Right: the expression being operated on
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns a string representation of the prefix expression in the format:
// "(<operator><right>)"
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents an infix operator expression (e.g., 5 + 5, true == false).
// It contains:
// - Token: the operator token
// - Left: the left-hand expression
// - Operator: the operator string
// - Right: the right-hand expression
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a string representation of the infix expression in the format:
// "(<left> <operator> <right>)"
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// Boolean represents a boolean literal expression.
// It contains:
// - Token: the boolean token
// - Value: the boolean value
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// String returns the boolean literal's value as a string.
func (b *Boolean) String() string {
	return b.Token.Literal
}

// IfExpression represents an if-else expression.
// It contains:
// - Token: the 'if' token
// - Condition: the condition expression
// - Consequence: the block to execute if condition is true
// - Alternative: the block to execute if condition is false (optional)
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a string representation of the if expression in the format:
// "if <condition> { <consequence> } else { <alternative> }"
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// BlockStatement represents a block of statements enclosed in curly braces.
// It contains:
// - Token: the opening brace token
// - Statements: the list of statements in the block
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// String returns a string representation of the block statement by concatenating
// the string representation of all its statements.
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// FunctionLiteral represents a function literal expression.
// It contains:
// - Token: the 'fn' token
// - Parameters: list of parameter identifiers
// - Body: the function body as a block statement
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns a string representation of the function literal in the format:
// "fn(<param1>, <param2>, ...) { <body> }"
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression represents a function call expression.
// It contains:
// - Token: the opening parenthesis token
// - Function: the function being called (Identifier or FunctionLiteral)
// - Arguments: list of argument expressions
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// String returns a string representation of the function call in the format:
// "<function>(<arg1>, <arg2>, ...)"
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}