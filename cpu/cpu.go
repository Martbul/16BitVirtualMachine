// 16 bit cpu
// each register is 16 bits wide
package cpu

import (
	"fmt"
	"github.com/martbul/memory"
	"log"
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

// Debug prints the register values in hexadecimal format
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

// Fetch fetches the next instruction from memory
func (cpu *CPU) Fetch() uint8 {
	// Get the current instruction pointer (ip)
	nextInstructionAddress := cpu.GetRegister("ip")

	// Fetch the instruction (byte) from memory at the current address
	instruction := cpu.memory.GetUint8(int(nextInstructionAddress))

	// Increment the ip register to the next instruction address
	cpu.SetRegister("ip", nextInstructionAddress+1)

	// Return the instruction
	return instruction
}
