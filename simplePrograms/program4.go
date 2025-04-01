package simpleprograms

import (
	"github.com/martbul/constants"
	"github.com/martbul/memory"
)

//INFO: The program:

// psh 0x3333
// psh 0x2222
// psh 0x1111
// mov 0x1234, r1
// mov 0x5678, r4
// psh 0x0000
// cal my_subroutine
// psh 0x4444
// ;; at address 0x3000
// my_subroutine:
//   psh 0x0102
//   psh 0x0304
//   psh 0x0506
//   mov 0x0708, r1
//   mov 0x090A, r8
//   ret

func Program4() {
	//creating the memory, with buffer(a byte slice that is with capacity of 256*256 bytes)
	memory := memory.CreateMemory(256 * 256)
	//getting the memmory's byte slice
	memoryBytes := memory.GetBuffer()

	//	cpu := cpuPack.NewCPU(memory)

	subroutineAddress := 0x3000

	i := 0

	memoryBytes[i] = constants.PSH_LIT
	i++
	memoryBytes[i] = 0x33
	i++
	memoryBytes[i] = 0x33
	i++

	memoryBytes[i] = constants.PSH_LIT
	i++
	memoryBytes[i] = 0x22
	i++
	memoryBytes[i] = 0x22
	i++

	memoryBytes[i] = constants.PSH_LIT
	i++
	memoryBytes[i] = 0x11
	i++
	memoryBytes[i] = 0x11
	i++

	memoryBytes[i] = constants.MOV_LIT_REG
	i++
	memoryBytes[i] = 0x12
	i++
	memoryBytes[i] = 0x34
	i++
	memoryBytes[i] = constants.R1
	i++

	memoryBytes[i] = constants.MOV_LIT_REG
	i++
	memoryBytes[i] = 0x56
	i++
	memoryBytes[i] = 0x78
	i++
	memoryBytes[i] = constants.R4
	i++

	memoryBytes[i] = constants.PSH_LIT
	i++
	memoryBytes[i] = 0x00
	i++
	memoryBytes[i] = 0x00
	i++

	// Call subroutine (CAL_LIT)
	memoryBytes[i] = constants.CAL_LIT
	i++
	memoryBytes[i] = byte((subroutineAddress & 0xFF00) >> 8) // High byte
	i++
	memoryBytes[i] = byte(subroutineAddress & 0x00FF) // Low byte
	i++

	// Push literal values (PSH_LIT)
	memoryBytes[i] = constants.PSH_LIT
	i++
	memoryBytes[i] = 0x44
	i++
	memoryBytes[i] = 0x44
	i++

	// Move to subroutine
	i = subroutineAddress

	// Push more literals in subroutine
	memoryBytes[i] = constants.PSH_LIT
	i++
	memoryBytes[i] = 0x01
	i++
	memoryBytes[i] = 0x02
	i++

	memoryBytes[i] = constants.PSH_LIT
	i++
	memoryBytes[i] = 0x03
	i++
	memoryBytes[i] = 0x04
	i++

	memoryBytes[i] = constants.PSH_LIT
	i++
	memoryBytes[i] = 0x05
	i++
	memoryBytes[i] = 0x06
	i++

	// Move literal to register (MOV_LIT_REG)
	memoryBytes[i] = constants.MOV_LIT_REG
	i++
	memoryBytes[i] = 0x07
	i++
	memoryBytes[i] = 0x08
	i++
	memoryBytes[i] = constants.R1
	i++

	memoryBytes[i] = constants.MOV_LIT_REG
	i++
	memoryBytes[i] = 0x09
	i++
	memoryBytes[i] = 0x0A
	i++
	memoryBytes[i] = constants.R8
	i++

	// Return (RET)
	memoryBytes[i] = constants.RET
	i++

	//	cpu.Debug()
	//	cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	//	cpu.ViewMemoryAt(0xffff-1-42, 44) //INFO: The start of the stack

	// scanner := bufio.NewScanner(os.Stdin)
	//
	//	for scanner.Scan() {
	//		cpu.Step()
	//		cpu.Debug()
	//		cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	//		cpu.ViewMemoryAt(0xffff-1-42, 44)
	//	}
}
