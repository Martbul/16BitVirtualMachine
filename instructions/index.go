// creating a lookup table (dictionary) from the array of instructions
//INFO: The end result is a function getInstructionByName where you can call for example MOV_LIT_REG and you receive its meadata.

package instructions

var InstructionByName map[string]MetaData
var InstructionByOpcode map[uint8]MetaData

func init() {
	// Initialize both maps
	InstructionByName = make(map[string]MetaData)
	InstructionByOpcode = make(map[uint8]MetaData)

	for _, inst := range Instructions {
		InstructionByName[inst.Instruction] = inst
		InstructionByOpcode[inst.Opcode] = inst
	}
}

func GetInstructionByName(name string) (MetaData, bool) {
	inst, ok := InstructionByName[name]
	return inst, ok
}
func GetInstructionByOpcode(opcode uint8) (MetaData, bool) {
	inst, ok := InstructionByOpcode[opcode]
	return inst, ok
}

// Example usage:
// instruction, found := GetInstructionByName("MOV_LIT_REG")
// if found {
//     fmt.Printf("Opcode: 0x%02X\n", instruction.Opcode)
// }
