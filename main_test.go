package main

import (
	"reflect"
	"testing"
)

func TestDecodeInstruction(t *testing.T) {
	hexTests := []struct {
		input    uint32
		expected Instruction
	}{
		{0x004000EF, &JType{Opcode: 0x6F, Rd: 1, Imm: 1024}},
		{0xFF810113, &IType{OpCode: 0x13, Rd: 2, Funct3: 0, Rs1: 2, Imm: 4088}},
		{0x002081B3, &RType{Opcode: 0x33, Rd: 3, Funct3: 0, Rs1: 1, Rs2: 2, Funct7: 0}},
		{0x00A12023, &SType{Opcode: 0x23, Funct3: 2, Rs1: 2, Rs2: 10, Imm: 0}},
		{0x00050663, &BType{Opcode: 0x63, Funct3: 0, Rs1: 10, Rs2: 0, Imm: 12}},
	}

	binTests := []struct {
		input    uint32
		expected Instruction
	}{
		{0b00000000000100000000000010010011, &IType{OpCode: 0x13, Rd: 1, Funct3: 0, Rs1: 0, Imm: 1}},
		{0b00000000001000001000000110010011, &IType{OpCode: 0x13, Rd: 3, Funct3: 0, Rs1: 1, Imm: 2}},
		{0b00000000101000100010001010100011, &SType{Opcode: 0x23, Funct3: 2, Rs1: 4, Rs2: 10, Imm: 5}},
		{0b00000000000000001000001101100011, &BType{Opcode: 0x63, Funct3: 0, Rs1: 1, Rs2: 0, Imm: 6}},
		{0b00000000010000000010000011101111, &JType{Opcode: 0x6F, Rd: 1, Imm: 1026}},
	}

	for _, tt := range hexTests {
		got := DecodeInstruction(tt.input)
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("HEX 0x%08X: expected %+v, received %+v", tt.input, tt.expected, got)
		}
	}

	for _, tt := range binTests {
		got := DecodeInstruction(tt.input)
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("BIN %032b: expected %+v, received %+v", tt.input, tt.expected, got)
		}
	}
}
