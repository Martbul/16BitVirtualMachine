package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
)

// DataNode represents a data declaration (8-bit or 16-bit)
type DataNode struct {
	Type     NodeType `json:"type"`
	Size     int      `json:"size"`     // 8 or 16
	IsExport bool     `json:"isExport"` // Whether the data is exported (prefixed with '+')
	Name     string   `json:"name"`     // Name of the data
	Values   []string `json:"values"`   // Array of hex values
}

// DataDeclaration represents the parsed structure of a data declaration
type DataDeclaration struct {
	IsExport bool     `parser:"@'+'?"`
	DataType string   `parser:"@('data8'|'data16')"` // Matches 'data8' or 'data16'
	Name     string   `parser:"@Ident"`
	Equals   string   `parser:"'='"`
	Open     string   `parser:"@Brace"`
	Values   []string `parser:"@HexDigit (',' @HexDigit)*"`
	Close    string   `parser:"@Brace"`
}

// AsNode converts the DataDeclaration to a DataNode
func (d *DataDeclaration) AsNode() *DataNode {
	size := 8
	if d.DataType == "data16" {
		size = 16
	}

	return &DataNode{
		Type:     "DATA",
		Size:     size,
		IsExport: d.IsExport,
		Name:     d.Name,
		Values:   d.Values,
	}
}

// ParseData8 parses an 8-bit data declaration
func ParseData8(input string) (*Node, error) {
	return parseData(input, 8)
}

// ParseData16 parses a 16-bit data declaration
func ParseData16(input string) (*Node, error) {
	return parseData(input, 16)
}

func parseData(input string, size int) (*Node, error) {
	// Build the parser
	parser, err := participle.Build[DataDeclaration](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, err
	}

	decl, err := parser.ParseString("", input)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(decl)
	//	fmt.Println(size)
	expectedType := fmt.Sprintf("data%d", size)

	//	fmt.Println(expectedType)
	//	fmt.Println(decl.DataType)
	if !strings.EqualFold(decl.DataType, expectedType) {
		return nil, fmt.Errorf("expected data type %s, got %s", expectedType, decl.DataType)
	}

	dataNode := decl.AsNode()

	return &Node{
		Type: "DATA_DECLARATION",
		Value: map[string]interface{}{
			"size":     dataNode.Size,
			"isExport": dataNode.IsExport,
			"name":     dataNode.Name,
			"values":   dataNode.Values,
		},
	}, nil
}

// Now update the lexer to support the new tokens needed for data declarations
func init() {
	// Make sure the lexer can recognize data8 and data16 keywords
	// This should be integrated with your existing lexer definitions
}
