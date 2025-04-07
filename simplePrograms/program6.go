package simpleprograms

import (
	"fmt"
	"os"

	"github.com/martbul/assembler/parser"
)

func SimpleProgram6() {
	input := "mov [$42 + !loc - ($05 * ($31 + !var) - $07)], r4"
	node, err := parser.ParseMovLitToReg(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
		return
	}

	fmt.Println("AST for:", input)
	//parser.DeepLog(node)
	// or
	parser.PrettyPrintNode(node)

}
