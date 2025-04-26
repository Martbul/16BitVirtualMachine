package parser

import (
	"github.com/alecthomas/participle/v2"
)

// Constant represents a constant definition in assembly
type Constant struct {
	IsExport bool        `parser:"@('+')?"`
	Name     string      `parser:"'constant' @Ident"`
	Value    *HexLiteral `parser:"'=' @@"`
}

// AsNode converts Constant to Node
func (c *Constant) AsNode() *Node {
	return &Node{
		Type: "CONSTANT",
		Value: map[string]interface{}{
			"isExport": c.IsExport,
			"name":     c.Name,
			"value":    c.Value.Value,
		},
	}
}

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
		return nil, err
	}

	return constant.AsNode(), nil
}
