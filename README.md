# Go-Interpreter

A Monkey programming language interpreter written in Go. This project implements a lexer, parser, and evaluator for the Monkey programming language.

## Features

### Lexer
- Tokenizes source code into a stream of tokens
- Supports identifiers, integers, operators, and keywords
- Handles whitespace and special characters

### Parser
The parser implements a Pratt parser (top-down operator precedence parser) that can handle:

#### Expression Types
- Identifiers
- Integer literals
- Prefix expressions (e.g., `-5`, `!true`)
- Infix expressions with operator precedence:
  - Addition (`+`)
  - Subtraction (`-`)
  - Multiplication (`*`)
  - Division (`/`)
  - Comparison operators (`==`, `!=`, `<`, `>`)

#### Statement Types
- Let statements (`let x = 5;`)
- Return statements (`return 10;`)
- Expression statements

### Operator Precedence
The parser implements the following precedence levels (from lowest to highest):
1. `LOWEST`
2. `EQUALS` (`==`, `!=`)
3. `LESSGREATER` (`<`, `>`)
4. `SUM` (`+`, `-`)
5. `PRODUCT` (`*`, `/`)
6. `PREFIX` (`-x`, `!x`)
7. `CALL` (function calls)

## Project Structure
```
.
├── lexer/         # Lexical analysis
├── parser/        # Syntax analysis
├── ast/          # Abstract Syntax Tree nodes
└── token/        # Token definitions
```

## Testing
The project includes comprehensive test coverage for each package. To run the tests:

### Run all tests
```bash
go test ./...
```

### Run tests for specific packages
```bash
# Test lexer package
go test ./lexer

# Test parser package
go test ./parser

# Test AST package
go test ./ast
```

### Run tests with coverage
```bash
# Generate coverage report for all packages
go test ./... -cover

# Generate coverage report for a specific package
go test ./lexer -cover
```

### Run tests with verbose output
```bash
go test ./... -v
```

## Usage
[Usage instructions to be added as the project develops]

## Development Status
The interpreter is currently under development. The lexer and parser are implemented with basic functionality for expressions and statements.

## License
[License information to be added]
