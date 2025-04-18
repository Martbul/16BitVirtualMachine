package simpleprograms

import (
	"github.com/martbul/constants"
	"github.com/martbul/instructions"
	"github.com/martbul/memory"
)

//INFO: The program:

// mov $5151, r1
// mov $4242, r2
// psh r1
// psh r2
// pop r1
// pop r2

func Program3() {
	//creating the memory, with buffer(a byte slice that is with capacity of 256*256 bytes)
	memory := memory.CreateMemory(256 * 256)
	//getting the memmory's byte slice
	memoryBytes := memory.GetBuffer()

	//	cpu := cpuPack.NewCPU(memory)

	i := 0

	memoryBytes[i] = instructions.MOV_LIT_REG
	i++
	memoryBytes[i] = 0x51
	i++
	memoryBytes[i] = 0x51
	i++
	memoryBytes[i] = constants.R1
	i++

	memoryBytes[i] = instructions.MOV_LIT_REG
	i++
	memoryBytes[i] = 0x42
	i++
	memoryBytes[i] = 0x42
	i++
	memoryBytes[i] = constants.R2
	i++

	memoryBytes[i] = instructions.PSH_REG
	i++
	memoryBytes[i] = constants.R1
	i++

	memoryBytes[i] = instructions.PSH_REG
	i++
	memoryBytes[i] = constants.R2
	i++

	memoryBytes[i] = instructions.POP
	i++
	memoryBytes[i] = constants.R1
	i++

	memoryBytes[i] = instructions.POP
	i++
	memoryBytes[i] = constants.R2
	i++

	//	cpu.Debug()
	//	cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	///	cpu.ViewMemoryAt(0xffff - 1 - 6) //INFO: The start of the stack

	// scanner := bufio.NewScanner(os.Stdin)
	//
	//	for scanner.Scan() {
	//		cpu.Step()
	//		cpu.Debug()
	//		cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	//		cpu.ViewMemoryAt(0xffff - 1 - 6)
	//	}
}
