package simpleprograms

import (
	"fmt"

	"github.com/martbul/constants"
	cpuPack "github.com/martbul/cpu"
	"github.com/martbul/memory"
)

const (
	IP  = 0 //Instruction Pointer
	ACC = 1 // Accumulator
	R1  = 2 // General register 1
	R2  = 3 // General Register 2
)

func Program1() {
	//creating the memory, with buffer(a byte slice that is with capacity of 256*256 bytes)
	memory := memory.CreateMemory(256 * 256)
	//getting the memmory's byte slice
	memoryBytes := memory.GetBuffer()

	cpu := cpuPack.NewCPU(memory)

	i := 0

	// Move Literal 0x1234 → R1
	memoryBytes[i] = constants.MOV_LIT_REG
	i++
	memoryBytes[i] = 0x12 // High byte of 0x1234
	i++
	memoryBytes[i] = 0x34 // Low byte of 0x1234
	i++
	memoryBytes[i] = R1
	i++

	// Move Literal 0xABCD → R2
	memoryBytes[i] = constants.MOV_LIT_REG
	i++
	memoryBytes[i] = 0xAB
	i++
	memoryBytes[i] = 0xCD
	i++
	memoryBytes[i] = R2
	i++

	// Add R1 + R2 → ACC
	memoryBytes[i] = constants.ADD_REG_REG
	i++
	memoryBytes[i] = R1
	i++
	memoryBytes[i] = R2
	i++

	// Move ACC → Memory at address 0x0100
	memoryBytes[i] = constants.MOV_REG_MEM
	i++
	memoryBytes[i] = ACC
	i++
	memoryBytes[i] = 0x01
	i++
	memoryBytes[i] = 0x00
	i++

	// ✅ Debug: Print the written instructions before execution
	fmt.Println("Instruction Memory:", memoryBytes[:i])

	// Debug initial state
	cpu.Debug()
	cpu.ViewMemoryAt(int(cpu.GetRegister("ip"))) // View memory at IP
	cpu.ViewMemoryAt(0x0100)                     // View memory at 0x0100

	// Step through the instructions and print state
	for step := 1; step <= 4; step++ {
		fmt.Printf("\n🔹 Step %d\n", step)
		cpu.Step()
		cpu.Debug()
		cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
		cpu.ViewMemoryAt(0x0100)
	}
}
