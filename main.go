package main

import (
	"fmt"

	"github.com/martbul/constants"
	cpuPack "github.com/martbul/cpu"
	"github.com/martbul/memory"
)

const (
	IP  = 0
	ACC = 1
	R1  = 2
	R2  = 3
)

func main() {

	memory := memory.CreateMemory(256 * 256)
	writableBytes := memory.GetBuffer()

	cpu := cpuPack.NewCPU(memory)

	i := 0

	writableBytes[i] = constants.MOV_LIT_REG
	i++
	writableBytes[i] = 0x12
	i++
	writableBytes[i] = 0x34
	i++
	writableBytes[i] = R1
	i++

	writableBytes[i] = constants.MOV_LIT_REG
	i++
	writableBytes[i] = 0xAB
	i++
	writableBytes[i] = 0xCD
	i++
	writableBytes[i] = R2
	i++

	writableBytes[i] = constants.ADD_REG_REG
	i++
	writableBytes[i] = R1
	i++
	writableBytes[i] = R2
	i++

	writableBytes[i] = constants.MOV_REG_MEM
	i++
	writableBytes[i] = ACC
	i++
	writableBytes[i] = 0x01
	i++
	writableBytes[i] = 0x00
	i++

	// ✅ Debug: Print the written instructions before execution
	fmt.Println("Instruction Memory:", writableBytes[:i])

	cpu.Debug()
	cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	cpu.ViewMemoryAt(0x0100)

	cpu.Step()
	cpu.Debug()
	cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	cpu.ViewMemoryAt(0x0100)

	cpu.Step()
	cpu.Debug()
	cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	cpu.ViewMemoryAt(0x0100)

	cpu.Step()
	cpu.Debug()
	cpu.ViewMemoryAt(int(cpu.GetRegister("ip")))
	cpu.ViewMemoryAt(0x0100)

}
