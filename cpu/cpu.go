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
	}

	// Create registers memory space with 2 bytes per register
	registers := memory.NewDataView(len(registerNames) * 2)
	registerMap := make(map[string]int)
	for i, name := range registerNames {
		registerMap[name] = i * 2
	}

	return &CPU{
		memory:        mem,
		registerNames: registerNames,
		registers:     registers,
		registerMap:   registerMap,
	}
}

func (cpu *CPU) Debug() {
	for _, name := range cpu.registerNames {
		value := cpu.GetRegister(name)
		fmt.Printf("%s: 0x%04X\n", name, value)
	}
	fmt.Println()
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

// Execute decodes and executes instructions
func (cpu *CPU) Execute(instruction uint8) {
	switch instruction {

	case constants.MOV_LIT_R1:
		literal := cpu.Fetch16()
		cpu.SetRegister("r1", literal)
		fmt.Printf("MOV_LIT_R1: Set r1 = 0x%04X\n", literal)

	case constants.MOV_LIT_R2:
		literal := cpu.Fetch16()
		cpu.SetRegister("r2", literal)
		fmt.Printf("MOV_LIT_R2: Set r2 = 0x%04X\n", literal)

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

	default:
		log.Fatalf("execute: Unknown instruction 0x%02X", instruction)
	}
}

func (cpu *CPU) Step() {
	instruction := cpu.Fetch()
	cpu.Execute(instruction)
}
