package assembler

import (
	"fmt"
	"os"
	"strings"

	"github.com/yourmodule/parser" // Replace with your actual module path
)

// RegisterMap maps register names to their numeric codes
var RegisterMap = map[string]byte{
	"ip":  0,
	"acc": 1,
	"r1":  2,
	"r2":  3,
	"r3":  4,
	"r4":  5,
	"r5":  6,
	"r6":  7,
	"r7":  8,
	"r8":  9,
	"sp":  10,
	"fp":  11,
}

// InstructionMeta contains metadata about each instruction
type InstructionMeta struct {
	Opcode byte
	Size   int
	Type   string
}

// Define instruction types (similar to the JS instructionTypes object)
const (
	LitReg    = "LIT_REG"
	MemReg    = "MEM_REG"
	RegLit8   = "REG_LIT8"
	RegLit    = "REG_LIT"
	RegMem    = "REG_MEM"
	LitMem    = "LIT_MEM"
	RegReg    = "REG_REG"
	RegPtrReg = "REG_PTR_REG"
	LitOffReg = "LIT_OFF_REG"
	SingleReg = "SINGLE_REG"
	SingleLit = "SINGLE_LIT"
	NoArgs    = "NO_ARGS"
)

// Instructions maps instruction mnemonics to their metadata
var Instructions = map[string]InstructionMeta{
	"MOV_LIT_REG":     {0x10, 4, LitReg},
	"MOV_REG_REG":     {0x11, 3, RegReg},
	"MOV_REG_MEM":     {0x12, 4, RegMem},
	"MOV_MEM_REG":     {0x13, 4, MemReg},
	"MOV_LIT_MEM":     {0x14, 5, LitMem},
	"MOV_REG_PTR_REG": {0x15, 3, RegPtrReg},
	"MOV_LIT_OFF_REG": {0x16, 5, LitOffReg},
	"ADD_REG_REG":     {0x20, 3, RegReg},
	"ADD_LIT_REG":     {0x21, 4, LitReg},
	"SUB_REG_REG":     {0x30, 3, RegReg},
	"SUB_LIT_REG":     {0x31, 4, LitReg},
	"MUL_REG_REG":     {0x40, 3, RegReg},
	"MUL_LIT_REG":     {0x41, 4, LitReg},
	"INC_REG":         {0x50, 2, SingleReg},
	"DEC_REG":         {0x51, 2, SingleReg},
	"LSF_REG_REG":     {0x60, 3, RegReg},
	"LSF_LIT_REG":     {0x61, 4, LitReg},
	"RSF_REG_REG":     {0x70, 3, RegReg},
	"RSF_LIT_REG":     {0x71, 4, LitReg},
	"AND_REG_REG":     {0x80, 3, RegReg},
	"AND_LIT_REG":     {0x81, 4, LitReg},
	"OR_REG_REG":      {0x90, 3, RegReg},
	"OR_LIT_REG":      {0x91, 4, LitReg},
	"XOR_REG_REG":     {0xA0, 3, RegReg},
	"XOR_LIT_REG":     {0xA1, 4, LitReg},
	"NOT":             {0xB0, 2, SingleReg},
	"JEQ_REG":         {0xC0, 4, RegMem},
	"JEQ_LIT":         {0xC1, 5, LitMem},
	"JNE_REG":         {0xD0, 4, RegMem},
	"JMP_NOT_EQ":      {0xD1, 5, LitMem}, // JNE_LIT
	"JLT_REG":         {0xE0, 4, RegMem},
	"JLT_LIT":         {0xE1, 5, LitMem},
	"JGT_REG":         {0xF0, 4, RegMem},
	"JGT_LIT":         {0xF1, 5, LitMem},
	"JLE_REG":         {0xF2, 4, RegMem},
	"JLE_LIT":         {0xF3, 5, LitMem},
	"JGE_REG":         {0xF4, 4, RegMem},
	"JGE_LIT":         {0xF5, 5, LitMem},
	"PSH_LIT":         {0x01, 3, SingleLit},
	"PSH_REG":         {0x02, 2, SingleReg},
	"POP_REG":         {0x03, 2, SingleReg},
	"CAL_LIT":         {0x04, 3, SingleLit},
	"CAL_REG":         {0x05, 2, SingleReg},
	"RET":             {0x06, 1, NoArgs},
	"HLT":             {0x07, 1, NoArgs},
}

func main() {
	// Example program (as in your JS code)
	exampleProgram := []string{
		"start:",
		"  mov $0A, &0050",
		"loop:",
		"  mov &0050, acc",
		"  dec acc",
		"  mov acc, &0050",
		"  inc r2",
		"  inc r2",
		"  inc r2",
		"  jne $00, &[!loop]",
		"end:",
		"  hlt",
	}
	programText := strings.Join(exampleProgram, "\n")

	// Parse the program
	parsedNodes, err := parser.ParseProgram(programText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing program: %v\n", err)
		os.Exit(1)
	}

	// Initialize machine code array and labels map
	machineCode := []byte{}
	labels := make(map[string]int)
	currentAddress := 0

	// First pass: resolve labels
	for _, node := range parsedNodes {
		if node.Type == "LABEL" {
			labelName, ok := node.Value.(map[string]interface{})["name"].(string)
			if !ok {
				fmt.Fprintf(os.Stderr, "Invalid label format\n")
				os.Exit(1)
			}
			labels[labelName] = currentAddress
		} else {
			// Must be an instruction
			instrValue, ok := node.Value.(map[string]interface{})
			if !ok {
				fmt.Fprintf(os.Stderr, "Invalid instruction format\n")
				os.Exit(1)
			}
			instrType := instrValue["instruction"].(string)
			metadata, exists := Instructions[instrType]
			if !exists {
				fmt.Fprintf(os.Stderr, "Unknown instruction: %s\n", instrType)
				os.Exit(1)
			}
			currentAddress += metadata.Size
		}
	}

	// Second pass: encode instructions
	for _, node := range parsedNodes {
		// Skip labels in second pass
		if node.Type == "LABEL" {
			continue
		}

		instrValue := node.Value.(map[string]interface{})
		instrType := instrValue["instruction"].(string)
		metadata := Instructions[instrType]

		// Add opcode
		machineCode = append(machineCode, metadata.Opcode)

		// Get arguments
		args := instrValue["args"].([]*parser.Node)

		// Encode arguments based on instruction type
		switch metadata.Type {
		case LitReg, MemReg:
			encodeLitOrMem(&machineCode, args[0], labels)
			encodeReg(&machineCode, args[1])

		case RegLit8:
			encodeReg(&machineCode, args[0])
			encodeLit8(&machineCode, args[1], labels)

		case RegLit, RegMem:
			encodeReg(&machineCode, args[0])
			encodeLitOrMem(&machineCode, args[1], labels)

		case LitMem:
			encodeLitOrMem(&machineCode, args[0], labels)
			encodeLitOrMem(&machineCode, args[1], labels)

		case RegReg, RegPtrReg:
			encodeReg(&machineCode, args[0])
			encodeReg(&machineCode, args[1])

		case LitOffReg:
			encodeLitOrMem(&machineCode, args[0], labels)
			encodeReg(&machineCode, args[1])
			encodeReg(&machineCode, args[2])

		case SingleReg:
			encodeReg(&machineCode, args[0])

		case SingleLit:
			encodeLitOrMem(&machineCode, args[0], labels)
		}
	}

	// Print the resulting machine code
	fmt.Println("Machine code:")
	for i, b := range machineCode {
		fmt.Printf("%02X ", b)
		if (i+1)%8 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()

	// Print resolved labels
	fmt.Println("Labels:")
	for label, addr := range labels {
		fmt.Printf("%s: 0x%04X\n", label, addr)
	}
}

// encodeLitOrMem encodes a literal or memory address
func encodeLitOrMem(machineCode *[]byte, node *parser.Node, labels map[string]int) {
	var hexVal int

	// Check if this is a label reference (VARIABLE in JS code)
	if node.Type == "VARIABLE" {
		labelName := node.Value.(string)
		addr, exists := labels[labelName]
		if !exists {
			fmt.Fprintf(os.Stderr, "Error: label '%s' wasn't resolved\n", labelName)
			os.Exit(1)
		}
		hexVal = addr
	} else {
		// Must be a literal (e.g., "$0A", "&0050")
		// Extract value from hex format and parse
		valueStr := node.Value.(string)
		// Remove "$" or "&" prefix and parse
		if strings.HasPrefix(valueStr, "$") || strings.HasPrefix(valueStr, "&") {
			valueStr = valueStr[1:]
		}
		fmt.Sscanf(valueStr, "%X", &hexVal)
	}

	// Push high byte and low byte
	highByte := byte((hexVal & 0xFF00) >> 8)
	lowByte := byte(hexVal & 0x00FF)
	*machineCode = append(*machineCode, highByte, lowByte)
}

// encodeLit8 encodes an 8-bit literal
func encodeLit8(machineCode *[]byte, node *parser.Node, labels map[string]int) {
	var hexVal int

	if node.Type == "VARIABLE" {
		labelName := node.Value.(string)
		addr, exists := labels[labelName]
		if !exists {
			fmt.Fprintf(os.Stderr, "Error: label '%s' wasn't resolved\n", labelName)
			os.Exit(1)
		}
		hexVal = addr
	} else {
		valueStr := node.Value.(string)
		if strings.HasPrefix(valueStr, "$") || strings.HasPrefix(valueStr, "&") {
			valueStr = valueStr[1:]
		}
		fmt.Sscanf(valueStr, "%X", &hexVal)
	}

	// Push just the low byte
	lowByte := byte(hexVal & 0x00FF)
	*machineCode = append(*machineCode, lowByte)
}

// encodeReg encodes a register reference
func encodeReg(machineCode *[]byte, node *parser.Node) {
	regName := strings.ToLower(node.Value.(string))
	regCode, exists := RegisterMap[regName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: unknown register '%s'\n", regName)
		os.Exit(1)
	}
	*machineCode = append(*machineCode, regCode)
}
