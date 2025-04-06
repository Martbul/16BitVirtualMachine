package main

import (
	"fmt"

	asseblerParser "github.com/martbul/assembler/parser"
)

func main() {
	// Example: Parse 'mov $42, r4'
	// Create the MOV parser
	parser, err := asseblerParser.MovLitToRegParser()
	if err != nil {
		fmt.Println("Error building MOV parser:", err)
		return
	}

	// Parse input instruction with filename (can be anything like "input")
	result, err := parser.ParseString("input", "mov $42, r4")
	if err != nil {
		fmt.Println("Error parsing instruction:", err)
		return
	}

	// Output parsed result
	fmt.Println(result)
}
