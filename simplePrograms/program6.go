package simpleprograms

import (
	"fmt"
	"os"

	"github.com/martbul/assembler/parser"
)

func SimpleProgram6() {
	//	input := "mov [$42 + !loc - ($05 * ($31 + !var) - $07)], r4"
	//WARN: input := "mov acc, &[!loc + $4200]" - doesn not work(mov reg mem)
	input := "mov &4200, r1"
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
