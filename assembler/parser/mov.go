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

// === Parser Constructors ===
// Modified RegToReg and LitToReg functions with debug output
// First, update your LitToReg function with more detailed debug info
func LitToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		fmt.Println("input: ", input, "mnemonic: ", mnemonic)
		fmt.Printf("Trying to parse %s instruction: %s\n", mnemonic, input)

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

// Create instances for different instruction types
var MovLitToReg = LitToReg("mov", "MOV_LIT_REG")
var AddLitToReg = LitToReg("add", "ADD_LIT_REG")
var SubLitToReg = LitToReg("sub", "SUB_LIT_REG")
var AndLitToReg = LitToReg("and", "AND_LIT_REG")
var OrLitToReg = LitToReg("or", "OR_LIT_REG")
var XorLitToReg = LitToReg("xor", "XOR_LIT_REG")

// Update your RegToReg function to make sure it's correctly using RegRegInstruction
func RegToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		fmt.Printf("Trying to parse %s reg-to-reg instruction: %s\n", mnemonic, input)

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

// Create instances for different instruction types
var MovRegToReg = RegToReg("mov", "MOV_REG_REG")
var AddRegToReg = RegToReg("add", "ADD_REG_REG")
var SubRegToReg = RegToReg("sub", "SUB_REG_REG")
var AndRegToReg = RegToReg("and", "AND_REG_REG")
var OrRegToReg = RegToReg("or", "OR_REG_REG")
var XorRegToReg = RegToReg("xor", "XOR_REG_REG")

// MovRegToMem parses "mov r1, &$42" or "mov r1, &[$42 + !loc]" expressions
func MovRegToMem(input string) (*Node, error) {
	// Build the parser for this specific instruction type
	parser, err := participle.Build[MovRegToMemInstruction](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, err
	}

	// Parse the input string
	movInstr, err := parser.ParseString("", input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reg-to-mem instruction: %v", err)
	}

	return movInstr.AsNode(), nil
}

// MovRegToMem parses "mov r1, &$42" or "mov r1, &[$42 + !loc]" expressions
func MovMemToReg(input string) (*Node, error) {
	// Build the parser for this specific instruction type
	parser, err := participle.Build[MovMemToRegInstruction](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, err
	}

	// Parse the input string
	movInstr, err := parser.ParseString("", input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reg-to-mem instruction: %v", err)
	}

	return movInstr.AsNode(), nil
}

func MovLitToMem(input string) (*Node, error) {
	// Build the parser for this specific instruction type
	parser, err := participle.Build[MovLitToMemInstruction](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, err
	}

	// Parse the input string
	movInstr, err := parser.ParseString("", input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lit-to-mem instruction: %v", err)
	}

	return movInstr.AsNode(), nil
}

// MovRegToReg parses "mov r1, r2" expressions
func MovRegPtrToReg(input string) (*Node, error) {
	parser, err := participle.Build[MovRegPtrToRegInstruction](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, err
	}

	// Parse the input string
	movInstr, err := parser.ParseString("", input)
	if err != nil {
		return nil, err
	}

	return movInstr.AsNode(), nil
}

// MovRegToReg parses "mov r1, r2" expressions
func MovLitOffToReg(input string) (*Node, error) {
	// Build the parser for this specific instruction type
	parser, err := participle.Build[MovLitOffToRegInstruction](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, err
	}

	// Parse the input string
	movInstr, err := parser.ParseString("", input)
	if err != nil {
		return nil, err
	}

	return movInstr.AsNode(), nil
}

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
