package simpleprograms

import (
	"fmt"

	"github.com/martbul/assembler/parser"
)

func SimpleProgram6() {
	assembly := "mov $42, r4"
	node, err := parser.ParseMovLitToReg(assembly)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Print the node (you would implement proper JSON/YAML pretty printing)
	fmt.Printf("%+v\n", node)

}
