package main

import (
	"github.com/martbul/assembler"
)

func main() {

	// Example program (as in your JS code)
	program := []string{
		"start:",
		"  mov $0A, &0050",
		"loop:",
		"  mov &0050, acc",
		"  dec acc",
		"  mov acc, &0050",
		"  inc r2",
		"  inc r2",
		"  inc r2",
		"  jne $00, &[!loop]",
		"end:",
		"  hlt",
	}

	assembler.AssembleProgram(program)

}
