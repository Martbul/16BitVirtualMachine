package main

import (
	"github.com/martbul/constants"
	cpuPack "github.com/martbul/cpu"
	"github.com/martbul/memory"
)

func main() {

	memory := memory.CreateMemory(256)

	cpu := cpuPack.NewCPU(memory)

	writableBytes := memory.GetBuffer()

	writableBytes[0] = constants.MOV_LIT_R1
	writableBytes[1] = 0x12 // 0x1234
	writableBytes[2] = 0x34

	writableBytes[3] = constants.MOV_LIT_R2
	writableBytes[4] = 0xAB // 0xABCD
	writableBytes[5] = 0xCD

	writableBytes[6] = constants.ADD_REG_REG
	writableBytes[7] = 2 //r1 index
	writableBytes[8] = 3 //r2 index

	cpu.Debug()
	cpu.Step()
	cpu.Debug()
	cpu.Step()
	cpu.Debug()
	cpu.Step()
	cpu.Debug()
}
