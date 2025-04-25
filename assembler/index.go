package assembler

import (
	"fmt"
	"os"
	"strings"

	"github.com/martbul/assembler/parser"
	"github.com/martbul/instructions"
	"github.com/martbul/registers"
)

// RegisterMap maps register names to their numeric codes
//v////ar RegisterMap = map[string]byte{
//	"ip":  0,
//	"acc": 1,
//	"r1":  2,
///	"r2":  3,
//	"r3":  4,
//	"r4":  5,
///	"r5":  6,
//	"r6":  7,
//	"r7":  8,
//	"r8":  9,
//	"sp":  10,
//	"fp":  11,
//}

func AssembleProgram(program []string) {
	programText := strings.Join(program, "\n")

	// Parse the program
	parsedNodes, err := parser.ParseProgram(programText)
	for _, n := range parsedNodes {

		parser.PrettyPrintNode(n)
	}

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
			labelName, ok := node.Value.(map[string]interface{})["label"].(string)
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
			metadata, exists := instructions.GetInstructionByName(instrType)
			if !exists {
				fmt.Fprintf(os.Stderr, "Unknown instruction: %s\n", instrType)
				os.Exit(1)
			}
			currentAddress += int(metadata.Size)
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
		metadata, _ := instructions.GetInstructionByName(instrType)
		// Add opcode
		machineCode = append(machineCode, metadata.Opcode)

		// Get arguments
		args := instrValue["args"].([]*parser.Node)

		// Encode arguments based on instruction type
		switch metadata.Type {
		case instructions.LitReg, instructions.MemReg:
			encodeLitOrMem(&machineCode, args[0], labels)
			encodeReg(&machineCode, args[1])

		case instructions.RegLit8:
			encodeReg(&machineCode, args[0])
			encodeLit8(&machineCode, args[1], labels)

		case instructions.RegLit, instructions.RegMem:
			encodeReg(&machineCode, args[0])
			encodeLitOrMem(&machineCode, args[1], labels)

		case instructions.LitMem:
			encodeLitOrMem(&machineCode, args[0], labels)
			encodeLitOrMem(&machineCode, args[1], labels)

		case instructions.RegReg, instructions.RegPtrReg:
			encodeReg(&machineCode, args[0])
			encodeReg(&machineCode, args[1])

		case instructions.LitOffReg:
			encodeLitOrMem(&machineCode, args[0], labels)
			encodeReg(&machineCode, args[1])
			encodeReg(&machineCode, args[2])

		case instructions.SingleReg:
			encodeReg(&machineCode, args[0])

		case instructions.SingleLit:
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
	regCode, exists := registers.Map[regName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: unknown register '%s'\n", regName)
		os.Exit(1)
	}
	*machineCode = append(*machineCode, byte(regCode))
}

// encodeLitOrMem encodes a literal or memory address
func encodeLitOrMem(machineCode *[]byte, node *parser.Node, labels map[string]int) {
	var hexVal int

	// Handle MEMORY_REFERENCE type which contains a nested VARIABLE
	if node.Type == "MEMORY_REFERENCE" {
		if nestedNode, ok := node.Value.(map[string]interface{}); ok {
			if nestedType, ok := nestedNode["type"].(string); ok && nestedType == "VARIABLE" {
				if labelName, ok := nestedNode["value"].(string); ok {
					addr, exists := labels[labelName]
					if !exists {
						fmt.Fprintf(os.Stderr, "Error: label '%s' wasn't resolved\n", labelName)
						os.Exit(1)
					}
					hexVal = addr
					goto encodeValue // Skip further processing
				}
			}
		}
	}

	// Check if this is a direct VARIABLE (label reference)
	if node.Type == "VARIABLE" {
		labelName := node.Value.(string)
		addr, exists := labels[labelName]
		if !exists {
			fmt.Fprintf(os.Stderr, "Error: label '%s' wasn't resolved\n", labelName)
			os.Exit(1)
		}
		hexVal = addr
	} else if node.Type == "HEX_LITERAL" || node.Type == "ADDRESS" {
		// Must be a literal (e.g., "$0A", "&0050")
		// Extract value from hex format and parse
		valueStr := node.Value.(string)
		// Remove "$" or "&" prefix if present
		if strings.HasPrefix(valueStr, "$") || strings.HasPrefix(valueStr, "&") {
			valueStr = valueStr[1:]
		}
		fmt.Sscanf(valueStr, "%X", &hexVal)
	}

encodeValue:
	// Push high byte and low byte
	highByte := byte((hexVal & 0xFF00) >> 8)
	lowByte := byte(hexVal & 0x00FF)
	*machineCode = append(*machineCode, highByte, lowByte)
}
