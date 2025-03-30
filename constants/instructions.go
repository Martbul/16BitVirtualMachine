package constants

//INFO: instructions are written in hex

const (
	MOV_LIT_REG = 0x10 //0x10 is opcode(it is 16 decimal). Each instruction needs an uinique identifiesr(opcode) in machine code
	MOV_REG_REG = 0x11
	MOV_REG_MEM = 0x12
	MOV_MEM_REG = 0x13
	ADD_REG_REG = 0x14
	JMP_NOT_EQ  = 0x15
)
