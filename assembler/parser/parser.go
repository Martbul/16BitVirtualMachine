package parser

import (
	"fmt"
	"strings"
)

// ParseProgram parses a complete program consisting of instructions and labels
func ParseProgram(input string) ([]*Node, error) {
	var nodes []*Node

	// Trim leading whitespace
	remaining := strings.TrimSpace(input)

	for len(remaining) > 0 {
		// Try to parse either an instruction or label
		node, rest, err := parseInstructionOrLabel(remaining)
		if err != nil {
			return nil, fmt.Errorf("parse error at '%s': %v", remaining, err)
		}

		nodes = append(nodes, node)

		// Update remaining input
		remaining = strings.TrimSpace(rest)
	}

	return nodes, nil
}

// parseInstructionOrLabel attempts to parse either an instruction or a label
// This is the Go equivalent of A.choice([instructionsParser, label])
func parseInstructionOrLabel(input string) (*Node, string, error) {
	// Try label first since it's simpler to detect
	if isLabelLine(input) {
		return parseLabel(input)
	}

	// If not a label, try parsing as an instruction
	instruction, restAfterInstruction, err := parseInstruction(input)
	if err != nil {
		return nil, input, fmt.Errorf("failed to parse as instruction or label: %v", err)
	}

	return instruction, restAfterInstruction, nil
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
func parseInstruction(input string) (*Node, string, error) {
	// Find where the instruction ends (newline or semicolon)
	endIndex := strings.IndexAny(input, "\n;")
	var instructionText string
	var rest string

	if endIndex == -1 {
		// No newline or semicolon, assume instruction takes the whole input
		instructionText = strings.TrimSpace(input)
		rest = ""
	} else {
		instructionText = strings.TrimSpace(input[:endIndex])
		rest = input[endIndex+1:]
	}

	// Use the existing ParseInstruction function from instructions.go
	node, err := ParseInstruction(instructionText)
	if err != nil {
		return nil, input, err
	}

	return node, rest, nil
}
