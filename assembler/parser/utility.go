// Add this to your assembler/parser/utility.go or create a new file

package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// This function extracts the variable name from a memory reference with a constant
// Example: [!code_constant] -> code_constant
func extractConstantName(memRef string) (string, bool) {
	// Check if the memory reference contains a constant reference with !
	if strings.Contains(memRef, "!") {
		// Extract the constant name between ! and ]
		re := regexp.MustCompile(`\[!([a-zA-Z0-9_]+)\]`)
		matches := re.FindStringSubmatch(memRef)
		if len(matches) >= 2 {
			return matches[1], true
		}
	}
	return "", false
}

// Add a function to parse memory reference with constants
func parseMemoryReferenceWithConstant(input string) (*Node, error) {
	// Extract the constant name
	constName, ok := extractConstantName(input)
	if !ok {
		return nil, fmt.Errorf("invalid memory reference format: %s", input)
	}

	// Create a memory reference node that includes the constant
	return &Node{
		Type: "MEMORY_REFERENCE",
		Value: map[string]interface{}{
			"type":  "VARIABLE",
			"value": constName,
		},
	}, nil
}
