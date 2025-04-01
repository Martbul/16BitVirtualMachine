package simpleprograms

import (
	"github.com/martbul/constants"
	"github.com/martbul/memory"
)

//INFO: THE PROGRAM

// start:
// mov #0x0100, r1
// mov 0x0001, r2
// add r1, r2
// mov acc, #0x0100
// jne 0x0003, start

func Program2() {
	//creating the memory, with buffer(a byte slice that is with capacity of 256*256 bytes)
	memory := memory.CreateMemory(256 * 256)
	//getting the memmory's byte slice
	memoryBytes := memory.GetBuffer()

	//	cpu := cpuPack.NewCPU(memory)

	i := 0

	memoryBytes[i] = constants.MOV_MEM_REG
	i++
	memoryBytes[i] = 0x01
	i++
	memoryBytes[i] = 0x00
	i++
	memoryBytes[i] = constants.R1
	i++

	memoryBytes[i] = constants.MOV_LIT_REG
	i++
	memoryBytes[i] = 0x00
	i++
	memoryBytes[i] = 0x01
	i++
	memoryBytes[i] = constants.R2
	i++

	memoryBytes[i] = constants.ADD_REG_REG
	i++
	memoryBytes[i] = constants.R1
	i++
	memoryBytes[i] = constants.R2
	i++

	memoryBytes[i] = constants.MOV_REG_MEM
	i++
	memoryBytes[i] = constants.ACC
	i++
	memoryBytes[i] = 0x01
	i++
	memoryBytes[i] = 0x00
	i++

	memoryBytes[i] = constants.JMP_NOT_EQ
	i++
	memoryBytes[i] = 0x00
	i++
	memoryBytes[i] = 0x03
	i++
	memoryBytes[i] = 0x00
	i++
	memoryBytes[i] = 0x00
	i++

	//	cpu.Debug()
	//	cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	//	cpu.ViewMemoryAt(0x0100)

	// scanner := bufio.NewScanner(os.Stdin)
	//
	//	for scanner.Scan() {
	//		cpu.Step()
	//		cpu.Debug()
	//		cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	//		cpu.ViewMemoryAt(0x0100)
	//	}
}
