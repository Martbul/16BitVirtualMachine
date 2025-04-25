package simpleprograms

import (
	cpuPack "github.com/martbul/cpu"
	"github.com/martbul/devices"
	"github.com/martbul/instructions"
	"github.com/martbul/memory"
	memMapper "github.com/martbul/memoryMapper"
)

func Program8() {

	//creating the memory, with buffer(a byte slice that is with capacity of 256*256 bytes)
	memory := memory.CreateMemory(256 * 256)
	//getting the memmory's byte slice
	memoryBytes := memory.GetBuffer()

	memoryMapper := memMapper.NewMemoryMapper()

}

func createBankedMemory(n, bankcedSIze)
