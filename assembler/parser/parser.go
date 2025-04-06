// mov $42, r1 -> (instruction literal hex value, register)

// move [$42 + (!loc - $1F)], r1 -> (instruction [literal hex value + (!variable - )])
//package parser

//import (
//	"github.com/alecthomas/participle/v2"
///	"github.com/alecthomas/participle/v2/lexer"
///)

// Register represents a register (e.g., R1, r2, SP)
//type Register struct {
//	Value string `parser:"@Register"`
//}

// HexLiteral represents a hex literal (e.g., $1F, $deadBEEF)
//type HexLiteral struct {
//	Value string `parser:"'$' @HexDigit+"`
//}

// SquareBracketExpr represents an expression in square brackets (e.g., [R1])
//type SquareBracketExpr struct {
//	Expr string `parser:"'[' @Ident ']'"`
//}

// Instruction represents a generic instruction
//type Instruction struct {
//	Instruction string        `parser:"@Ident"`
//	Args        []interface{} `parser:"@HexLiteral | @Register | @SquareBracketExpr"`
//	Type        string        `parser:""`
//}

// Custom lexer for both registers, hex literals, and square bracket expressions
//var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
///	{Name: "Register", Pattern: `(?i)\b(r[1-8]|sp|fp|ip|acc)\b`},
//	{Name: "HexDigit", Pattern: `[0-9A-Fa-f]`},         // hex digit
///	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`}, // identifier
//	{Name: "Char", Pattern: `\$`},                      // '$' literal
//	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},        // skip whitespace
///})

// RegisterParser returns a parser for register names
//func RegisterParser() (*participle.Parser[Register], error) {
//	return participle.Build[Register](
//		participle.Lexer(lexerDef),
//	)
//}

// HexParser returns a parser for hex literals
//func HexParser() (*participle.Parser[HexLiteral], error) {
//	return participle.Build[HexLiteral](
///		participle.Lexer(lexerDef),
//		participle.Elide("Whitespace"),
//	)
//}

// SquareBracketParser returns a parser for expressions in square brackets
//func SquareBracketParser() (*participle.Parser[SquareBracketExpr], error) {
//	return participle.Build[SquareBracketExpr](
//		participle.Lexer(lexerDef),
///		participle.Elide("Whitespace"),
//	)
//}

// MovLitToRegParser parses a MOV_LIT_REG instruction
//
//	func MovLitToRegParser() (*participle.Parser[Instruction], error) {
//		return participle.Build[Instruction](
//
// /		participle.Lexer(lexerDef),
//
//			participle.Elide("Whitespace"),
//		)
//	}/
package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Register represents a register (e.g., r1, r2, acc, etc.)
type Register struct {
	Value string `parser:"@Register"`
}

// HexLiteral represents a hexadecimal literal (e.g., $42, $deadBEEF)
type HexLiteral struct {
	Value string `parser:"'$' @HexDigit+"`
}

// Variable represents a variable (e.g., !myVar)
type Variable struct {
	Name string `parser:"'!' @Ident"`
}

// Operator represents an operator (e.g., +, -, *)
type Operator struct {
	Type string `parser:"@('+'|'-'|'*')"`
}

// SquareBracketExpr represents an expression in square brackets (e.g., [r1 + r2])
type SquareBracketExpr struct {
	Elements []interface{} `parser:"'[' (@HexLiteral | @Register | @Variable) ( @Operator (@HexLiteral | @Register | @Variable) )* ']'"`
}

// Instruction represents an instruction (e.g., MOV_LIT_REG)
type Instruction struct {
	Instruction string        `parser:"@Ident"`
	Args        []interface{} `parser:"@HexLiteral | @Register | @SquareBracketExpr"`
}

var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
	// Register (r1 to r8, sp, fp, ip, acc)
	{Name: "Register", Pattern: `(?i)\b(r[1-8]|sp|fp|ip|acc)\b`},
	// Hexadecimal digits
	{Name: "HexDigit", Pattern: `[0-9A-Fa-f]`},
	// Identifiers (e.g., variable names)
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	// Whitespace
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
})

// RegisterParser returns a parser for register names
func RegisterParser() (*participle.Parser[Register], error) {
	return participle.Build[Register](
		participle.Lexer(lexerDef),
	)
}

// HexParser returns a parser for hex literals
func HexParser() (*participle.Parser[HexLiteral], error) {
	return participle.Build[HexLiteral](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
}

// VariableParser returns a parser for variables
func VariableParser() (*participle.Parser[Variable], error) {
	return participle.Build[Variable](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
}

// SquareBracketParser parses expressions in square brackets
func SquareBracketParser() (*participle.Parser[SquareBracketExpr], error) {
	return participle.Build[SquareBracketExpr](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
}

// MovLitToRegParser parses a MOV_LIT_REG instruction
func MovLitToRegParser() (*participle.Parser[Instruction], error) {
	return participle.Build[Instruction](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
}
