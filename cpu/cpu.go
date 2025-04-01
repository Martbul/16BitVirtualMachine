//INFO: This CPU implementation reads binary instructions from memory, decodes them and executes them. The instructions are stored  as bytes and the memory is undernned a big slice

// 16 bit cpu
// each register is 16 bits wide
package cpu

import (
	"fmt"
	"log"

	"github.com/martbul/constants"
	"github.com/martbul/memory"
	memorymapper "github.com/martbul/memoryMapper"
)

//register hold state in the CPU
//instructuin pointer is a special register
//registers are small peice of memory

type CPU struct {
	memory         *memorymapper.MemoryMapper
	registers      *memory.DataView
	registerMap    map[string]int
	registerNames  []string
	stackFrameSize int
}

func NewCPU(mem *memorymapper.MemoryMapper) *CPU {
	registerNames := []string{
		"ip", "acc",
		"r1", "r2", "r3", "r4",
		"r5", "r6", "r7", "r8",
		"sp", "fp",
	}

	// Create registers memory space with 2 bytes per register
	registers := memory.NewDataView(len(registerNames) * 2)
	registerMap := make(map[string]int)
	for i, name := range registerNames {
		registerMap[name] = i * 2
	}

	stackFrameSize := 0

	cpu := &CPU{
		memory:         mem,
		registerNames:  registerNames,
		registers:      registers,
		registerMap:    registerMap,
		stackFrameSize: stackFrameSize,
	}

	// Stack grows downward, so set SP and FP at the end of memory
	cpu.SetRegister("sp", 0xffff-1)
	cpu.SetRegister("fp", 0xffff-1)

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
	//return: 0x0f01: 0x04 0xA3 0xFE 0x13 0x0
	nDefaultVal := 8

	if len(n) > 0 {
		nDefaultVal = n[0] // Use the provided value if given
	}

	nextNBytes := make([]string, nDefaultVal)

	for i := 0; i < nDefaultVal; i++ {
		val, _ := cpu.memory.GetUint8(address + i)

		nextNBytes[i] = fmt.Sprintf("0X%02X", val)
	}

	fmt.Printf("0x%04X: %s\n", address, nextNBytes)
}

// GetRegister gets the value of a register
func (cpu *CPU) GetRegister(name string) uint16 {
	offset, exists := cpu.registerMap[name]
	if !exists {
		log.Fatalf("getRegister: No such register '%s'", name)
	}
	return cpu.registers.GetUint16(offset)
}

// SetRegister sets the value of a register
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
func (cpu *CPU) Execute(instruction uint8) (bool, string) {
	switch instruction {

	case constants.MOV_LIT_REG:
		literal := cpu.Fetch16()
		register := cpu.FetachRegisterIndex()
		cpu.registers.SetUint16(int(register), literal)
		return false, ""

		//moving register to register
	case constants.MOV_REG_REG:
		registerFrom := cpu.FetachRegisterIndex()
		registerTo := cpu.FetachRegisterIndex()
		value := cpu.registers.GetUint16(int(registerFrom))
		cpu.registers.SetUint16(int(registerTo), value)
		return false, ""

	//move register to memory
	case constants.MOV_REG_MEM:
		registerFrom := cpu.FetachRegisterIndex()
		address := cpu.Fetch16()
		value := cpu.registers.GetUint16(int(registerFrom))
		cpu.memory.SetUint16(int(address), value)

		//	fmt.Printf("MOV_REG_MEM: Writing 0x%X ('%c') to memory address 0x%X\n", value, value, address) // Debug
		return false, ""

	//move memory to register
	case constants.MOV_MEM_REG:
		address := cpu.Fetch16()
		value, _ := cpu.memory.GetUint16(int(address))
		registerTo := cpu.FetachRegisterIndex()
		cpu.registers.SetUint16(int(registerTo), value)
		return false, ""

	// add register to register
	case constants.ADD_REG_REG:
		r1 := cpu.FetachRegisterIndex()
		r2 := cpu.FetachRegisterIndex()

		registerValue1 := cpu.registers.GetUint16(r1)
		registerValue2 := cpu.registers.GetUint16(r2)

		sum := registerValue1 + registerValue2
		cpu.SetRegister("acc", sum)

		///		fmt.Printf("ADD_REG_REG: Added r%d (0x%04X) + r%d (0x%04X) = 0x%04X (stored in acc)\n",
		//			r1, registerValue1, r2, registerValue2, sum)
		return false, ""

	case constants.JMP_NOT_EQ:
		value := cpu.Fetch16()
		address := cpu.Fetch16()

		if value != cpu.GetRegister("acc") {
			cpu.SetRegister("ip", address)
		}
		return false, ""

	//push literal value on the stack
	case constants.PSH_LIT:
		value := cpu.Fetch16()
		cpu.Push(value)
		return false, ""

	//push val from register on the stack
	case constants.PSH_REG:
		registerIndex := cpu.FetachRegisterIndex()
		cpu.Push(cpu.registers.GetUint16(registerIndex))
		return false, ""

	case constants.POP:
		registerIndex := cpu.FetachRegisterIndex()
		value := cpu.Pop()
		cpu.registers.SetUint16(registerIndex, value)
		return false, ""

	case constants.CAL_LIT:
		address := cpu.Fetch16()

		//saving the cpu state
		cpu.PushState()

		cpu.SetRegister("ip", address)
		return false, ""

	case constants.CAL_REG:
		registerIndex := cpu.FetachRegisterIndex()
		address := cpu.registers.GetUint16(registerIndex)
		cpu.PushState()
		cpu.SetRegister("ip", address)
		return false, ""

		//return from subroutinw
	case constants.RET:
		cpu.PopState()
		return false, ""

	case constants.HLT:
		return true, "manual hlt"

		// default case to handle unknown instructions
	default:
		fmt.Printf("Unknown instruction: 0x%X at address: 0x%X\n", instruction, cpu.GetRegister("ip")-1)
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
			//fmt.Println(hltReason)
			break
		}
	}
}
