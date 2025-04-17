package instructions

//INFO: instructions are written in hex

// InstructionType represents the type of CPU instruction
type InstructionType int

// InstructionSize represents the size of an instruction type in bytes
type InstructionSize int

const (
	LitReg InstructionType = iota
	RegLit
	RegLit8
	RegReg
	RegMem
	MemReg
	LitMem
	RegPtrReg
	LitOffReg
	NoArgs
	SingleReg
	SingleLit
)

// Define instruction sizes for each type
const (
	SizeNoArgs    InstructionSize = 1
	SizeSingleReg InstructionSize = 2
	SizeSingleLit InstructionSize = 3
	SizeRegReg    InstructionSize = 3
	SizeRegMem    InstructionSize = 3
	SizeMemReg    InstructionSize = 3
	SizeRegLit    InstructionSize = 4
	SizeRegLit8   InstructionSize = 3
	SizeRegPtrReg InstructionSize = 3
	SizeLitReg    InstructionSize = 4 // opcode - 1 byte, literal - 2 bytes, reg - 1 byte= 4 total bytes
	SizeLitMem    InstructionSize = 5
	SizeLitOffReg InstructionSize = 5
)

const (
	MOV_LIT_REG     = 0x10 //0x10 is opcode(it is 16 decimal). Each instruction needs an uinique identifiesr(opcode) in machine code
	MOV_REG_REG     = 0x11
	MOV_REG_MEM     = 0x12
	MOV_MEM_REG     = 0x13
	MOV_LIT_MEM     = 0x1B
	MOV_REG_PTR_REG = 0x1C
	MOV_LIT_OFF_REG = 0x1D

	ADD_REG_REG = 0x14
	ADD_LIT_REG = 0x3F
	SUB_LIT_REG = 0x16
	SUB_REG_LIT = 0x1E
	SUB_REG_REG = 0x1F
	INC_REG     = 0x35
	DEC_REG     = 0x36
	MUL_LIT_REG = 0x20
	MUL_REG_REG = 0x21

	LSF_REG_LIT = 0x26
	LSF_REG_REG = 0x27
	RSF_REG_LIT = 0x2A
	RSF_REG_REG = 0x2B
	AND_REG_LIT = 0x2E
	AND_REG_REG = 0x2F
	OR_REG_LIT  = 0x30
	OR_REG_REG  = 0x31
	XOR_REG_LIT = 0x32
	XOR_REG_REG = 0x33
	NOT         = 0x34

	JMP_NOT_EQ = 0x15
	JNE_REG    = 0x40
	JEQ_REG    = 0x3E
	JEQ_LIT    = 0x41
	JLT_REG    = 0x42
	JLT_LIT    = 0x43
	JGT_REG    = 0x44
	JGT_LIT    = 0x45
	JLE_REG    = 0x46
	JLE_LIT    = 0x47
	JGE_REG    = 0x48
	JGE_LIT    = 0x49

	PSH_LIT = 0x17
	PSH_REG = 0x18
	POP     = 0x1A
	CAL_LIT = 0x5E
	CAL_REG = 0x5F
	RET     = 0x60
	HLT     = 0xFF
)

type MetaData struct {
	Instruction string
	Opcode      byte
	Type        InstructionType
	Size        InstructionSize
	Mnemonic    string
}

var Instructions = []MetaData{
	{Instruction: "MOV_LIT_REG", Opcode: 0x10, Type: LitReg, Size: SizeLitReg, Mnemonic: "mov"},
	{Instruction: "MOV_REG_REG", Opcode: 0x11, Type: RegReg, Size: SizeRegReg, Mnemonic: "mov"},
	{Instruction: "MOV_REG_MEM", Opcode: 0x12, Type: RegMem, Size: SizeRegMem, Mnemonic: "mov"},
	{Instruction: "MOV_MEM_REG", Opcode: 0x13, Type: MemReg, Size: SizeMemReg, Mnemonic: "mov"},
	{Instruction: "MOV_LIT_MEM", Opcode: 0x1B, Type: LitMem, Size: SizeLitMem, Mnemonic: "mov"},
	{Instruction: "MOV_REG_PTR_REG", Opcode: 0x1C, Type: RegPtrReg, Size: SizeRegPtrReg, Mnemonic: "mov"},
	{Instruction: "MOV_LIT_OFF_REG", Opcode: 0x1D, Type: LitOffReg, Size: SizeLitOffReg, Mnemonic: "mov"},
	{Instruction: "ADD_REG_REG", Opcode: 0x14, Type: RegReg, Size: SizeRegReg, Mnemonic: "add"},
	{Instruction: "ADD_LIT_REG", Opcode: 0x3F, Type: LitReg, Size: SizeLitReg, Mnemonic: "add"},
	{Instruction: "SUB_LIT_REG", Opcode: 0x16, Type: LitReg, Size: SizeLitReg, Mnemonic: "sub"},
	{Instruction: "SUB_REG_LIT", Opcode: 0x1E, Type: RegLit, Size: SizeRegLit, Mnemonic: "sub"},
	{Instruction: "SUB_REG_REG", Opcode: 0x1F, Type: RegReg, Size: SizeRegReg, Mnemonic: "sub"},
	{Instruction: "INC_REG", Opcode: 0x35, Type: SingleReg, Size: SizeSingleReg, Mnemonic: "inc"},
	{Instruction: "DEC_REG", Opcode: 0x36, Type: SingleReg, Size: SizeSingleReg, Mnemonic: "dec"},
	{Instruction: "MUL_LIT_REG", Opcode: 0x20, Type: LitReg, Size: SizeLitReg, Mnemonic: "mul"},
	{Instruction: "MUL_REG_REG", Opcode: 0x21, Type: RegReg, Size: SizeRegReg, Mnemonic: "mul"},
	{Instruction: "LSF_REG_LIT", Opcode: 0x26, Type: RegLit8, Size: SizeRegLit8, Mnemonic: "lsf"},
	{Instruction: "LSF_REG_REG", Opcode: 0x27, Type: RegReg, Size: SizeRegReg, Mnemonic: "lsf"},
	{Instruction: "RSF_REG_LIT", Opcode: 0x2A, Type: RegLit8, Size: SizeRegLit8, Mnemonic: "rsf"},
	{Instruction: "RSF_REG_REG", Opcode: 0x2B, Type: RegReg, Size: SizeRegReg, Mnemonic: "rsf"},
	{Instruction: "AND_REG_LIT", Opcode: 0x2E, Type: RegLit, Size: SizeRegLit, Mnemonic: "and"},
	{Instruction: "AND_REG_REG", Opcode: 0x2F, Type: RegReg, Size: SizeRegReg, Mnemonic: "and"},
	{Instruction: "OR_REG_LIT", Opcode: 0x30, Type: RegLit, Size: SizeRegLit, Mnemonic: "or"},
	{Instruction: "OR_REG_REG", Opcode: 0x31, Type: RegReg, Size: SizeRegReg, Mnemonic: "or"},
	{Instruction: "XOR_REG_LIT", Opcode: 0x32, Type: RegLit, Size: SizeRegLit, Mnemonic: "xor"},
	{Instruction: "XOR_REG_REG", Opcode: 0x33, Type: RegReg, Size: SizeRegReg, Mnemonic: "xor"},
	{Instruction: "NOT", Opcode: 0x34, Type: SingleReg, Size: SizeSingleReg, Mnemonic: "not"},
	{Instruction: "JMP_NOT_EQ", Opcode: 0x15, Type: LitMem, Size: SizeLitMem, Mnemonic: "jne"},
	{Instruction: "JNE_REG", Opcode: 0x40, Type: RegMem, Size: SizeRegMem, Mnemonic: "jne"},
	{Instruction: "JEQ_REG", Opcode: 0x3E, Type: RegMem, Size: SizeRegMem, Mnemonic: "jeq"},
	{Instruction: "JEQ_LIT", Opcode: 0x41, Type: LitMem, Size: SizeLitMem, Mnemonic: "jeq"},
	{Instruction: "JLT_REG", Opcode: 0x42, Type: RegMem, Size: SizeRegMem, Mnemonic: "jlt"},
	{Instruction: "JLT_LIT", Opcode: 0x43, Type: LitMem, Size: SizeLitMem, Mnemonic: "jlt"},
	{Instruction: "JGT_REG", Opcode: 0x44, Type: RegMem, Size: SizeRegMem, Mnemonic: "jgt"},
	{Instruction: "JGT_LIT", Opcode: 0x45, Type: LitMem, Size: SizeLitMem, Mnemonic: "jgt"},
	{Instruction: "JLE_REG", Opcode: 0x46, Type: RegMem, Size: SizeRegMem, Mnemonic: "jle"},
	{Instruction: "JLE_LIT", Opcode: 0x47, Type: LitMem, Size: SizeLitMem, Mnemonic: "jle"},
	{Instruction: "JGE_REG", Opcode: 0x48, Type: RegMem, Size: SizeRegMem, Mnemonic: "jge"},
	{Instruction: "JGE_LIT", Opcode: 0x49, Type: LitMem, Size: SizeLitMem, Mnemonic: "jge"},
	{Instruction: "PSH_LIT", Opcode: 0x17, Type: SingleLit, Size: SizeSingleLit, Mnemonic: "psh"},
	{Instruction: "PSH_REG", Opcode: 0x18, Type: SingleReg, Size: SizeSingleReg, Mnemonic: "psh"},
	{Instruction: "POP", Opcode: 0x1A, Type: SingleReg, Size: SizeSingleReg, Mnemonic: "pop"},
	{Instruction: "CAL_LIT", Opcode: 0x5E, Type: SingleLit, Size: SizeSingleLit, Mnemonic: "cal"},
	{Instruction: "CAL_REG", Opcode: 0x5F, Type: SingleReg, Size: SizeSingleReg, Mnemonic: "cal"},
	{Instruction: "RET", Opcode: 0x60, Type: NoArgs, Size: SizeNoArgs, Mnemonic: "ret"},
	{Instruction: "HLT", Opcode: 0xFF, Type: NoArgs, Size: SizeNoArgs, Mnemonic: "hlt"},
}

var InstructionMap map[byte]MetaData

func init() {
	InstructionMap = make(map[byte]MetaData)
	for _, inst := range Instructions {
		InstructionMap[inst.Opcode] = inst
	}
}

func GetInstructionName(opcode byte) string {
	if inst, ok := InstructionMap[opcode]; ok {
		return inst.Instruction
	}
	return "UNKNOWN"
}

func GetInstructionSize(opcode byte) InstructionSize {
	if inst, ok := InstructionMap[opcode]; ok {
		return inst.Size
	}
	return 0
}
