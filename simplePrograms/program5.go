package simpleprograms

import (
	"bufio"
	"os"

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
	memoryMapper.Map(memory, 0, 0xffff) // memory is the device, 0 is the start and 0xffff is the end

	//Map 0xFF bytes of the address space to an "output device" - standart stdout
	memoryMapper.Map(devices.CreateScreenDevice(), 0x3000, 0x30ff, true)

	cpu := cpuPack.NewCPU(memoryMapper) //WARN: changed the type that newcpu receives, possible error in the future
	i := 0
	
	memoryBytes[i] = constants.MOV_LIT_REG
	i++
	memoryBytes[i] = 0x00 //high byte
	i++
	memoryBytes[i] = 0x
}
