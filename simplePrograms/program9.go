package simpleprograms

import "github.com/martbul/assembler"

func SimpleProgram9() {

	//program :=
	//	`constant code_constant = $C0DE
	//	+data8 bytes = { $01,   $02,   $03,   $04   }
	//	data16 words = { $0506, $0708, $090A, $0B0C }
	//	code:
	//	mov [!code_constant], &1234`

	program := `constant code_constant = $C0DE
mov [!code_constant], &1234`

	assembler.AssembleProgram(program)
}
