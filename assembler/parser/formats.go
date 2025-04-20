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
			//		fmt.Printf("Error building parser for %s: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//		fmt.Printf("Error parsing with %s: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//	fmt.Printf("Parsed instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

func RegToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[RegRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			//	fmt.Printf("Error building parser for %s reg-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//	fmt.Printf("Error parsing with %s reg-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//fmt.Printf("Parsed reg-to-reg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

func RegToMem(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[RegMemInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			//	fmt.Printf("Error building parser for %s reg-to-mem: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//	fmt.Printf("Error parsing with %s reg-to-mem: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//fmt.Printf("Parsed reg-to-mem instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

func MemToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[MemRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			//	fmt.Printf("Error building parser for %s mem-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//	fmt.Printf("Error parsing with %s mem-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//fmt.Printf("Parsed mem-to-reg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

func LitToMem(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[LitMemInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			//	fmt.Printf("Error building parser for %s lit-to-mem: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//	fmt.Printf("Error parsing with %s lit-to-mem: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//fmt.Printf("Parsed lit-to-mem instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

func RegPtrToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[RegPtrToRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			//	fmt.Printf("Error building parser for %s regptr-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//	fmt.Printf("Error parsing with %s regptr-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//		/fmt.Printf("Parsed regptr-to-reg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)
		//
		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

func LitOffToReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[LitOffToRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			//			fmt.Printf("Error building parser for %s litOff-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//			fmt.Printf("Error parsing with %s litoff-to-reg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//		fmt.Printf("Parsed litof-to-reg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

// NoArg handles operations that don't take any arguments (like HLT, RET)
func NoArg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[NoArgsInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			//			fmt.Printf("Error building parser for %s no-arg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//			fmt.Printf("Error parsing with %s no-arg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//		fmt.Printf("Parsed no-arg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

// SingleReg handles operations that operate only on a single register
func SingleReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[SingleRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			//			fmt.Printf("Error building parser for %s single-reg: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//			fmt.Printf("Error parsing with %s single-reg: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//		fmt.Printf("Parsed single-reg instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

// SingleReg handles operations that operate only on a single register
func SingleLit(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[SingleLitInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			//			fmt.Printf("Error building parser for %s single-lit: %v\n", mnemonic, err)
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			//			fmt.Printf("Error parsing with %s single-lit: %v\n", mnemonic, err)
			return nil, err
		}

		// Debug what instruction was actually parsed
		//		fmt.Printf("Parsed single-lit instruction name: %s, expected: %s\n", instr.Instr, mnemonic)

		// Validate that the parsed instruction matches our expected mnemonic
		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}
