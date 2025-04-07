package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

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
