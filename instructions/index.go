package instructions

var InstructionByName map[string]MetaData

func init() {
	InstructionByName = make(map[string]MetaData)
	for _, inst := range Instructions {
		InstructionByName[inst.Instruction] = inst
	}
}

func GetInstructionByName(name string) (MetaData, bool) {
	inst, ok := InstructionByName[name]
	return inst, ok
}

// Example usage:
// instruction, found := GetInstructionByName("MOV_LIT_REG")
// if found {
//     fmt.Printf("Opcode: 0x%02X\n", instruction.Opcode)
// }
