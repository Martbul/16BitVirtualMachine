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

//package parser

//import (
//	"github.com/alecthomas/participle/v2"
//	"github.com/alecthomas/participle/v2/lexer"
//)

// === Lexer ===

//var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
//	{Name: "Register", Pattern: `(?i)\b(r[1-8]|sp|fp|ip|acc)\b`},
//	{Name: "HexDigit", Pattern: `[0-9A-Fa-f]+`},
//	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
//	{Name: "Operator", Pattern: `[\+\-\*]`},
//	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
//	{Name: "Punct", Pattern: `[\[\],!$]`},
//})

// === AST Nodes ===

//type Register struct {
//	Value string `parser:"@Register"`
//}

//type HexLiteral struct {
//	Value string `parser:"'$' @HexDigit"`
//}

//type Variable struct {
//	Name string `parser:"'!' @Ident"`
//}

//type Operator struct {
//	Symbol string `parser:"@Operator"`
//}

//type Expr struct {
//	Hex    *HexLiteral        `parser:"  @@"`
//	Reg    *Register          `parser:"| @@"`
//	Var    *Variable          `parser:"| @@"`
//	Square *SquareBracketExpr `parser:"| @@"`
//}

//type SquareBracketExpr struct {
//	Open  string  `parser:"'['"`
//	Parts []*Expr `parser:"@@ ( @Operator @@ )*"`
//	Close string  `parser:"']'"`
//}

//type Instruction struct {
//	Instr string `parser:"@Ident"`
//	Arg1  *Expr  `parser:"@@ ','"`
//	Arg2  *Expr  `parser:"@@"`
//}

// === Parser Constructors ===

// NewParser returns a parser for full instructions like "mov $42, r4"
//func NewParser() (*participle.Parser[Instruction], error) {
//	return participle.Build[Instruction](
///		participle.Lexer(lexerDef),
//		participle.Elide("Whitespace"),
///	)
//}

// RegisterParser returns a parser for register names (e.g. r1, sp)
//func RegisterParser() (*participle.Parser[Register], error) {
//	return participle.Build[Register](
//		participle.Lexer(lexerDef),
//	)
//}

// HexParser returns a parser for hex literals like "$42"
//func HexParser() (*participle.Parser[HexLiteral], error) {
//	return participle.Build[HexLiteral](
//		participle.Lexer(lexerDef),
//		participle.Elide("Whitespace"),
//	)
//}

// VariableParser returns a parser for variables like "!foo"
//func VariableParser() (*participle.Parser[Variable], error) {
///	return participle.Build[Variable](
//		participle.Lexer(lexerDef),
///		participle.Elide("Whitespace"),
///	)
//}

// SquareBracketParser parses expressions inside brackets like "[r1 + $10]"
//func SquareBracketParser() (*participle.Parser[SquareBracketExpr], error) {
///	return participle.Build[SquareBracketExpr](
//		participle.Lexer(lexerDef),
//		participle.Elide("Whitespace"),
//	)
//}

// MovLitToRegParser parses "mov $42, r4"
//func MovLitToRegParser() (*participle.Parser[Instruction], error) {
//	return participle.Build[Instruction](
//		participle.Lexer(lexerDef),
//		participle.Elide("Whitespace"),
//	)
//}

package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

// === Lexer ===
var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Register", Pattern: `(?i)\b(r[1-8]|sp|fp|ip|acc)\b`},
	{Name: "HexDigit", Pattern: `[0-9A-Fa-f]+`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Operator", Pattern: `[\+\-\*]`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
	{Name: "Punct", Pattern: `[\[\],!$]`},
})

// NodeType represents the type of AST node
type NodeType string

// Node types corresponding to JavaScript asType values
const (
	TypeRegister          NodeType = "REGISTER"
	TypeHexLiteral        NodeType = "HEX_LITERAL"
	TypeVariable          NodeType = "VARIABLE"
	TypeOpPlus            NodeType = "OP_PLUS"
	TypeOpMinus           NodeType = "OP_MINUS"
	TypeOpMultiply        NodeType = "OP_MULTIPLY"
	TypeSquareBracketExpr NodeType = "SQUARE_BRACKET_EXPRESSION"
	TypeInstruction       NodeType = "INSTRUCTION"
)

// Node represents a generic AST node with type and value
type Node struct {
	Type  NodeType    `json:"type"`
	Value interface{} `json:"value"`
}

// === AST Nodes ===
type Register struct {
	Value string `parser:"@Register" json:"value"`
}

// AsNode converts Register to Node
func (r *Register) AsNode() *Node {
	return &Node{
		Type:  TypeRegister,
		Value: r.Value,
	}
}

type HexLiteral struct {
	Value string `parser:"'$' @HexDigit" json:"value"`
}

// AsNode converts HexLiteral to Node
func (h *HexLiteral) AsNode() *Node {
	return &Node{
		Type:  TypeHexLiteral,
		Value: h.Value,
	}
}

type Variable struct {
	Name string `parser:"'!' @Ident" json:"name"`
}

// AsNode converts Variable to Node
func (v *Variable) AsNode() *Node {
	return &Node{
		Type:  TypeVariable,
		Value: v.Name,
	}
}

type Operator struct {
	Symbol string `parser:"@Operator" json:"symbol"`
}

// AsNode converts Operator to Node
func (o *Operator) AsNode() *Node {
	switch o.Symbol {
	case "+":
		return &Node{Type: TypeOpPlus, Value: "+"}
	case "-":
		return &Node{Type: TypeOpMinus, Value: "-"}
	case "*":
		return &Node{Type: TypeOpMultiply, Value: "*"}
	default:
		return &Node{Type: "OP_UNKNOWN", Value: o.Symbol}
	}
}

// ExprElement represents either an expression or an operator
type ExprElement struct {
	Expr     *Expr     `parser:"  @@"`
	Operator *Operator `parser:"| @@"`
}

// AsNode converts ExprElement to Node
func (e *ExprElement) AsNode() *Node {
	if e.Expr != nil {
		return e.Expr.AsNode()
	}
	return e.Operator.AsNode()
}

type Expr struct {
	Hex    *HexLiteral        `parser:"  @@"`
	Var    *Variable          `parser:"| @@"`
	Square *SquareBracketExpr `parser:"| @@"`
}

// AsNode converts Expr to Node
func (e *Expr) AsNode() *Node {
	if e.Hex != nil {
		return e.Hex.AsNode()
	}
	if e.Var != nil {
		return e.Var.AsNode()
	}
	if e.Square != nil {
		return e.Square.AsNode()
	}
	return nil
}

type SquareBracketExpr struct {
	Open     string         `parser:"'['"`
	Elements []*ExprElement `parser:"@@ ( @@ )*"`
	Close    string         `parser:"']'"`
}

// AsNode converts SquareBracketExpr to Node
func (s *SquareBracketExpr) AsNode() *Node {
	var elements []*Node
	for _, elem := range s.Elements {
		elements = append(elements, elem.AsNode())
	}
	return &Node{
		Type:  TypeSquareBracketExpr,
		Value: elements,
	}
}

type MovInstruction struct {
	Instr string    `parser:"@('mov'|'MOV')"`
	Arg1  *Expr     `parser:"@@"`
	Comma string    `parser:"','"`
	Arg2  *Register `parser:"@@"`
}

// AsNode converts MovInstruction to Node
func (m *MovInstruction) AsNode() *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": "MOV_LIT_REG",
			"args":        []*Node{m.Arg1.AsNode(), m.Arg2.AsNode()},
		},
	}
}

// DeepLog prints a structure with full depth (similar to JS deepLog)
func DeepLog(data interface{}) {
	// In a real implementation, you'd use a JSON pretty printer or similar
	// For now we'll leave this as a placeholder
}

// Upper or lowercase string helpers
func UpperOrLowerStr(s string) []string {
	return []string{strings.ToUpper(s), strings.ToLower(s)}
}

// === Parser Constructors ===
// ParseMovLitToReg parses "mov $42, r4"
func ParseMovLitToReg(input string) (*Node, error) {
	parser, err := participle.Build[MovInstruction](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, err
	}

	movInstr, err := parser.ParseString("", input)
	if err != nil {
		return nil, err
	}

	return movInstr.AsNode(), nil
}
