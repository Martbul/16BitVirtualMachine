package parser

import (
	"fmt"
)

// ParseInstruction tries all instruction parsers in sequence
func ParseInstruction(input string) (*Node, error) {
	parsers := []Parser{
		// MOV instructions
		{"MovRegToReg", MovRegToReg},
		{"MovLitToReg", MovLitToReg},
		{"MovMemToReg", MovMemToReg},
		{"MovRegToMem", MovRegToMem},
		{"MovLitToMem", MovLitToMem},
		{"MovRegPtrToReg", MovRegPtrToReg},
		{"MovLitOffToReg", MovLitOffToReg},

		// ADD instructions
		{"AddRegToReg", AddRegToReg},
		{"AddLitToReg", AddLitToReg},

		// SUB instructions
		{"SubRegToReg", SubRegToReg},
		{"SubLitToReg", SubLitToReg},

		// MUL instructions
		{"MulRegToReg", MulRegToReg},
		{"MulLitToReg", MulLitToReg},

		// LSF instructions (Left Shift)
		{"LsfRegToReg", LsfRegToReg},
		{"LsfLitToReg", LsfLitToReg},

		// RSF instructions (Right Shift)
		{"RsfRegToReg", RsfRegToReg},
		{"RsfLitToReg", RsfLitToReg},

		// AND instructions
		{"AndRegToReg", AndRegToReg},
		{"AndLitToReg", AndLitToReg},

		// OR instructions
		{"OrRegToReg", OrRegToReg},
		{"OrLitToReg", OrLitToReg},

		// XOR instructions
		{"XorRegToReg", XorRegToReg},
		{"XorLitToReg", XorLitToReg},

		// Single register operations
		{"IncReg", IncReg},
		{"DecReg", DecReg},
		{"NotReg", NotReg},

		// Jump instructions
		{"JeqReg", JeqReg},
		{"JeqLit", JeqLit},
		{"JneReg", JneReg},
		{"JneLit", JneLit},
		{"JltReg", JltReg},
		{"JltLit", JltLit},
		{"JgtReg", JgtReg},
		{"JgtLit", JgtLit},
		{"JleReg", JleReg},
		{"JleLit", JleLit},
		{"JgeReg", JgeReg},
		{"JgeLit", JgeLit},

		// Stack operations
		{"PshLit", PshLit},
		{"PshReg", PshReg},
		{"PopReg", PopReg},

		// Call instructions
		{"CalLit", CalLit},
		{"CalReg", CalReg},

		// No arguments instructions
		{"Ret", Ret},
		{"Hlt", Hlt},
	}

	// Try each parser in sequence
	var errors []string
	for _, parser := range parsers {
		node, err := parser.Fn(input)
		if err == nil {
			return node, nil
		}
		// Collect error for reporting
		errors = append(errors, fmt.Sprintf("%s: %v", parser.Name, err))
	}

	// If we reach here, all parsers failed
	return nil, fmt.Errorf("failed to parse instruction. Errors: %v", errors)
}

// MUL instructions
func RegToReg2(mnemonic, instructionType string) func(string) (*Node, error) {
	return RegToReg(mnemonic, instructionType)
}

// ParseInstructionGeneral provides a choice-based parser similar to the JS version
func ParseInstructionGeneral(input string) (*Node, error) {
	// Try each parser in a way similar to the JS choice mechanism
	// First try MOV
	if node, err := tryParserGroup(input, []Parser{
		{"MovRegToReg", MovRegToReg},
		{"MovLitToReg", MovLitToReg},
		{"MovMemToReg", MovMemToReg},
		{"MovRegToMem", MovRegToMem},
		{"MovLitToMem", MovLitToMem},
		{"MovRegPtrToReg", MovRegPtrToReg},
		{"MovLitOffToReg", MovLitOffToReg},
	}); err == nil {
		return node, nil
	}

	// Then try ADD
	if node, err := tryParserGroup(input, []Parser{
		{"AddRegToReg", AddRegToReg},
		{"AddLitToReg", AddLitToReg},
	}); err == nil {
		return node, nil
	}

	// Then try SUB
	if node, err := tryParserGroup(input, []Parser{
		{"SubRegToReg", SubRegToReg},
		{"SubLitToReg", SubLitToReg},
	}); err == nil {
		return node, nil
	}

	// Then try MUL
	if node, err := tryParserGroup(input, []Parser{
		{"MulRegToReg", MulRegToReg},
		{"MulLitToReg", MulLitToReg},
	}); err == nil {
		return node, nil
	}

	// Continue with other instruction groups
	// ... (similar pattern for other instruction types)

	// If all fails, return an error
	return nil, fmt.Errorf("no parser matched the input: %s", input)
}

// tryParserGroup tries a group of parsers and returns the first successful result
func tryParserGroup(input string, parsers []Parser) (*Node, error) {
	for _, parser := range parsers {
		if node, err := parser.Fn(input); err == nil {
			return node, nil
		}
	}
	return nil, fmt.Errorf("no parser in group matched")
}

// ===================== PARSER FUNCTIONS FOR INSTRUCTIONS ======================== //

// MOV
var MovRegToReg = RegToReg("mov", "MOV_REG_REG")
var MovLitToReg = LitToReg("mov", "MOV_LIT_REG")
var MovMemToReg = MemToReg("mov", "MOV_MEM_REG")
var MovRegToMem = RegToMem("mov", "MOV_REG_MEM")
var MovLitToMem = LitToMem("mov", "MOV_LIT_MEM")
var MovRegPtrToReg = RegPtrToReg("mov", "MOV_REG_PTR_REG")
var MovLitOffToReg = LitOffToReg("mov", "MOV_LIT_OFF_REG")

// ADD
var AddRegToReg = RegToReg("add", "ADD_REG_REG")
var AddLitToReg = LitToReg("add", "ADD_LIT_REG")

// SUB
var SubRegToReg = RegToReg("sub", "SUB_REG_REG")
var SubLitToReg = LitToReg("sub", "SUB_LIT_REG")

// MUL
var MulRegToReg = RegToReg2("mul", "MUL_REG_REG")
var MulLitToReg = LitToReg("mul", "MUL_LIT_REG")

// LSF
var LsfRegToReg = RegToReg("lsf", "LSF_REG_REG")
var LsfLitToReg = LitToReg("lsf", "LSF_LIT_REG")

// RSF
var RsfRegToReg = RegToReg("rsf", "RSF_REG_REG")
var RsfLitToReg = LitToReg("rsf", "RSF_LIT_REG")

// AND
var AndRegToReg = RegToReg("and", "AND_REG_REG")
var AndLitToReg = LitToReg("and", "AND_LIT_REG")

// OR
var OrRegToReg = RegToReg("or", "OR_REG_REG")
var OrLitToReg = LitToReg("or", "OR_LIT_REG")

// XOR
var XorRegToReg = RegToReg("xor", "XOR_REG_REG")
var XorLitToReg = LitToReg("xor", "XOR_LIT_REG")

// SINGLE
var IncReg = SingleReg("inc", "INC_REG")
var DecReg = SingleReg("dec", "DEC_REG")
var NotReg = SingleReg("not", "NOT")

// JEQ
var JeqReg = RegToMem("jeq", "JEQ_REG")
var JeqLit = LitToMem("jeq", "JEQ_LIT")

// JNE
var JneReg = RegToMem("jne", "JNE_REG")
var JneLit = LitToMem("jne", "JMP_NOT_EQ")

// JET
var JltReg = RegToMem("jlt", "JLT_REG")
var JltLit = LitToMem("jlt", "JLT_LIT")

// JGT
var JgtReg = RegToMem("jgt", "JGT_REG")
var JgtLit = LitToMem("jgt", "JGT_LIT")

// JLE
var JleReg = RegToMem("jle", "JLE_REG")
var JleLit = LitToMem("jle", "JLE_LIT")

// JGE
var JgeReg = RegToMem("jge", "JGE_REG")
var JgeLit = LitToMem("jge", "JGE_LIT")

// PUSH/POP
var PshLit = SingleLit("psh", "PSH_LIT")
var PshReg = SingleReg("psh", "PSH_REG")
var PopReg = SingleReg("pop", "POP_REG")

// CALL
var CalLit = SingleLit("cal", "CAL_LIT")
var CalReg = SingleReg("cal", "CAL_REG")

// NO ARGS
var Ret = NoArg("ret", "RET")
var Hlt = NoArg("hlt", "HLT")
