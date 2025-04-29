package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

//	var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
//		{Name: "Register", Pattern: `(?i)\b(r[1-8]|sp|fp|ip|acc)\b`},
//		{Name: "HexDigit", Pattern: `[0-9A-Fa-f]+`},
//		{Name: "Instruction", Pattern: `(?i)\b(mov|add|sub|inc|dec|mul|lsf|rsf|and|or|xor|not|jmp|jne|jeq|jlt|jgt|jle|jge|psh|pop|cal|ret|hlt)\b`},
//		{Name: "Constant", Pattern: `(?i)constant`},
//		{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
//		{Name: "Operator", Pattern: `[\+\-\*]`},
//		{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
//		{Name: "Punct", Pattern: `[\[\],!$&=+]`},
//		{Name: "DataType", Pattern: `data(8|16)`},
//		{Name: "Export", Pattern: `\+`},
//
// {Name: "Keyword", Pattern: `(?i)(constant)`},
// })
// WARN: SOMEHOW THE "constant code_const = $C0DE" was parseed, if there is a future error it is likely to be in the lexerDef or HexLiteral or constant struct or the way the lexer handles punctuation and $ signs or other...
var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
	// Keywords and symbols FIRST
	{Name: "Constant", Pattern: `(?i)constant`},
	{Name: "Export", Pattern: `\+`},
	{Name: "Equals", Pattern: `=`},
	{Name: "HexDigit", Pattern: `\$[0-9A-Fa-f]+`},

	// Then instruction, register, etc.
	{Name: "Instruction", Pattern: `(?i)\b(mov|add|sub|inc|dec|mul|lsf|rsf|and|or|xor|not|jmp|jne|jeq|jlt|jgt|jle|jge|psh|pop|cal|ret|hlt)\b`},
	{Name: "Register", Pattern: `(?i)\b(r[1-8]|sp|fp|ip|acc)\b`},
	{Name: "DataType", Pattern: `data(8|16)`},

	// THEN Ident last
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},

	// Punctuation
	{Name: "Punct", Pattern: `[\[\],!$&]`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
})

//var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
//	// Keywords and symbols FIRST
//	{Name: "Constant", Pattern: `(?i)constant`},
//	{Name: "Dollar", Pattern: `\$`},
//	{Name: "HexDigit", Pattern: `[0-9A-Fa-f]+`},
//	{Name: "Export", Pattern: `\+`},
//	{Name: "Equals", Pattern: `=`},
//
// Then instruction, register, etc.
//	{Name: "Instruction", Pattern: `(?i)\b(mov|add|sub|inc|dec|mul|lsf|rsf|and|or|xor|not|jmp|jne|jeq|jlt|jgt|jle|jge|psh|pop|cal|ret|hlt)\b`},
//	{Name: "Register", Pattern: `(?i)\b(r[1-8]|sp|fp|ip|acc|im)\b`},
///	{Name: "DataType", Pattern: `data(8|16)`},

//	// THEN Ident last
///	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},

//	// Punctuation
//	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
//	{Name: "Punct", Pattern: `[\[\],!&]`}, // Removed $ here!
//})

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

func DeepLog(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal data: %v\n", err)
		return
	}

	fmt.Println(string(jsonBytes))
}

func PrettyPrintNode(node *Node) {
	jsonBytes, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal node: %v\n", err)
		return
	}

	fmt.Println(string(jsonBytes))
}
