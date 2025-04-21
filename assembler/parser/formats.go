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
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

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
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

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
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

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
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

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
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

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
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

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
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

func NoArg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[NoArgsInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

func SingleReg(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[SingleRegInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}

func SingleLit(mnemonic, instructionType string) func(string) (*Node, error) {
	return func(input string) (*Node, error) {
		parser, err := participle.Build[SingleLitInstruction](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
		if err != nil {
			return nil, err
		}

		instr, err := parser.ParseString("", input)
		if err != nil {
			return nil, err
		}

		if !strings.EqualFold(instr.Instr, mnemonic) {
			return nil, fmt.Errorf("expected instruction %s, got %s", mnemonic, instr.Instr)
		}

		return instr.AsNode(instructionType), nil
	}
}
