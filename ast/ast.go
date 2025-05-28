package ast 

//abstract syntax tree0 certain details visible in the source code are omitted in the AST
// ;, /n, whitespace, comments, braces, bracket and parentheses are not included
// the AST is a tree of nodes, each node is an instance of a type that implements the Node interface

// recursive descent parsern - top dowwn operator precedence 
// not care about speed yet, no dectection of erroneous syntax

// Every node in AST have to implement the Node interface -> provide a Tokenliteral() method
// that return the literal value of the token it's associate with 
import "monkey/token"

type Node interface {
	TokenLiteral() string
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
	Statements []Statement
}


type LetStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}


func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal}






func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

