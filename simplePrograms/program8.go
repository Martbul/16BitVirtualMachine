package simpleprograms

import (
	"bytes"
	"fmt"

	cpuPack "github.com/martbul/cpu"
	"github.com/martbul/memory"
	memMapper "github.com/martbul/memoryMapper"
)

var dataViewMethods = []string{"getUint8", "getUint16", "setUint8", "setUint16"}

func Program8() {

	memoryMapper := memMapper.NewMemoryMapper()
	cpu := cpuPack.NewCPU(memoryMapper)
	bankSize := 0xff
	var nBanks uint16 = 8
	memoryBankDevice := createBankedMemory(nBanks, bankSize, cpu)
	memoryMapper.Map(memoryBankDevice, 0, bankSize)
	regularMemory := memory.CreateMemory(0xff00)
	memoryMapper.Map(regularMemory, bankSize, 0xffff, true)

	fmt.Println("writing value 1 at addr 0")
	memoryMapper.SetUint16(0, 1)

	v, _ := memoryMapper.GetUint16(0)
	fmt.Println("reading value at addr 0: ", v)

	fmt.Println(":::switching memory bank 0 -> 1")
	cpu.SetRegister("mb", 1)

	v, _ = memoryMapper.GetUint16(0)
	fmt.Println("reading value at addr 0: ", v)

	fmt.Println("writing value 42 at addr 0")
	memoryMapper.SetUint16(0, 42)

	fmt.Println(":::switching memory bank 1 -> 2")
	cpu.SetRegister("mb", 2)

	v, _ = memoryMapper.GetUint16(0)
	fmt.Println("reading value at addr 0: ", v)

	fmt.Println(":::switching memory bank 2 -> 1")
	cpu.SetRegister("mb", 1)

	v, _ = memoryMapper.GetUint16(0)
	fmt.Println("reading value at addr 0: ", v)

	fmt.Println(":::switching memory bank 1 -> 0")
	cpu.SetRegister("mb", 0)

	v, _ = memoryMapper.GetUint16(0)
	fmt.Println("reading value at addr 0: ", v)

}

func createBankedMemory(n uint16, bankSize int, cpu *cpuPack.CPU) memMapper.MemoryDevice {
	// Create bank buffers
	bankBuffers := make([][]byte, n)
	for i := range bankBuffers {
		bankBuffers[i] = make([]byte, bankSize)
	}

	// Create views for each bank buffer
	banks := make([]*bytes.Buffer, n)
	for i, buffer := range bankBuffers {
		banks[i] = bytes.NewBuffer(buffer)
	}

	// Create and return a struct that satisfies the MemoryDevice interface
	return &BankedMemory{
		banks:    banks,
		n:        n,
		cpu:      cpu,
		bankSize: bankSize,
	}
}

// BankedMemory implements MemoryDevice interface
type BankedMemory struct {
	banks    []*bytes.Buffer
	n        uint16
	cpu      *cpuPack.CPU
	bankSize int
}

// GetUint8 reads a byte from the currently active memory bank
func (bm *BankedMemory) GetUint8(address int) uint8 {
	bankIndex := bm.cpu.GetRegister("mb") % bm.n
	memoryBankToUse := bm.banks[bankIndex]
	bytes := memoryBankToUse.Bytes()

	if address >= len(bytes) {
		panic(fmt.Sprintf("Memory access out of bounds: %d", address))
	}

	return bytes[address]
}

// SetUint8 writes a byte to the currently active memory bank
func (bm *BankedMemory) SetUint8(address int, value uint8) {
	bankIndex := bm.cpu.GetRegister("mb") % bm.n
	memoryBankToUse := bm.banks[bankIndex]
	bytes := memoryBankToUse.Bytes()

	if address >= len(bytes) {
		panic(fmt.Sprintf("Memory access out of bounds: %d", address))
	}

	bytes[address] = value
}

// GetUint16 reads a 16-bit value from the currently active memory bank
func (bm *BankedMemory) GetUint16(address int) uint16 {
	bankIndex := bm.cpu.GetRegister("mb") % bm.n
	memoryBankToUse := bm.banks[bankIndex]
	bytes := memoryBankToUse.Bytes()

	if address+1 >= len(bytes) {
		panic(fmt.Sprintf("Memory access out of bounds: %d", address))
	}

	// Assuming big-endian by default
	return (uint16(bytes[address]) << 8) | uint16(bytes[address+1])
}

// SetUint16 writes a 16-bit value to the currently active memory bank
func (bm *BankedMemory) SetUint16(address int, value uint16) {
	bankIndex := bm.cpu.GetRegister("mb") % bm.n
	memoryBankToUse := bm.banks[bankIndex]
	bytes := memoryBankToUse.Bytes()

	if address+1 >= len(bytes) {
		panic(fmt.Sprintf("Memory access out of bounds: %d", address))
	}

	// Assuming big-endian by default
	bytes[address] = byte((value >> 8) & 0xFF)
	bytes[address+1] = byte(value & 0xFF)
}
