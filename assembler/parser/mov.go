package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
)

type Parser struct {
	Name string
	Fn   func(string) (*Node, error)
}

func LitToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[LitRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			fmt.Printf("Error building parser for %s: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			fmt.Printf("Error parsing with %s: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		fmt.Printf("Parsed instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

var MovLitToReg = LitToReg("mov", "MOV_LIT_REG")
var AddLitToReg = LitToReg("add", "ADD_LIT_REG")
var SubLitToReg = LitToReg("sub", "SUB_LIT_REG")
var AndLitToReg = LitToReg("and", "AND_LIT_REG")
var OrLitToReg = LitToReg("or", "OR_LIT_REG")
var XorLitToReg = LitToReg("xor", "XOR_LIT_REG")

func RegToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[RegRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			fmt.Printf("Error building parser for %s reg-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			fmt.Printf("Error parsing with %s reg-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		fmt.Printf("Parsed reg-to-reg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

var MovRegToReg = RegToReg("mov", "MOV_REG_REG")
var AddRegToReg = RegToReg("add", "ADD_REG_REG")
var SubRegToReg = RegToReg("sub", "SUB_REG_REG")
var AndRegToReg = RegToReg("and", "AND_REG_REG")
var OrRegToReg = RegToReg("or", "OR_REG_REG")
var XorRegToReg = RegToReg("xor", "XOR_REG_REG")

func RegToMem(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[RegMemInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			fmt.Printf("Error building parser for %s reg-to-mem: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			fmt.Printf("Error parsing with %s reg-to-mem: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		fmt.Printf("Parsed reg-to-mem instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

var MovRegToMem = RegToMem("mov", "MOV_REG_MEM")
var AddRegToMem = RegToMem("add", "ADD_REG_MEM")
var SubRegToMem = RegToMem("sub", "SUB_REG_MEM")
var AndRegToMem = RegToMem("and", "AND_REG_MEM")
var OrRegToMem = RegToMem("or", "OR_REG_MEM")
var XorRegToMem = RegToMem("xor", "XOR_REG_MEM")

func MemToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[MemRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			fmt.Printf("Error building parser for %s mem-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			fmt.Printf("Error parsing with %s mem-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		fmt.Printf("Parsed mem-to-reg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

var MovMemToReg = MemToReg("mov", "MOV_MEM_REG")
var AddMemToReg = MemToReg("add", "ADD_MEM_REG")
var SubMemToReg = MemToReg("sub", "SUB_MEM_REG")
var AndMemToReg = MemToReg("and", "AND_MEM_REG")
var OrMemToReg = MemToReg("or", "OR_MEM_REG")
var XorMemToReg = MemToReg("xor", "XOR_MEM_REG")

func LitToMem(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[LitMemInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			fmt.Printf("Error building parser for %s lit-to-mem: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			fmt.Printf("Error parsing with %s lit-to-mem: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		fmt.Printf("Parsed lit-to-mem instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

var MovLitToMem = LitToMem("mov", "MOV_LIT_MEM")
var AddLitToMem = LitToMem("add", "ADD_LIT_MEM")
var SumLitToMem = LitToMem("sub", "SUM_LIT_MEM")
var AndLitToMem = LitToMem("and", "AND_LIT_MEM")
var OrLitToMem = LitToMem("or", "OR_LIT_MEM")
var XorLitToMem = LitToMem("xor", "XOR_LIT_MEM")

func RegPtrToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[RegPtrToRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			fmt.Printf("Error building parser for %s regptr-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			fmt.Printf("Error parsing with %s regptr-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		fmt.Printf("Parsed regptr-to-reg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

var MovRegPtrToReg = RegPtrToReg("mov", "MOV_REG_PTR_REG")

//var AddLitMem = LitToMem("add", "ADD_LIT_MEM")
//var SumLitMem = LitToMem("sub", "SUM_LIT_MEM")
//var AndLitMem = LitToMem("and", "AND_LIT_MEM")
//var OrLitMem = LitToMem("or", "OR_LIT_MEM")
//var XorLitMem = LitToMem("xor", "XOR_LIT_MEM")

func LitOffToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[LitOffToRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			fmt.Printf("Error building parser for %s litOff-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			fmt.Printf("Error parsing with %s litoff-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		fmt.Printf("Parsed litof-to-reg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

var MovLitOffToReg = LitOffToReg("mov", "MOV_LIT_OFF_REG")
var AddLitOffToReg = LitOffToReg("add", "ADD_LIT_OFF_REG")
var SumLitOffToReg = LitOffToReg("sub", "SUM_LIT_OFF_REG")
var AndLitOffToReg = LitOffToReg("and", "AND_LIT_OFF_REG")
var OrLitOffToReg = LitOffToReg("or", "OR_LIT_OFF_REG")
var XorLitOffToReg = LitOffToReg("xor", "XOR_LIT_OFF_REG")

func ParseInstruction(input string) (*Node, error) {
	var parsers = []Parser{
		// MOV instructions
		{"MovLitToReg", MovLitToReg},
		{"MovRegToReg", MovRegToReg},
		{"MovRegToMem", MovRegToMem},
		{"MovMemToReg", MovMemToReg},
		{"MovLitToMem", MovLitToMem},
		{"MovRegPtrToReg", MovRegPtrToReg},
		{"MovLitOffToReg", MovLitOffToReg},

		// ADD instructions
		{"AddLitToReg", AddLitToReg},
		{"AddRegToReg", AddRegToReg},
		{"AddRegToMem", AddRegToMem},
		{"AddMemToReg", AddMemToReg},
		{"AddLitToMem", AddLitToMem},
		//	{"AddRegPtrToReg", AddRegPtrToReg},
		{"AddLitOffToReg", AddLitOffToReg},

		// Future operations (commented out for now)
		// {"SubLitToReg",    SubLitToReg},
		// {"AndLitToReg",    AndLitToReg},
		// {"OrLitToReg",     OrLitToReg},
		// {"XorLitToReg",    XorLitToReg},
		// Note: Consider adding RegToReg versions for each operation above
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
