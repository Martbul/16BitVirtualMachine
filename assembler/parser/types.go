package parser

// === TypeSystem === //
// NodeType represents the type of AST node
type NodeType string

// Node types corresponding to JavaScript asType values
const (
	TypeRegister          NodeType = "REGISTER"
	TypeLabel             NodeType = "LABEL"
	TypeRegisterPointer   NodeType = "REGISTER_POINTER"
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

// Address represents a direct memory address
type Address struct {
	Value string `parser:"'&' @HexDigit" json:"value"`
}

// AsNode converts Address to Node
func (a *Address) AsNode() *Node {
	return &Node{
		Type:  "ADDRESS",
		Value: a.Value,
	}
}

// RegisterPointer represents a register pointer (e.g., &r1)
type RegisterPointer struct {
	Register *Register `parser:"'&' @@" json:"register"`
}

// AsNode converts RegisterPointer to Node
func (rp *RegisterPointer) AsNode() *Node {
	return &Node{
		Type:  TypeRegisterPointer,
		Value: rp.Register.Value,
	}
}

// MemoryReference represents either a direct address or a square bracket expression with an ampersand
type MemoryReference struct {
	DirectAddress *Address           `parser:"  @@"`
	Expression    *SquareBracketExpr `parser:"| '&' @@"`
}

// AsNode converts MemoryReference to Node
func (m *MemoryReference) AsNode() *Node {
	if m.DirectAddress != nil {
		return m.DirectAddress.AsNode()
	}

	// For expressions, create a memory reference node that wraps the square bracket expression
	exprNode := m.Expression.AsNode()
	return &Node{
		Type:  "MEMORY_REFERENCE",
		Value: exprNode,
	}
}

//WARN: Personal fuckup

// MemoryReference represents either a direct address or a square bracket expression with an ampersand
type LiteralReference struct {
	DirectLiteral *HexLiteral        `parser:"  @@"`
	Expression    *SquareBracketExpr `parser:"| '&' @@"`
}

// AsNode converts MemoryReference to Node
func (m *LiteralReference) AsNode() *Node {
	if m.DirectLiteral != nil {
		return m.DirectLiteral.AsNode()
	}

	// For expressions, create a memory reference node that wraps the square bracket expression
	exprNode := m.Expression.AsNode()
	return &Node{
		Type:  "LITERAL_REFERENCE",
		Value: exprNode,
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

type LitRegInstruction struct {
	//Instr string `parser:"@(/[a-zA-Z]+/)"` WARN: Doesnt work
	//	Instr string `parser:"@Instruction"` //WARN: Doesnt work
	Instr string `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
	//Instr string `parser:"@Ident"` // Generic identifier instead of hardcoded mnemonic WARN: Doesnt work
	//Instr string    `parser:"@('add'|'ADD')"` WARN: When i use this it works
	Arg1  *Expr     `parser:"@@"`
	Comma string    `parser:"','"`
	Arg2  *Register `parser:"@@"`
}

// Fixed AsNode method with proper signature and variable references
func (instr *LitRegInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{instr.Arg1.AsNode(), instr.Arg2.AsNode()},
		},
	}
}

type RegRegInstruction struct {
	Instr string `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
	//Instr string    `parser:"@Ident"` // Generic identifier instead of hardcoded mnemonic
	Reg1  *Register `parser:"@@"`
	Comma string    `parser:"','"`
	Reg2  *Register `parser:"@@"`
}

func (instr *RegRegInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{instr.Reg1.AsNode(), instr.Reg2.AsNode()},
		},
	}
}

type RegMemInstruction struct {
	Instr  string           `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
	Reg    *Register        `parser:"@@"`
	Comma  string           `parser:"','"`
	Memory *MemoryReference `parser:"@@"`
}

func (instr *RegMemInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{instr.Reg.AsNode(), instr.Memory.AsNode()},
		},
	}
}

type MemRegInstruction struct {
	Instr  string           `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
	Memory *MemoryReference `parser:"@@"`
	Comma  string           `parser:"','"`
	Reg    *Register        `parser:"@@"`
}

func (instr *MemRegInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{instr.Memory.AsNode(), instr.Reg.AsNode()},
		},
	}
}

type LitMemInstruction struct {
	Instr  string            `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
	Lit    *LiteralReference `parser:"@@"`
	Comma  string            `parser:"','"`
	Memory *MemoryReference  `parser:"@@"`
}

func (instr *LitMemInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{instr.Lit.AsNode(), instr.Memory.AsNode()},
		},
	}
}

type RegPtrToRegInstruction struct {
	Instr  string           `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
	RegPtr *RegisterPointer `parser:"@@"`
	Comma  string           `parser:"','"`
	Reg    *Register        `parser:"@@"`
}

func (instr *RegPtrToRegInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{instr.RegPtr.AsNode(), instr.Reg.AsNode()},
		},
	}
}

type LitOffToRegInstruction struct {
	Instr  string            `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
	Lit    *LiteralReference `parser:"@@"`
	Comma1 string            `parser:"','"`
	RegPtr *RegisterPointer  `parser:"@@"`
	Comma2 string            `parser:"','"`
	Reg    *Register         `parser:"@@"`
}

func (instr *LitOffToRegInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{instr.Lit.AsNode(), instr.RegPtr.AsNode(), instr.Reg.AsNode()},
		},
	}
}

type NoArgsInstruction struct {
	Instr string `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
}

func (instr *NoArgsInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{},
		},
	}
}

// NoArgInstruction represents instructions without arguments (like HLT, RET)
type NoArgInstruction struct {
	Instr string `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
}

type SingleRegInstruction struct {
	Instr string    `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
	Reg   *Register `parser:"@@"`
}

func (instr *SingleRegInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{instr.Reg.AsNode()},
		},
	}

}

type SingleLitInstruction struct {
	Instr string            `parser:"@('MOV'|'mov'|'ADD'|'add'|'SUB'|'sub'|'MUL'|'mul'|'LSF'|'lsf'|'RSF'|'rsf'|'AND'|'and'|'OR'|'or'|'XOR'|'xor'|'JMP'|'jmp'|'JNE'|'jne'|'JEQ'|'jeq'|'JLT'|'jlt'|'JGT'|'jgt'|'JLE'|'jle'|'JGE'|'jge'|'PSH'|'psh'|'POP'|'pop'|'CAL'|'cal'|'RET'|'ret'|'HLT'|'hlt'|'INC'|'inc'|'DEC'|'dec'|'NOT'|'not')"`
	Lit   *LiteralReference `parser:"@@"`
}

func (instr *SingleLitInstruction) AsNode(instructionType string) *Node {
	return &Node{
		Type: TypeInstruction,
		Value: map[string]interface{}{
			"instruction": instructionType,
			"args":        []*Node{instr.Lit.AsNode()},
		},
	}

}

// BinaryOperation represents a binary operation with left operand, operator, and right operand
type BinaryOperation struct {
	A  *Node
	Op *Node
	B  *Node
}
