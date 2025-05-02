package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

//type Constant struct {
//	IsExport bool        `parser:"@('+' )?"`
//	Keyword  string      `parser:"@Constant"`
//	Name     string      `parser:"@Ident"`
//	Value    *HexLiteral `parser:"'=' @@"`
//}

//type Constant struct {
//	IsExport bool        `parser:"@('+' )?"`
//	Name     string      `parser:"@Constant @Ident"`
//	Value    *HexLiteral `parser:"'=' @@"`
//}

//type Constant struct {
//	IsExport bool        `parser:"@('+' )?"`  // Optional export
//	Keyword  string      `parser:"@Constant"` // Matches the 'constant' keyword explicitly
//	Name     string      `parser:"@Ident"`    // Matches the identifier for the constant name
//	Value    *HexLiteral `parser:"'=' @@"`    // Matches the equals sign and then the HexLiteral
//}

//type Constant struct {
//	IsExport bool        `parser:"@('+' )?"`  // Optional '+' (export)
///	Keyword  string      `parser:"@Constant"` // Matches the 'constant' keyword
///	Name     string      `parser:"@Ident"`    // Matches the identifier (e.g., code_const)
///	Equals   string      `parser:"'='"`       // Matches the '=' sign explicitly
//	Value    *HexLiteral `parser:"@@"`        // Matches the HexLiteral (e.g., $C0DE)
//}

//type Constant struct {
//	IsExport bool        `parser:"@('+' )?"`
//	Keyword  string      `parser:"@Constant"` // This matches the token `Constant`
//	Name     string      `parser:"@Ident"`
//	Value    *HexLiteral `parser:"'=' @@"`
//}
// AsNode converts Constant to Node
//func (c *Constant) AsNode() *Node {
//	return &Node{
//		Type: "CONSTANT",
//		Value: map[string]interface{}{
//			"isExport": c.IsExport,
//			"name":     c.Name,
//			"value":    c.Value.Value,
//		},
//	}
//}

// ParseConstant attempts to parse a constant declaration
func ParseConstant(input string) (*Node, error) {
	parser, err := participle.Build[Constant](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {

		return nil, err
	}

	constant, err := parser.ParseString("", input)
	if err != nil {

		fmt.Println(err)
		return nil, err
	}

	return constant.AsNode(), nil
}
