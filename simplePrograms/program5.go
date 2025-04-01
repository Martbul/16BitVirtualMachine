package simpleprograms

import (
	"github.com/martbul/constants"
	cpuPack "github.com/martbul/cpu"
	"github.com/martbul/devices"
	"github.com/martbul/memory"
	memMapper "github.com/martbul/memoryMapper"
)

//INFO: The program:

func Program5() {

	//creating the memory, with buffer(a byte slice that is with capacity of 256*256 bytes)
	memory := memory.CreateMemory(256 * 256)
	//getting the memmory's byte slice
	memoryBytes := memory.GetBuffer()

	memoryMapper := memMapper.NewMemoryMapper()
	memoryMapper.Map(memory, 0, 0xffff) //INFO: the range from 0 to 0xffff is maped to be RAM

	//Map 0xFF bytes of the address space to an "output device" - standart stdout
	memoryMapper.Map(devices.CreateScreenDevice(), 0x3000, 0x30ff, true) //INFO: Writes to this range will be displayed as output.

	cpu := cpuPack.NewCPU(memoryMapper) //WARN: changed the type that newcpu receives, possible error in the future
	ip := 0

	str := "hello world"
	for index, char := range str {
		writeCharToScreen(&memoryBytes, &ip, char, index)
	}

	// Writing the halt instruction to stop the program
	memoryBytes[ip] = constants.HLT
	ip++

	cpu.Run() //INFO: Starts executing the instructions stored in memory.:

}

// Function to write a character to the screen memory
func writeCharToScreen(memoryBytes *[]byte, ip *int, char rune, position int) {
	// Writing MOV_LIT_REG instruction (load immediate value into register)
	(*memoryBytes)[*ip] = constants.MOV_LIT_REG
	*ip++

	// Writing 0x00 as high byte (it's part of the instruction)
	(*memoryBytes)[*ip] = 0x00
	*ip++

	// Writing the character 'P' (ASCII value 80) to memory (lower byte)
	(*memoryBytes)[*ip] = byte(char)
	*ip++

	// Writing to register 1
	(*memoryBytes)[*ip] = constants.R1
	*ip++

	// Writing MOV_REG_MEM instruction (move value from register to memory)
	(*memoryBytes)[*ip] = constants.MOV_REG_MEM
	*ip++

	// Writing register R1 to memory address 0x3000 (mapped to screen)
	(*memoryBytes)[*ip] = constants.R1
	*ip++

	// Writing memory address 0x3000 (the screen device's address)
	(*memoryBytes)[*ip] = 0x30
	*ip++
	(*memoryBytes)[*ip] = byte(position)
	*ip++
}
