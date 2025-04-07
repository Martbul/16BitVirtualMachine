package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
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

// Upper or lowercase string helpers
func UpperOrLowerStr(s string) []string {
	return []string{strings.ToUpper(s), strings.ToLower(s)}
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

// DeepLog prints a structure with full depth (similar to JS deepLog)
func DeepLog(data interface{}) {
	// Marshal the data to JSON with indentation
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal data: %v\n", err)
		return
	}

	// Print the JSON representation
	fmt.Println(string(jsonBytes))
}

// PrettyPrintNode prints a Node and its children with proper indentation
func PrettyPrintNode(node *Node) {
	jsonBytes, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal node: %v\n", err)
		return
	}

	fmt.Println(string(jsonBytes))
}
