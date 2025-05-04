package parser

import (
	"fmt"
	"strings"
)

var constantsMap = make(map[string]string)

// ParseProgram parses a complete program consisting of instructions and labels
//func ParseProgram(input string) ([]*Node, error) {
//	var nodes []*Node
//	constantsMap = make(map[string]string) // Reset constants map
//	// Trim leading whitespace
//	remaining := strings.TrimSpace(input)

//	for len(remaining) > 0 {
//		if isConstantLine(remaining) {
//			node, rest, err := parseConstant(remaining)
//			if err != nil {
//				return nil, fmt.Errorf("error parsing constant: %v", err)
//			}
//			// Store the constant in our map
//			constantValue := node.Value.(map[string]interface{})
//			constantsMap[constantValue["name"].(string)] = constantValue["value"].(map[string]interface{})["value"].(string)
//			nodes = append(nodes, node)
//			remaining = strings.TrimSpace(rest)
//			continue
//		}
// Try to parse either an instruction or label
//		node, rest, err := parseInstructionOrLabel(remaining)
//		if err != nil {

//			return nil, fmt.Errorf("parse error at '%s': %v", remaining, err)
//		}

//		nodes = append(nodes, node)

// Update remaining input
///		remaining = strings.TrimSpace(rest)
//	}

//		return nodes, nil
//	}
//
// ParseProgram parses a complete program consisting of instructions and labels
func ParseProgram(input string) ([]*Node, error) {
	var nodes []*Node
	constantsMap = make(map[string]string) // Reset constants map
	remaining := strings.TrimSpace(input)

	for len(remaining) > 0 {
		if isConstantLine(remaining) {
			node, rest, err := parseConstant(remaining)
			if err != nil {
				return nil, fmt.Errorf("error parsing constant: %v", err)
			}

			// Safely extract constant value
			constantValue, ok := node.Value.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid constant format")
			}

			name, ok := constantValue["name"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid constant name")
			}

			value, ok := constantValue["value"]
			if !ok {
				return nil, fmt.Errorf("missing constant value")
			}

			// Handle both string and map value types
			var valueStr string
			switch v := value.(type) {
			case string:
				valueStr = v
			case map[string]interface{}:
				if val, exists := v["value"]; exists {
					if strVal, ok := val.(string); ok {
						valueStr = strVal
					} else {
						return nil, fmt.Errorf("invalid constant value type")
					}
				} else {
					return nil, fmt.Errorf("missing nested constant value")
				}
			default:
				return nil, fmt.Errorf("unexpected constant value type")
			}

			constantsMap[name] = valueStr
			nodes = append(nodes, node)
			remaining = strings.TrimSpace(rest)
			continue
		}

		node, rest, err := parseInstructionOrLabel(remaining)
		if err != nil {
			return nil, fmt.Errorf("parse error at '%s': %v", remaining, err)
		}

		nodes = append(nodes, node)
		remaining = strings.TrimSpace(rest)
	}

	return nodes, nil
}
func resolveConstantReferences(input string) (string, error) {
	var result strings.Builder
	i := 0

	for i < len(input) {
		if i+1 < len(input) && input[i] == '[' && input[i+1] == '!' {
			// Found constant reference
			end := strings.IndexByte(input[i:], ']')
			if end == -1 {
				return "", fmt.Errorf("unclosed constant reference")
			}
			end += i

			constantName := input[i+2 : end]
			constantValue, exists := constantsMap[constantName]
			if !exists {
				return "", fmt.Errorf("undefined constant: %s", constantName)
			}

			result.WriteString(constantValue)
			i = end + 1
		} else {
			result.WriteByte(input[i])
			i++
		}
	}

	return result.String(), nil
}

// Add to parseInstructionOrLabel in parser.go
func parseInstructionOrLabel(input string) (*Node, string, error) {

	if isDataDeclaration(input) {
		return parseDataDeclaration(input)
	}

	if isLabelLine(input) {
		return parseLabel(input)
	}
	if isConstantLine(input) {
		return parseConstant(input)
	}
	// If not a label or data declaration, try parsing as an instruction
	instruction, restAfterInstruction, err := parseInstruction(input)
	if err != nil {
		fmt.Println("hhh")

		return nil, input, fmt.Errorf("failed to parse as instruction, label, or data declaration: %v", err)
	}

	return instruction, restAfterInstruction, nil
}

func isConstantLine(input string) bool {
	trimmed := strings.TrimSpace(input)
	return strings.HasPrefix(trimmed, "constant") || strings.HasPrefix(trimmed, "+constant")
}

func parseConstant(input string) (*Node, string, error) {
	endIndex := strings.IndexAny(input, "\n;") //gets only the single line with the constant
	var constantText string
	var rest string

	if endIndex == -1 {
		// No newline or semicolon, assume constant takes the whole input
		constantText = strings.TrimSpace(input)
		rest = ""
	} else {
		constantText = strings.TrimSpace(input[:endIndex])
		rest = input[endIndex+1:]
	}

	node, err := ParseConstant(constantText)
	if err != nil {
		return nil, input, err

	}

	return node, rest, nil
}

// isDataDeclaration checks if the line starts with a data declaration
func isDataDeclaration(input string) bool {
	trimmed := strings.TrimSpace(input)
	return strings.HasPrefix(trimmed, "data8") ||
		strings.HasPrefix(trimmed, "data16") ||
		strings.HasPrefix(trimmed, "+data8") ||
		strings.HasPrefix(trimmed, "+data16")
}

// parseDataDeclaration attempts to parse a data declaration
func parseDataDeclaration(input string) (*Node, string, error) {
	// Find where the declaration ends (closing brace followed by whitespace or EOF)
	closeBraceIndex := strings.Index(input, "}")
	if closeBraceIndex == -1 {
		return nil, input, fmt.Errorf("missing closing brace in data declaration")
	}

	// Include the closing brace in the declaration
	declarationText := strings.TrimSpace(input[:closeBraceIndex+1])
	rest := input[closeBraceIndex+1:]

	// Determine if it's data8 or data16
	isData8 := strings.Contains(declarationText, "data8")

	var node *Node
	var err error

	if isData8 {
		node, err = ParseData8(declarationText)
	} else {
		node, err = ParseData16(declarationText)
	}

	if err != nil {
		return nil, input, err
	}

	return node, rest, nil
}

// isLabelLine checks if the line starts with a label
func isLabelLine(input string) bool {
	// In most assembly languages, labels end with a colon
	// And are not indented (for simplicity we'll ignore whitespace)
	trimmed := strings.TrimSpace(input)

	// Label can't start with a digit
	if len(trimmed) == 0 || (trimmed[0] >= '0' && trimmed[0] <= '9') {
		return false
	}

	// Look for a colon that's not inside a string or comment
	colonIndex := strings.Index(trimmed, ":")
	if colonIndex == -1 {
		return false
	}

	// Check if what's before the colon is a valid identifier
	labelName := trimmed[:colonIndex]
	return isValidIdentifier(labelName)
}

// isValidIdentifier checks if a string is a valid label identifier
func isValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Must start with a letter or underscore
	if !((s[0] >= 'a' && s[0] <= 'z') ||
		(s[0] >= 'A' && s[0] <= 'Z') ||
		s[0] == '_') {
		return false
	}

	// Rest can be alphanumeric or underscore
	for i := 1; i < len(s); i++ {
		if !((s[i] >= 'a' && s[i] <= 'z') ||
			(s[i] >= 'A' && s[i] <= 'Z') ||
			(s[i] >= '0' && s[i] <= '9') ||
			s[i] == '_') {
			return false
		}
	}

	return true
}

// parseLabel attempts to parse a label
func parseLabel(input string) (*Node, string, error) {
	// Find the end of the label (colon)
	colonIndex := strings.Index(input, ":")
	if colonIndex == -1 {
		return nil, input, fmt.Errorf("not a label")
	}

	labelName := strings.TrimSpace(input[:colonIndex])
	if len(labelName) == 0 {
		return nil, input, fmt.Errorf("empty label name")
	}

	// Create a label node
	labelNode := &Node{
		Type: TypeLabel, // Use constant instead of string literal
		Value: map[string]interface{}{
			"label": labelName,
		},
	}

	// Return the rest of the input after the colon
	rest := input[colonIndex+1:]
	return labelNode, rest, nil
}

// parseInstruction attempts to parse an instruction and returns the consumed input
//func parseInstruction(input string) (*Node, string, error) {
// Find where the instruction ends (newline or semicolon)
//	endIndex := strings.IndexAny(input, "\n;")
//	var instructionText string
//	var rest string
//
//	if endIndex == -1 {
//		// No newline or semicolon, assume instruction takes the whole input
//		instructionText = strings.TrimSpace(input)
//		rest = ""
//	} else {
//		instructionText = strings.TrimSpace(input[:endIndex])
//		rest = input[endIndex+1:]
//	}

//	fmt.Println(instructionText)
// Use the existing ParseInstruction function from instructions.go
//	node, err := ParseInstruction(instructionText)
//	if err != nil {
//		return nil, input, err
//	}

//	return node, rest, nil
////}

func parseInstruction(input string) (*Node, string, error) {
	endIndex := strings.IndexAny(input, "\n;")
	var instructionText string
	var rest string

	if endIndex == -1 {
		instructionText = strings.TrimSpace(input)
		rest = ""
	} else {
		instructionText = strings.TrimSpace(input[:endIndex])
		rest = input[endIndex+1:]
	}

	// Resolve constant references before parsing
	resolvedInstruction, err := resolveConstantReferences(instructionText)
	if err != nil {
		return nil, input, fmt.Errorf("error resolving constants: %v", err)
	}

	node, err := ParseInstruction(resolvedInstruction)
	if err != nil {
		return nil, input, err
	}

	return node, rest, nil
}
