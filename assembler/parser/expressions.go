package parser

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
