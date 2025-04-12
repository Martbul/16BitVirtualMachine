package simpleprograms

import (
	"fmt"
	"os"

	"github.com/martbul/assembler/parser"
)

func SimpleProgram6() {
	//	input := "mov [$42 + !loc - ($05 * ($31 + !var) - $07)], r4"

	//input := "mov $42, &C0DE"
	//	input := "mov acc, r1"
	//	input := "mov acc, &[!loc + $4200]"
	//input := "mov &4200, r1"
	//input := "mov &r3, acc"
	input := "mov $42, &r1, r4"
	node, err := parser.ParseMovInstruction(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
		return
	}

	fmt.Println("AST for:", input)
	//parser.DeepLog(node)
	// or
	parser.PrettyPrintNode(node)

}
