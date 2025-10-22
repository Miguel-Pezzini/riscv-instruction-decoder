package decoder

import (
	"bufio"
	"fmt"
	"log"
	"os"
	isa "riscv-instruction-encoder/instructions"
	"riscv-instruction-encoder/instructions/BTypeInstructions"
	"riscv-instruction-encoder/instructions/ITypeInstructions"
	"riscv-instruction-encoder/instructions/JTypeInstructions"
	"riscv-instruction-encoder/instructions/RTypeInstructions"
	"riscv-instruction-encoder/instructions/STypeInstructions"
	"riscv-instruction-encoder/instructions/UTypeInstructions"
	"strconv"
)

const (
	OpRType  = 0x33
	OpIType1 = 0x13
	OpIType2 = 0x03
	OpIType3 = 0x67
	OpIType4 = 0x73
	OpSType  = 0x23
	OpBType  = 0x63
	OpUType1 = 0x37
	OpUType2 = 0x17
	OpJType  = 0x6F
)

const (
	FORMAT_BIN = "bin"
	FORMAT_HEX = "hex"
)

func DecodeInstruction(inst uint32) isa.Instruction {
	op := uint8(inst & 0x7F)

	switch op {
	case OpRType:
		return new(RTypeInstructions.Type).Decode(inst)
	case OpIType1, OpIType2, OpIType3, OpIType4:
		return new(ITypeInstructions.Type).Decode(inst)
	case OpSType:
		return new(STypeInstructions.Type).Decode(inst)
	case OpBType:
		return new(BTypeInstructions.Type).Decode(inst)
	case OpUType1, OpUType2:
		return new(UTypeInstructions.Type).Decode(inst)
	case OpJType:
		return new(JTypeInstructions.Type).Decode(inst)
	default:
		return nil
	}
}

func DecodeInstructionFromUInt32(encodedInstructions []isa.RawInstruction) []isa.Instruction {
	var instructions = make([]isa.Instruction, len(encodedInstructions))
	for i, inst := range encodedInstructions {
		decoded := DecodeInstruction(inst.Value)
		if decoded != nil {
			instructions[i] = decoded
			fmt.Printf("instrução %s -> %s\n", inst.Origin, decoded.String())
		} else {
			fmt.Printf("Opcode %02X não reconhecido\n", inst.Value&0x7F)
		}
	}

	return instructions
}

func DecodeFromFile(filePath string, format string) []isa.RawInstruction {
	var instructions []isa.RawInstruction

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var base int
	switch format {
	case FORMAT_BIN:
		base = 2
	case FORMAT_HEX:
		base = 16
	default:
		log.Fatalf("formato inválido: %s (use 'bin' ou 'hex')", format)
	}

	for scanner.Scan() {
		row := scanner.Text()
		num, err := strconv.ParseUint(row, base, 32)
		if err != nil {
			panic(err)
		}
		instructions = append(instructions, isa.RawInstruction{
			Origin: row,
			Value:  uint32(num),
		})
	}

	return instructions
}
