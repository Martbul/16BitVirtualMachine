//INFO: This CPU implementation reads binary instructions from memory, decodes them and executes them. The instructions are stored  as bytes and the memory is undernned a big slice

// 16 bit cpu
// each register is 16 bits wide
package cpu

import (
	"fmt"
	"log"

	"github.com/martbul/constants"
	"github.com/martbul/memory"
)

//register hold state in the CPU
//instructuin pointer is a special register
//registers are small peice of memory

type CPU struct {
	memory        *memory.DataView
	registers     *memory.DataView
	registerMap   map[string]int
	registerNames []string
}

func NewCPU(mem *memory.DataView) *CPU {
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

	cpu := &CPU{
		memory:        mem,
		registerNames: registerNames,
		registers:     registers,
		registerMap:   registerMap,
	}
	// Stack grows downward, so set SP and FP at the end of memory
	stackStart := len(mem.GetBuffer()) - 1 - 1
	cpu.SetRegister("sp", uint16(stackStart))
	cpu.SetRegister("fp", uint16(stackStart))

	return cpu
}

func (cpu *CPU) Debug() {
	for _, name := range cpu.registerNames {
		value := cpu.GetRegister(name)
		fmt.Printf("%s: 0x%04X\n", name, value)
	}
	fmt.Println()
}

func (cpu *CPU) ViewMemoryAt(address int) {
	//return: 0x0f01: 0x04 0xA3 0xFE 0x13 0x0

	nextEightBytes := make([]string, 8)

	for i := 0; i < 8; i++ {

		nextEightBytes[i] = fmt.Sprintf("0X%02X", cpu.memory.GetUint8(address+i))
	}

	fmt.Printf("0x%04X: %s\n", address, nextEightBytes)
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
	instruction := cpu.memory.GetUint8(int(nextInstructionAddress))

	cpu.SetRegister("ip", nextInstructionAddress+1)

	return instruction
}

// Fetch16 fetches a 16-bit instruction from memory
func (cpu *CPU) Fetch16() uint16 {
	nextInstructionAddress := cpu.GetRegister("ip")
	instruction := cpu.memory.GetUint16(int(nextInstructionAddress))
	cpu.SetRegister("ip", nextInstructionAddress+2)
	return instruction
}

func (cpu *CPU) Push(value uint16) {
	spAddress := cpu.GetRegister("sp")

	cpu.memory.SetUint16(int(spAddress), value)
	cpu.SetRegister("sp", spAddress-2)

}

func (cpu *CPU) FetachRegisterIndex() int {
	return (int(cpu.Fetch()) % len(cpu.registerNames)) * 2
}

// Execute decodes and executes instructions
func (cpu *CPU) Execute(instruction uint8) {
	switch instruction {

	case constants.MOV_LIT_REG:
		literal := cpu.Fetch16()
		register := cpu.FetachRegisterIndex()
		cpu.registers.SetUint16(int(register), literal)

		//moving register to register
	case constants.MOV_REG_REG:
		registerFrom := cpu.FetachRegisterIndex()
		registerTo := cpu.FetachRegisterIndex()
		value := cpu.registers.GetUint16(int(registerFrom))
		cpu.registers.SetUint16(int(registerTo), value)

	//move register to memory
	case constants.MOV_REG_MEM:
		registerFrom := cpu.FetachRegisterIndex()
		address := cpu.Fetch16()
		value := cpu.registers.GetUint16(int(registerFrom))
		cpu.memory.SetUint16(int(address), value)

	//move memory to register
	case constants.MOV_MEM_REG:
		address := cpu.Fetch16()
		value := cpu.memory.GetUint16(int(address))
		registerTo := cpu.FetachRegisterIndex()
		cpu.registers.SetUint16(int(registerTo), value)

	// add register to register
	case constants.ADD_REG_REG:
		r1 := cpu.Fetch()
		r2 := cpu.Fetch()

		r1Offset := int(r1) * 2
		r2Offset := int(r2) * 2

		registerValue1 := cpu.registers.GetUint16(r1Offset)
		registerValue2 := cpu.registers.GetUint16(r2Offset)

		sum := registerValue1 + registerValue2
		cpu.SetRegister("acc", sum)

		fmt.Printf("ADD_REG_REG: Added r%d (0x%04X) + r%d (0x%04X) = 0x%04X (stored in acc)\n",
			r1, registerValue1, r2, registerValue2, sum)

	case constants.JMP_NOT_EQ:
		value := cpu.Fetch16()
		address := cpu.Fetch16()

		if value != cpu.GetRegister("acc") {
			cpu.SetRegister("ip", address)
		}

	//push literal value on the stack
	case constants.PSH_LIT:
		value := cpu.Fetch16()
		cpu.Push(value)

	//push val from register on the stack
	case constants.PSH_REG:
		registerIndex := cpu.FetachRegisterIndex()
		cpu.Push(cpu.registers.GetUint16(registerIndex))
	}
}

func (cpu *CPU) Step() {
	instruction := cpu.Fetch()
	cpu.Execute(instruction)
}
