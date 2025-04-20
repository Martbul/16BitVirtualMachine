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
func parseInstructionOrLabel(input string) (*Node, string, error) {
	// Try label first (based on your existing label parser)
	labelNode, restAfterLabel, labelErr := parseLabel(input)
	if labelErr == nil {
		return labelNode, restAfterLabel, nil
	}

	// If label parsing failed, try instruction parsing
	for _, parser := range instructionParsers {
		node, err := parser.Fn(input)
		if err == nil {
			// Need to figure out how much input was consumed
			// This depends on your parser implementation
			// For now, let's assume your parser consumes up to newline or semicolon
			restIndex := strings.IndexAny(input, "\n;")
			if restIndex == -1 {
				return node, "", nil // Consumed all input
			}
			return node, input[restIndex+1:], nil
		}
	}

	return nil, input, fmt.Errorf("failed to parse instruction or label")
}

// parseLabel attempts to parse a label
func parseLabel(input string) (*Node, string, error) {
	// Implementation depends on your label format
	// For example, if labels are in the format "label_name:"
	labelEnd := strings.Index(input, ":")
	if labelEnd == -1 {
		return nil, input, fmt.Errorf("not a label")
	}

	labelName := strings.TrimSpace(input[:labelEnd])
	if len(labelName) == 0 {
		return nil, input, fmt.Errorf("empty label name")
	}

	// Create a label node
	labelNode := &Node{
		Type: "LABEL", // Assuming you have a label type
		Value: map[string]interface{}{
			"name": labelName,
		},
	}

	return labelNode, input[labelEnd+1:], nil
}

// instructionParsers contains all your instruction parsers
var instructionParsers = []Parser{
	{Name: "MOV_LIT_REG", Fn: LitToReg("mov", "MOV_LIT_REG")},
	{Name: "MOV_REG_REG", Fn: RegToReg("mov", "MOV_REG_REG")},
	{Name: "MOV_REG_MEM", Fn: RegToMem("mov", "MOV_REG_MEM")},
	{Name: "MOV_MEM_REG", Fn: MemToReg("mov", "MOV_MEM_REG")},
	{Name: "MOV_LIT_MEM", Fn: LitToMem("mov", "MOV_LIT_MEM")},
	// Add all other instruction parsers here
}
