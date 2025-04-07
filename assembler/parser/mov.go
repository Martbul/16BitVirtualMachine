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

// ParseMovInstruction parses both "mov $42, r4" and "mov r1, r2" patterns
//func ParseMovInstruction(input string) (*Node, error) {
//	// Try MovLitToReg first
//	litNode, litErr := MovLitToReg(input)
//	if litErr == nil {
//		return litNode, nil
//	}

// If that fails, try MovRegToReg
//	regNode, regErr := MovRegToReg(input)
//	if regErr == nil {
//		return regNode, nil
//	}

// If both fail, return the most appropriate error
//	return nil, fmt.Errorf("failed to parse mov instruction: %v, %v", litErr, regErr)
//}/
