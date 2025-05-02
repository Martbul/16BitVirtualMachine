package parser

import (
	"fmt"
)

// WARN: For constants you can eaither get the value before calling this func or you can make mov,add and so on instr for costant
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
		fmt.Println("sdasdas")
		// Collect error for reporting
		errors = append(errors, fmt.Sprintf("%s: %v", parser.Name, err))
	}

	// all parsers failed
	return nil, fmt.Errorf("failed to parse instruction. Errors: %v", errors)
}

func RegToReg2(mnemonic, instructionType string) func(string) (*Node, error) {
	return RegToReg(mnemonic, instructionType)
}

func ParseInstructionGeneral(input string) (*Node, error) {
	if node, err := ParseConstant(input); err == nil {
		return node, nil
	}
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

	if node, err := tryParserGroup(input, []Parser{
		{"AddRegToReg", AddRegToReg},
		{"AddLitToReg", AddLitToReg},
	}); err == nil {
		return node, nil
	}
	if node, err := tryParserGroup(input, []Parser{
		{"SubRegToReg", SubRegToReg},
		{"SubLitToReg", SubLitToReg},
	}); err == nil {
		return node, nil
	}
	if node, err := tryParserGroup(input, []Parser{
		{"MulRegToReg", MulRegToReg},
		{"MulLitToReg", MulLitToReg},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"LsfRegToReg", LsfRegToReg},
		{"LsfLitToReg", LsfLitToReg},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"RsfRegToReg", RsfRegToReg},
		{"RsfLitToReg", RsfLitToReg},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"AndRegToReg", AndRegToReg},
		{"AndLitToReg", AndLitToReg},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"OrRegToReg", OrRegToReg},
		{"OrLitToReg", OrLitToReg},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"XorRegToReg", XorRegToReg},
		{"XorLitToReg", XorLitToReg},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"IncReg", IncReg},
		{"DecReg", DecReg},
		{"NotReg", NotReg},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"JeqReg", JeqReg},
		{"JeqLit", JeqLit},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"JneReg", JneReg},
		{"JneLit", JneLit},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"JltReg", JltReg},
		{"JltLit", JltLit},
	}); err == nil {
		return node, nil
	}

	if node, err := tryParserGroup(input, []Parser{
		{"JgtReg", JgtReg},
		{"JgtLit", JgtLit},
	}); err == nil {
		return node, nil
	}

	// Try jump less than or equal operations
	if node, err := tryParserGroup(input, []Parser{
		{"JleReg", JleReg},
		{"JleLit", JleLit},
	}); err == nil {
		return node, nil
	}

	// Try jump greater than or equal operations
	if node, err := tryParserGroup(input, []Parser{
		{"JgeReg", JgeReg},
		{"JgeLit", JgeLit},
	}); err == nil {
		return node, nil
	}

	// Try stack operations
	if node, err := tryParserGroup(input, []Parser{
		{"PshLit", PshLit},
		{"PshReg", PshReg},
		{"PopReg", PopReg},
	}); err == nil {
		return node, nil
	}

	// Try call operations
	if node, err := tryParserGroup(input, []Parser{
		{"CalLit", CalLit},
		{"CalReg", CalReg},
	}); err == nil {
		return node, nil
	}

	// Try no-arg operations
	if node, err := tryParserGroup(input, []Parser{
		{"Ret", Ret},
		{"Hlt", Hlt},
	}); err == nil {
		return node, nil
	}

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
