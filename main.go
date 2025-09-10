package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

type Instruction interface {
	Encode() uint32
	String() string
	Decode(inst uint32) Instruction
}

type RType struct {
	Opcode uint8 // 7 bits
	Rd     uint8 // 5 bits
	Funct3 uint8 // 3 bits
	Rs1    uint8 // 5 bits
	Rs2    uint8 // 5 bits
	Funct7 uint8 // 7 bits
}

type IType struct {
	OpCode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

type SType struct {
	Opcode uint8  // 7 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Rs2    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

type BType struct {
	Opcode uint8  // 7 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Rs2    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

type UType struct {
	Opcode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Imm    uint32 // 20 bits
}

type JType struct {
	Opcode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Imm    uint32 // 20 bits
}

func (r *RType) Encode() uint32 {
	var inst uint32
	inst |= uint32(r.Opcode & 0x7F)
	inst |= uint32(r.Rd&0x1F) << 7
	inst |= uint32(r.Funct3&0x7) << 12
	inst |= uint32(r.Rs1&0x1F) << 15
	inst |= uint32(r.Rs2&0x1F) << 20
	inst |= uint32(r.Funct7&0x7F) << 25
	return inst
}

func (r *RType) Decode(inst uint32) Instruction {
	r.Opcode = uint8(inst & 0x7F)
	r.Rd = uint8((inst >> 7) & 0x1F)
	r.Funct3 = uint8((inst >> 12) & 0x7)
	r.Rs1 = uint8((inst >> 15) & 0x1F)
	r.Rs2 = uint8((inst >> 20) & 0x1F)
	r.Funct7 = uint8((inst >> 25) & 0x7F)
	return r
}

func (r *RType) String() string {
	return fmt.Sprintf("RType {opcode=%02X, rd=%d, funct3=%d, rs1=%d, rs2=%d, funct7=%d}",
		r.Opcode, r.Rd, r.Funct3, r.Rs1, r.Rs2, r.Funct7)
}

func (i *IType) Encode() uint32 {
	var inst uint32
	inst |= uint32(i.OpCode & 0x7F)
	inst |= uint32(i.Rd&0x1F) << 7
	inst |= uint32(i.Funct3&0x7) << 12
	inst |= uint32(i.Rs1&0x1F) << 15
	inst |= uint32(i.Imm&0xFFF) << 20
	return inst
}

func (i *IType) Decode(inst uint32) Instruction {
	i.OpCode = uint8(inst & 0x7F)
	i.Rd = uint8((inst >> 7) & 0x1F)
	i.Funct3 = uint8((inst >> 12) & 0x7)
	i.Rs1 = uint8((inst >> 15) & 0x1F)
	i.Imm = uint16((inst >> 20) & 0xFFF)
	return i
}

func (i *IType) String() string {
	return fmt.Sprintf("IType {opcode=%02X, rd=%d, funct3=%d, rs1=%d, imm=%d}",
		i.OpCode, i.Rd, i.Funct3, i.Rs1, i.Imm)
}

func DecodeInstruction(inst uint32) Instruction {
	op := uint8(inst & 0x7F)

	switch op {
	case OpRType: // R-type
		return new(RType).Decode(inst)
	case OpIType1, OpIType2, OpIType3, OpIType4: // I-type
		return new(IType).Decode(inst)
	// case 0x23: return new(SType).Decode(inst)
	// case 0x63: return new(BType).Decode(inst)
	// case 0x37, 0x17: return new(UType).Decode(inst)
	// case 0x6F: return new(JType).Decode(inst)
	default:
		return nil
	}
}

func DecodeInstructionFromUInt32(encodedInstructions []uint32) {
	for _, inst := range encodedInstructions {
		decoded := DecodeInstruction(inst)
		if decoded != nil {
			fmt.Println(decoded.String())
		} else {
			fmt.Printf("Opcode %02X não reconhecido\n", inst&0x7F)
		}
	}
}

func DecodeFromBinaryFile() []uint32 {
	var instructions []uint32
	file, err := os.Open("fib_rec_binario.txt")
	if err != nil {
		log.Fatalf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		num, err := strconv.ParseUint(row, 2, 32)
		if err != nil {
			panic(err)
		}
		instructions = append(instructions, uint32(num))
	}

	return instructions
}

func DecodeFromHexFile() []uint32 {
	var instructions []uint32
	file, err := os.Open("fib_rec_hexadecimal.txt")
	if err != nil {
		log.Fatalf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		num, err := strconv.ParseUint(row, 16, 32)
		if err != nil {
			panic(err)
		}
		instructions = append(instructions, uint32(num))
	}

	return instructions
}

func main() {
	var instructionsFromBinaryFile []uint32 = DecodeFromBinaryFile()
	var instructionsFromHexFile []uint32 = DecodeFromHexFile()

	// for i, v := range instructionsFromBinaryFile {
	// 	fmt.Printf("nums[%d] = %d\n", i, v)
	// }

	// for i, v := range instructionsFromHexFile {
	// 	fmt.Printf("nums[%d] = %d\n", i, v)
	// }

	DecodeInstructionFromUInt32(instructionsFromBinaryFile)
	DecodeInstructionFromUInt32(instructionsFromHexFile)
}
