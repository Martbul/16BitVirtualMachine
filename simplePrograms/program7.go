package simpleprograms

import (
	"fmt"
	"os"

	"github.com/martbul/assembler/parser"
)

func SimpleProgram7() {
	//input := "pop acc" WARN: DOESNT WORK
	input := "pop acc"
	node, err := parser.ParseInstruction(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
		return
	}

	fmt.Println("AST for:", input)
	//parser.DeepLog(node)
	// or
	parser.PrettyPrintNode(node)

}
