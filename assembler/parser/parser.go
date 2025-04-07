// mov $42, r1 -> (instruction literal hex value, register)

// move [$42 + (!loc - $1F)], r1 -> (instruction [literal hex value + (!variable - )])
package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

// === Lexer === //
var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Register", Pattern: `(?i)\b(r[1-8]|sp|fp|ip|acc)\b`}, // Matches register names(r1-r8, sp,fp,ip,acc)
	{Name: "HexDigit", Pattern: `[0-9A-Fa-f]+`},                  // Matches hexadecimal digits(0-9, A-F,a-f)
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},           // Matches identifiers(variable names)
	{Name: "Operator", Pattern: `[\+\-\*]`},                      // Matches arithmetic operations(+,-,*)
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},                  // Matches spaces, tabs, newlines
	{Name: "Punct", Pattern: `[\[\],!$()\-]`},                    // Matches punctuation chars (brackets, commas, ! and $ symbols)
})

// === TypeSystem === //
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
	TypeBracketedExpr     NodeType = "BRACKETED_EXPRESSION"
	TypeBinaryOperation   NodeType = "BINARY_OPERATION"
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
	Hex     *HexLiteral        `parser:"  @@"`
	Var     *Variable          `parser:"| @@"`
	Square  *SquareBracketExpr `parser:"| @@"`
	Bracket *BracketedExpr     `parser:"| @@"`
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
	if e.Bracket != nil {
		return e.Bracket.AsNode()
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
	node := &Node{
		Type:  TypeSquareBracketExpr,
		Value: elements,
	}
	return DisambiguateOrderOfOperations(node)
}

type BracketedExpr struct { //INFO: Representing the different states Ope, Closed and elems slice is for operators, brackets and elements
	Open     string         `parser:"'('"`
	Elements []*ExprElement `parser:"@@ ( @@ )*"`
	Close    string         `parser:"')'"`
}

// AsNode converts BracketedExpr to Node
func (b *BracketedExpr) AsNode() *Node {
	var elements []*Node
	for _, elem := range b.Elements {
		elements = append(elements, elem.AsNode())
	}
	node := &Node{
		Type:  TypeBracketedExpr,
		Value: elements,
	}
	return DisambiguateOrderOfOperations(node)
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

// BinaryOperation represents a binary operation with left operand, operator, and right operand
type BinaryOperation struct {
	A  *Node
	Op *Node
	B  *Node
}

// CreateBinaryOperation creates a binary operation node
func CreateBinaryOperation(a, op, b *Node) *Node {
	return &Node{
		Type: TypeBinaryOperation,
		Value: map[string]interface{}{
			"a":  a,
			"op": op,
			"b":  b,
		},
	}
}

// DisambiguateOrderOfOperations transforms the AST to respect operator precedence
func DisambiguateOrderOfOperations(expr *Node) *Node {
	// Check if the expression is not a bracketed expression
	if expr.Type != TypeSquareBracketExpr && expr.Type != TypeBracketedExpr {
		return expr
	}

	// Get the value as a slice of nodes
	elements, ok := expr.Value.([]*Node)
	if !ok {
		return expr
	}

	// If there's only one element, return it directly
	if len(elements) == 1 {
		return elements[0]
	}

	// Define operator priorities
	priorities := map[NodeType]int{
		TypeOpMultiply: 2,
		TypeOpPlus:     1,
		TypeOpMinus:    0,
	}

	// Find the highest priority operator
	candidateIndex := -1
	highestPriority := -1

	for i := 1; i < len(elements); i += 2 {
		priority, exists := priorities[elements[i].Type]
		if exists && priority > highestPriority {
			highestPriority = priority
			candidateIndex = i
		}
	}

	if candidateIndex == -1 {
		return expr // No operators found
	}

	// Create binary operation with the highest priority operator
	leftOperand := DisambiguateOrderOfOperations(elements[candidateIndex-1])
	rightOperand := DisambiguateOrderOfOperations(elements[candidateIndex+1])
	operator := elements[candidateIndex]

	binaryOp := CreateBinaryOperation(leftOperand, operator, rightOperand)

	// Create a new expression with the binary operation replacing the original elements
	var newElements []*Node
	newElements = append(newElements, elements[:candidateIndex-1]...)
	newElements = append(newElements, binaryOp)
	newElements = append(newElements, elements[candidateIndex+2:]...)

	if len(newElements) == 1 {
		return newElements[0]
	}

	newExpr := &Node{
		Type:  TypeBracketedExpr,
		Value: newElements,
	}

	// Recursively process the resulting expression
	return DisambiguateOrderOfOperations(newExpr)
}

// Upper or lowercase string helpers
func UpperOrLowerStr(s string) []string {
	return []string{strings.ToUpper(s), strings.ToLower(s)}
}

// === Parser Constructors ===
// ParseMovLitToReg parses "mov $42, r4" and more complex expressions
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
