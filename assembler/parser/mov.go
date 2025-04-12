package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

// === Parser Constructors ===
// ParseMovLitToReg parses "mov $42, r4" and more complex expressions
func MovLitToReg(input string) (*Node, error) {
	parser, err := participle.Build[MovLitToRegInstruction](
		participle.Lexer(lexerDef),
		participle.Elide("Whitespace"),
	)
	if err != nil {
		return nil, err
	}

	movInstr, err := parser.ParseString("", input)
	if err != nil {
		return nil, err
	}

	return movInstr.AsNode(), nil
}

// MovRegToReg parses "mov r1, r2" expressions
func MovRegToReg(input string) (*Node, error) {
	// Build the parser for this specific instruction type
	parser, err := participle.Build[MovRegToRegInstruction](
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

// ParseMovInstruction tries all possible MOV instruction patterns and returns the first successful parse

func ParseMovInstruction(input string) (*Node, error) {
	// Define a slice of parser functions to try
	parsers := []struct {
		name string
		fn   func(string) (*Node, error)
	}{
		{"MovLitToReg", MovLitToReg}, //INFO: WORKS
		{"MovRegToReg", MovRegToReg}, //INFO WORKS
		{"MovRegToMem", MovRegToMem},
		{"MovMemToReg", MovMemToReg},
		{"MovLitToMem", MovLitToMem},
		{"MovRegPtrToReg", MovRegPtrToReg},
		{"MovLitOffToReg", MovLitOffToReg},
	}

	// Try each parser in sequence
	var errors []string
	for _, parser := range parsers {
		node, err := parser.fn(input)
		if err == nil {
			// Success! Return the parsed node
			return node, nil
		}
		// Collect error for reporting
		errors = append(errors, fmt.Sprintf("%s: %v", parser.name, err))
	}

	// If we reach here, all parsers failed
	return nil, fmt.Errorf("failed to parse mov instruction. Errors: %v", errors)
}
