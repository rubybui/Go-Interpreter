package ast 

//abstract syntax tree certain details visible in the source code are omitted in the AST
// ;, /n, whitespace, comments, braces, bracket and parentheses are not included
// the AST is a tree of nodes, each node is an instance of a type that implements the Node interface

// recursive descent parsern - top dowwn operator precedence 
// not care about speed yet, no dectection of erroneous syntax

// Every node in AST have to implement the Node interface -> provide a Tokenliteral() method
// that return the literal value of the token it's associate with 
import (
	"monkey/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string // for easier printing of AST node 
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements 	[]Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
	out.WriteString(s.String())
	}
	return out.String()
}


type LetStatement struct {
	Token 		token.Token
	Name 		*Identifier
	Value 		Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal}
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


type ReturnStatement struct {
	Token 		token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}


type ExpressionStatement struct {
	Token 		token.Token
	Expression 	Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}


type Identifier struct {
	Token 		token.Token
	Value 		string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal}
func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token 		token.Token
	Value 		int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token 		token.Token // The prefix token, e.g. !
	Operator 	string
	Right 		Expression
}
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token 		token.Token // The prefix token, e.g. !
	Left		Expression
	Operator 	string
	Right 		Expression
}
func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}