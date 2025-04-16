package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
)

// Label represents a label in the assembly code
type Label struct {
	Name string `@Ident ":"`
}

// AsNode converts a label to a Node
func (l *Label) AsNode() *Node {
	return &Node{
		Type: "LABEL",
		Value: map[string]interface{}{
			"name": l.Name,
		},
	}
}

// ProgramLine represents either an instruction or a label in the program
type ProgramLine struct {
	Instruction *Node  // This will be populated by our instruction parser
	Label       *Label `@@`
}

// ParseLabel parses a label in the assembly code
func ParseLabel(input string) (*Node, error) {
	parser, err := participle.Build[Label](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, fmt.Errorf("error building label parser: %v", err)
	}

	label, err := parser.ParseString("", input)
	if err != nil {
		return nil, err
	}

	return label.AsNode(), nil
}

// ParseProgram parses an entire program consisting of both instructions and labels
func ParseProgram(input string) ([]*Node, error) {
	lines := strings.Split(input, "\n")
	var nodes []*Node

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue // Skip empty lines
		}

		// First try to parse as an instruction
		instrNode, instrErr := ParseInstruction(line)
		if instrErr == nil {
			nodes = append(nodes, instrNode)
			continue
		}

		// If that fails, try to parse as a label
		labelNode, labelErr := ParseLabel(line)
		if labelErr == nil {
			nodes = append(nodes, labelNode)
			continue
		}

		// If both fail, report the error
		return nil, fmt.Errorf("failed to parse line '%s': instruction error: %v, label error: %v",
			line, instrErr, labelErr)
	}

	return nodes, nil
}

// ParseFile parses an entire assembly file
func ParseFile(filename string) ([]*Node, error) {
	content, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParseProgram(content)
}

// Helper function to read a file
func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
