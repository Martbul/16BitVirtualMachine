//INFO: This CPU implementation reads binary instructions from memory, decodes them and executes them. The instructions are stored  as bytes and the memory is undernned a big slice

// 16 bit cpu
// each register is 16 bits wide
package cpu

import (
	"fmt"
	"log"

	"github.com/martbul/instructions"
	"github.com/martbul/memory"
	memorymapper "github.com/martbul/memoryMapper"
	"github.com/martbul/registers"
)

//register hold state in the CPU
//instructuin pointer is a special register
//registers are small peice of memory

type CPU struct {
	memory                *memorymapper.MemoryMapper
	registers             *memory.DataView
	registerMap           map[string]int
	registerNames         []string
	stackFrameSize        int
	interuptVectorAddress int
	isInInteruptedHandler bool
}

func NewCPU(mem *memorymapper.MemoryMapper, interuptVectorAddress ...int) *CPU {
	vector := 0xFFFE // default interupVectorAddress
	if len(interuptVectorAddress) > 0 {
		vector = interuptVectorAddress[0]
	}
	registerNames := registers.Registers
	// Create registers memory space with 2 bytes per register
	registers := memory.NewDataView(len(registerNames) * 2)
	registerMap := make(map[string]int)
	for i, name := range registerNames {
		registerMap[name] = i * 2
	}

	stackFrameSize := 0

	cpu := &CPU{
		memory:                mem,
		registerNames:         registerNames,
		registers:             registers,
		registerMap:           registerMap,
		stackFrameSize:        stackFrameSize,
		interuptVectorAddress: vector,
		isInInteruptedHandler: false,
	}

	// Stack grows downward, so set SP and FP at the end of memory
	cpu.SetRegister("sp", 0xffff-1)
	cpu.SetRegister("fp", 0xffff-1)

	cpu.SetRegister("im", 0xffff)

	return cpu
}

func (cpu *CPU) Debug() {
	for _, name := range cpu.registerNames {
		value := cpu.GetRegister(name)
		fmt.Printf("%s: 0x%04X\n", name, value)
	}
	fmt.Println()
}

func (cpu *CPU) ViewMemoryAt(address int, n ...int) {
	nDefaultVal := 8

	if len(n) > 0 {
		nDefaultVal = n[0]
	}

	nextNBytes := make([]string, nDefaultVal)

	for i := 0; i < nDefaultVal; i++ {
		val, _ := cpu.memory.GetUint8(address + i)

		nextNBytes[i] = fmt.Sprintf("0X%02X", val)
	}

	fmt.Printf("0x%04X: %s\n", address, nextNBytes)
}

func (cpu *CPU) GetRegister(name string) uint16 {
	offset, exists := cpu.registerMap[name]
	if !exists {
		log.Fatalf("getRegister: No such register '%s'", name)
	}
	return cpu.registers.GetUint16(offset)
}

func (cpu *CPU) SetRegister(name string, value uint16) {
	offset, exists := cpu.registerMap[name]
	if !exists {
		log.Fatalf("setRegister: No such register '%s'", name)
	}
	cpu.registers.SetUint16(offset, value)
}

// when fetch() is called it fetches the instruction and moves the instruction pointer 1 byte
func (cpu *CPU) Fetch() uint8 {
	// Get the current instruction pointer
	nextInstructionAddress := cpu.GetRegister("ip")

	// Fetch the instruction (byte) from memory at the current address
	instruction, _ := cpu.memory.GetUint8(int(nextInstructionAddress))

	cpu.SetRegister("ip", nextInstructionAddress+1)

	return instruction
}

// fetches the 16-bit instruction where the instruction pointer is located
func (cpu *CPU) Fetch16() uint16 {
	nextInstructionAddress := cpu.GetRegister("ip")
	instruction, _ := cpu.memory.GetUint16(int(nextInstructionAddress))
	cpu.SetRegister("ip", nextInstructionAddress+2)
	return instruction
}

func (cpu *CPU) FetachRegisterIndex() int {
	return (int(cpu.Fetch()) % len(cpu.registerNames)) * 2
}

func (cpu *CPU) Push(value uint16) {
	spAddress := cpu.GetRegister("sp")

	cpu.memory.SetUint16(int(spAddress), value)
	cpu.SetRegister("sp", spAddress-2)
	cpu.stackFrameSize += 2
}
func (cpu *CPU) Pop() uint16 {

	nextSpAddress := cpu.GetRegister("sp") + 2
	cpu.SetRegister("sp", nextSpAddress)
	value, _ := cpu.memory.GetUint16(int(nextSpAddress))
	cpu.stackFrameSize -= 2
	return value
}

// Pushes the cpu's state onto the stack
func (cpu *CPU) PushState() {

	//saving the cpu state
	cpu.Push(cpu.GetRegister("r1"))
	cpu.Push(cpu.GetRegister("r2"))
	cpu.Push(cpu.GetRegister("r3"))
	cpu.Push(cpu.GetRegister("r4"))
	cpu.Push(cpu.GetRegister("r5"))
	cpu.Push(cpu.GetRegister("r6"))
	cpu.Push(cpu.GetRegister("r7"))
	cpu.Push(cpu.GetRegister("r8"))
	cpu.Push(cpu.GetRegister("ip")) //INFO: ip is the return address of this subroutine
	cpu.Push(uint16(cpu.stackFrameSize) + 2)

	cpu.SetRegister("fp", cpu.GetRegister("sp")) //INFO: Moving the framePointer to where the stackPointer points
	cpu.stackFrameSize = 0                       //INFO: Reseting the stackFrameSize
}

func (cpu *CPU) PopState() {
	framePointerAddress := cpu.GetRegister("fp")
	cpu.SetRegister("sp", framePointerAddress)

	cpu.stackFrameSize = int(cpu.Pop())
	stackFrameSize := cpu.stackFrameSize

	cpu.SetRegister("ip", cpu.Pop())
	cpu.SetRegister("r8", cpu.Pop())
	cpu.SetRegister("r7", cpu.Pop())
	cpu.SetRegister("r6", cpu.Pop())
	cpu.SetRegister("r5", cpu.Pop())
	cpu.SetRegister("r4", cpu.Pop())
	cpu.SetRegister("r3", cpu.Pop())
	cpu.SetRegister("r2", cpu.Pop())
	cpu.SetRegister("r1", cpu.Pop())

	nArgs := cpu.Pop()

	for i := 0; i < int(nArgs); i++ {
		cpu.Pop()
	}

	cpu.SetRegister("fp", framePointerAddress+uint16(stackFrameSize))
}

// Execute decodes and executes instructions
func (cpu *CPU) Execute(instr uint8) (bool, string) {
	// Look up the instruction by opcode
	_, found := instructions.GetInstructionByOpcode(instr)
	if !found {
		return true, fmt.Sprintf("Unknown instruction: 0x%X", instr)
	}
	switch instr {

	//software triggered interupt
	case instructions.INT:
		interuptValue := cpu.Fetch16()
		cpu.HandleInterupt(interuptValue)
		return false, ""
	case instructions.RET_INT:
		cpu.isInInteruptedHandler = false
		cpu.PopState()
		return false, ""

	case instructions.MOV_LIT_REG:
		literal := cpu.Fetch16()
		register := cpu.FetachRegisterIndex()
		cpu.registers.SetUint16(int(register), literal)
		return false, ""

		//moving register to register
	case instructions.MOV_REG_REG:
		registerFrom := cpu.FetachRegisterIndex()
		registerTo := cpu.FetachRegisterIndex()
		value := cpu.registers.GetUint16(int(registerFrom))
		cpu.registers.SetUint16(int(registerTo), value)
		return false, ""

	//move register to memory
	case instructions.MOV_REG_MEM:
		registerFrom := cpu.FetachRegisterIndex()
		address := cpu.Fetch16()
		value := cpu.registers.GetUint16(int(registerFrom))
		cpu.memory.SetUint16(int(address), value)

		//	fmt.Printf("MOV_REG_MEM: Writing 0x%X ('%c') to memory address 0x%X\n", value, value, address) // Debug
		return false, ""

	//move memory to register
	case instructions.MOV_MEM_REG:
		address := cpu.Fetch16()
		value, _ := cpu.memory.GetUint16(int(address))
		registerTo := cpu.FetachRegisterIndex()
		cpu.registers.SetUint16(int(registerTo), value)
		return false, ""

	case instructions.MOV_LIT_MEM:
		value := cpu.Fetch16()
		address := cpu.Fetch16()
		cpu.memory.SetUint16(int(address), value)
		return false, ""

	case instructions.MOV_REG_PTR_REG:
		r1 := cpu.FetachRegisterIndex() //register we want the vaklue form
		r2 := cpu.FetachRegisterIndex() //destination regisster
		ptr := cpu.registers.GetUint16(r1)
		value, _ := cpu.memory.GetUint16(int(ptr)) //WARN: Bad error hnadlg fix later
		cpu.registers.SetUint16(r2, value)
		return false, ""

	// move value at [literal + register ] to register
	case instructions.MOV_LIT_OFF_REG:
		baseAddress := cpu.Fetch16()
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex() //destination reg
		offset := cpu.registers.GetUint16(r1)

		value, _ := cpu.memory.GetUint16(int(baseAddress) + int(offset))
		cpu.registers.SetUint16(r2, value)
		return false, ""

	// add register to register
	case instructions.ADD_REG_REG:
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex()

		registerValue1 := cpu.registers.GetUint16(r1)
		registerValue2 := cpu.registers.GetUint16(r2)

		sum := registerValue1 + registerValue2
		cpu.SetRegister("acc", sum)

		///		fmt.Printf("ADD_REG_REG: Added r%d (0x%04X) + r%d (0x%04X) = 0x%04X (stored in acc)\n",
		//			r1, registerValue1, r2, registerValue2, sum)
		return false, ""

	//add literal value to register value
	case instructions.ADD_LIT_REG:
		literal := cpu.Fetch16()
		r1 := cpu.FetachRegisterIndex()
		registerValue := cpu.registers.GetUint16(r1)
		cpu.SetRegister("acc", literal+registerValue)
		return false, ""

	case instructions.SUB_LIT_REG:
		literal := cpu.Fetch16()
		r1 := cpu.FetachRegisterIndex()
		registerValue := cpu.registers.GetUint16(r1)
		res := registerValue - literal
		cpu.SetRegister("acc", res)
		return false, ""

	// subtract register value from a literal value
	case instructions.SUB_REG_LIT:
		r1 := cpu.FetachRegisterIndex()
		registerValue := cpu.registers.GetUint16(r1)
		literal := cpu.Fetch16()
		res := literal - registerValue
		cpu.SetRegister("acc", res)
		return false, ""

	case instructions.SUB_REG_REG:
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex()
		registerValue1 := cpu.registers.GetUint16(r1)
		registerValue2 := cpu.registers.GetUint16(r2)
		res := registerValue1 - registerValue2
		cpu.SetRegister("acc", res)
		return false, ""

	case instructions.MUL_LIT_REG:
		literal := cpu.Fetch16()
		r1 := cpu.FetachRegisterIndex()
		r1Value := cpu.registers.GetUint16(r1)
		res := literal * r1Value
		cpu.SetRegister("acc", res)
		return false, ""

	case instructions.MUL_REG_REG:
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex()
		registerValue1 := cpu.registers.GetUint16(r1)
		registerValue2 := cpu.registers.GetUint16(r2)
		res := registerValue1 * registerValue2
		cpu.SetRegister("acc", res)
		return false, ""

	case instructions.INC_REG:
		r1 := cpu.FetachRegisterIndex()
		oldValue := cpu.registers.GetUint16(r1)
		newValue := oldValue + 1
		cpu.registers.SetUint16(r1, newValue)
		return false, ""

	case instructions.DEC_REG:
		r1 := cpu.FetachRegisterIndex()
		oldValue := cpu.registers.GetUint16(r1)
		newValue := oldValue - 1
		cpu.registers.SetUint16(r1, newValue)
		return false, ""

	// left shift register by literal value (in place)
	//INFO: Any bits that are outside the boudnaries of 16 are lost
	//in left shif 9 << 2 is equal to 9 * 2^2
	case instructions.LSF_REG_LIT:
		r1 := cpu.FetachRegisterIndex()
		literal := cpu.Fetch16()
		registerValue := cpu.registers.GetUint16(r1)
		res := registerValue << literal
		cpu.registers.SetUint16(r1, res)
		return false, ""

	// left shift register by register (in place)
	case instructions.LSF_REG_REG:
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex()
		registerValue1 := cpu.registers.GetUint16(r1)
		registerValue2 := cpu.registers.GetUint16(r2)
		res := registerValue1 << registerValue2
		cpu.registers.SetUint16(r1, res)
		return false, ""

	//INFO: in right shift 9 >> 2 is equal to 9 / 2^2
	case instructions.RSF_REG_LIT:
		r1 := cpu.FetachRegisterIndex()
		literal := cpu.Fetch16()
		registerValue := cpu.registers.GetUint16(r1)
		res := registerValue >> literal
		cpu.registers.SetUint16(r1, res)
		return false, ""

	// right shift register by register (in place)
	case instructions.RSF_REG_REG:
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex()
		registerValue1 := cpu.registers.GetUint16(r1)
		registerValue2 := cpu.registers.GetUint16(r2)
		res := registerValue1 >> registerValue2
		cpu.registers.SetUint16(r1, res)
		return false, ""

	// and register with literal
	//INFO: AND takes 2 binary numbers and produce a new binary number wher ieatch place is one where both of the numbers BOTH of the 2 binariy nums have a 1 in the same place, otherwise it is 0(useful for isolating a particular part of a number like the bottom or the top byte)
	case instructions.AND_REG_LIT:
		r1 := cpu.FetachRegisterIndex()
		literal := cpu.Fetch16()
		registerValue := cpu.registers.GetUint16(r1)

		res := registerValue & literal
		cpu.SetRegister("acc", res)
		return false, ""

	case instructions.AND_REG_REG:
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex()
		registerValue1 := cpu.registers.GetUint16(r1)
		registerValue2 := cpu.registers.GetUint16(r2)

		res := registerValue1 & registerValue2
		cpu.SetRegister("acc", res)
		return false, ""

	//INFO: OR takes 2 binary numbers and produce a new binary number where each place is a one if eather of the numbers have a 1 in that place, otherwise it is a 0.
	case instructions.OR_REG_LIT:
		r1 := cpu.FetachRegisterIndex()
		literal := cpu.Fetch16()
		registerValue := cpu.registers.GetUint16(r1)

		res := registerValue | literal
		cpu.SetRegister("acc", res)
		return false, ""

	case instructions.OR_REG_REG:
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex()
		registerValue1 := cpu.registers.GetUint16(r1)
		registerValue2 := cpu.registers.GetUint16(r2)

		res := registerValue1 | registerValue2
		cpu.SetRegister("acc", res)
		return false, ""

	//INFO: XOR works almost like OR, exapt what it is exclusive -> each place is a 1 if ever a number is a 1, if both are 1 or if both are 0 it is a 0.
	// A ^ B = C
	// A ^ C = B
	// B ^ C = A
	case instructions.XOR_REG_LIT:
		r1 := cpu.FetachRegisterIndex()
		literal := cpu.Fetch16()
		registerValue := cpu.registers.GetUint16(r1)

		res := registerValue ^ literal
		cpu.SetRegister("acc", res)
		return false, ""

	case instructions.XOR_REG_REG:
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex()
		registerValue1 := cpu.registers.GetUint16(r1)
		registerValue2 := cpu.registers.GetUint16(r2)

		res := registerValue1 ^ registerValue2
		cpu.SetRegister("acc", res)
		return false, ""

	//INFO: NOT flips all bits in a number 0->1 and 1->0
	case instructions.NOT:
		r1 := cpu.FetachRegisterIndex()
		registerValue := cpu.registers.GetUint16(r1)
		res := ^registerValue & 0xffff // selecting just the bottmo 16 bits so that the number is not abouve the 16 bits the vm works on
		cpu.SetRegister("acc", res)
		return false, ""

		//jump if literal not equal
	case instructions.JMP_NOT_EQ:
		value := cpu.Fetch16()
		address := cpu.Fetch16()

		if value != cpu.GetRegister("acc") {
			cpu.SetRegister("ip", address)
		}
		return false, ""

		//jump if register not equal
	case instructions.JNE_REG:
		r1 := cpu.FetachRegisterIndex()
		r1Value := cpu.registers.GetUint16(r1)
		addressToJupmTo := cpu.Fetch16()

		if r1Value != cpu.GetRegister("acc") {
			cpu.SetRegister("ip", addressToJupmTo)
		}
		return false, ""

		//jump if literal equal
	case instructions.JEQ_LIT:
		value := cpu.Fetch16()
		address := cpu.Fetch16()

		if value == cpu.GetRegister("acc") {
			cpu.SetRegister("ip", address)
		}
		return false, ""

		//jump if register equal
	case instructions.JEQ_REG:
		r1 := cpu.FetachRegisterIndex()
		r1Value := cpu.registers.GetUint16(r1)
		addressToJupmTo := cpu.Fetch16()

		if r1Value == cpu.GetRegister("acc") {
			cpu.SetRegister("ip", addressToJupmTo)
		}
		return false, ""

		//jump if literal less than
	case instructions.JLT_LIT:
		value := cpu.Fetch16()
		address := cpu.Fetch16()

		if value < cpu.GetRegister("acc") {
			cpu.SetRegister("ip", address)
		}
		return false, ""

		//jump if register less than
	case instructions.JLT_REG:
		r1 := cpu.FetachRegisterIndex()
		r1Value := cpu.registers.GetUint16(r1)
		addressToJupmTo := cpu.Fetch16()

		if r1Value < cpu.GetRegister("acc") {
			cpu.SetRegister("ip", addressToJupmTo)
		}
		return false, ""

		//jump if literal greater than
	case instructions.JGT_LIT:
		value := cpu.Fetch16()
		address := cpu.Fetch16()

		if value > cpu.GetRegister("acc") {
			cpu.SetRegister("ip", address)
		}
		return false, ""

		//jump if register greater than
	case instructions.JGT_REG:
		r1 := cpu.FetachRegisterIndex()
		r1Value := cpu.registers.GetUint16(r1)
		addressToJupmTo := cpu.Fetch16()

		if r1Value > cpu.GetRegister("acc") {
			cpu.SetRegister("ip", addressToJupmTo)
		}
		return false, ""

		//jump if literal less than or equal to
	case instructions.JLE_LIT:
		value := cpu.Fetch16()
		address := cpu.Fetch16()

		if value <= cpu.GetRegister("acc") {
			cpu.SetRegister("ip", address)
		}
		return false, ""

		//jump if register less than or equal to
	case instructions.JLE_REG:
		r1 := cpu.FetachRegisterIndex()
		r1Value := cpu.registers.GetUint16(r1)
		addressToJupmTo := cpu.Fetch16()

		if r1Value <= cpu.GetRegister("acc") {
			cpu.SetRegister("ip", addressToJupmTo)
		}
		return false, ""

		//jump if literal greater than or equal to
	case instructions.JGE_LIT:
		value := cpu.Fetch16()
		address := cpu.Fetch16()

		if value >= cpu.GetRegister("acc") {
			cpu.SetRegister("ip", address)
		}
		return false, ""

		//jump if register greater than or equal to
	case instructions.JGE_REG:
		r1 := cpu.FetachRegisterIndex()
		r1Value := cpu.registers.GetUint16(r1)
		addressToJupmTo := cpu.Fetch16()

		if r1Value >= cpu.GetRegister("acc") {
			cpu.SetRegister("ip", addressToJupmTo)
		}
		return false, ""

	//push literal value on the stack
	case instructions.PSH_LIT:
		value := cpu.Fetch16()
		cpu.Push(value)
		return false, ""

	//push val from register on the stack
	case instructions.PSH_REG:
		registerIndex := cpu.FetachRegisterIndex()
		cpu.Push(cpu.registers.GetUint16(registerIndex))
		return false, ""

	case instructions.POP:
		registerIndex := cpu.FetachRegisterIndex()
		value := cpu.Pop()
		cpu.registers.SetUint16(registerIndex, value)
		return false, ""

	case instructions.CAL_LIT:
		address := cpu.Fetch16()

		//saving the cpu state
		cpu.PushState()

		cpu.SetRegister("ip", address)
		return false, ""

	case instructions.CAL_REG:
		registerIndex := cpu.FetachRegisterIndex()
		address := cpu.registers.GetUint16(registerIndex)
		cpu.PushState()
		cpu.SetRegister("ip", address)
		return false, ""

		//return from subroutinw
	case instructions.RET:
		cpu.PopState()
		return false, ""

	case instructions.HLT:
		return true, "manual hlt"

		// default case to handle unknown instructions
	default:
		fmt.Printf("Unknown instruction: 0x%X at address: 0x%X\n", instr, cpu.GetRegister("ip")-1)
		return true, "error hlt"

	}
}

func (cpu *CPU) Step() (bool, string) {
	instruction := cpu.Fetch()
	isRunning, hltReason := cpu.Execute(instruction)
	return isRunning, hltReason
}

func (cpu *CPU) Run() {
	for {
		isRunning, _ := cpu.Step()
		if isRunning {
			break
		}
	}
}

func (cpu *CPU) HandleInterupt(value uint16) {
	//INFO: 0xf = 15 in decimal
	interuptVectorIndex := value % 0xf

	isUnmasked := (1<<interuptVectorIndex)&cpu.GetRegister("im") != 0

	if !isUnmasked {
		return
	}

	addressPoninter := cpu.interuptVectorAddress + (int(value) * 2)
	address, err := cpu.memory.GetUint16(addressPoninter)
	if err != nil {
		fmt.Println(err)
	}

	if !cpu.isInInteruptedHandler {
		// push 0 to the stack for the number of arguments
		// 0 indicates that there are no arguments beeing pased via the stack
		cpu.Push(0)
		cpu.PushState()
	}

	cpu.isInInteruptedHandler = true
	cpu.SetRegister("ip", address) // seting the instruction ponter to the of the interup vectoe

}
