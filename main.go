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
	String() string
	Decode(inst uint32) Instruction
}

type RawInstruction struct {
	Origin string
	Value  uint32
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
	return fmt.Sprintf("formato = R {opcode=%02X, rd=%d, funct3=%d, rs1=%d, rs2=%d, funct7=%d}",
		r.Opcode, r.Rd, r.Funct3, r.Rs1, r.Rs2, r.Funct7)
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
	return fmt.Sprintf("formato = I {opcode=%02X, rd=%d, funct3=%d, rs1=%d, imm=%d}",
		i.OpCode, i.Rd, i.Funct3, i.Rs1, i.Imm)
}
func (s *SType) Decode(inst uint32) Instruction {
	s.Opcode = uint8(inst & 0x7F)
	imm4_0 := (inst >> 7) & 0x1F
	s.Funct3 = uint8((inst >> 12) & 0x7)
	s.Rs1 = uint8((inst >> 15) & 0x1F)
	s.Rs2 = uint8((inst >> 20) & 0x1F)
	imm11_5 := (inst >> 25) & 0x7F
	s.Imm = uint16((imm11_5 << 5) | imm4_0)
	return s
}

func (s *SType) String() string {
	return fmt.Sprintf("formato = S {opcode=%02X, funct3=%d, rs1=%d, rs2=%d, imm=%d}",
		s.Opcode, s.Funct3, s.Rs1, s.Rs2, s.Imm)
}

func (b *BType) Decode(inst uint32) Instruction {
	b.Opcode = uint8(inst & 0x7F)
	imm11 := (inst >> 7) & 0x1
	imm4_1 := (inst >> 8) & 0xF
	b.Funct3 = uint8((inst >> 12) & 0x7)
	b.Rs1 = uint8((inst >> 15) & 0x1F)
	b.Rs2 = uint8((inst >> 20) & 0x1F)
	imm10_5 := (inst >> 25) & 0x3F
	imm12 := (inst >> 31) & 0x1
	b.Imm = uint16((imm12 << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1))
	return b
}

func (b *BType) String() string {
	return fmt.Sprintf("formato = B {opcode=%02X, funct3=%d, rs1=%d, rs2=%d, imm=%d}",
		b.Opcode, b.Funct3, b.Rs1, b.Rs2, b.Imm)
}

func (u *UType) Decode(inst uint32) Instruction {
	u.Opcode = uint8(inst & 0x7F)
	u.Rd = uint8((inst >> 7) & 0x1F)
	u.Imm = uint32(inst>>12) & 0xFFFFF
	return u
}

func (u *UType) String() string {
	return fmt.Sprintf("formato = U {opcode=%02X, rd=%d, imm=%d}",
		u.Opcode, u.Rd, u.Imm)
}

func (j *JType) Decode(inst uint32) Instruction {
	j.Opcode = uint8(inst & 0x7F)
	j.Rd = uint8((inst >> 7) & 0x1F)
	j.Imm = uint32(inst>>12) & 0xFFFFF

	return j
}

func (j *JType) String() string {
	return fmt.Sprintf("formato = J {opcode=%02X, rd=%d, imm=%d}",
		j.Opcode, j.Rd, j.Imm)
}

func DecodeInstruction(inst uint32) Instruction {
	op := uint8(inst & 0x7F)

	switch op {
	case OpRType:
		return new(RType).Decode(inst)
	case OpIType1, OpIType2, OpIType3, OpIType4:
		return new(IType).Decode(inst)
	case OpSType:
		return new(SType).Decode(inst)
	case OpBType:
		return new(BType).Decode(inst)
	case OpUType1, OpUType2:
		return new(UType).Decode(inst)
	case OpJType:
		return new(JType).Decode(inst)
	default:
		return nil
	}
}

func DecodeInstructionFromUInt32(encodedInstructions []RawInstruction) {
	for _, inst := range encodedInstructions {
		decoded := DecodeInstruction(inst.Value)
		if decoded != nil {
			fmt.Printf("instrução %s -> %s\n", inst.Origin, decoded.String())
		} else {
			fmt.Printf("Opcode %02X não reconhecido\n", inst.Value&0x7F)
		}
	}
}

func DecodeFromBinaryFile() []RawInstruction {
	var instructions []RawInstruction
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
		instructions = append(instructions, RawInstruction{
			Origin: row,
			Value:  uint32(num),
		})
	}

	return instructions
}

func DecodeFromHexFile() []RawInstruction {
	var instructions []RawInstruction
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
		instructions = append(instructions, RawInstruction{
			Origin: row,
			Value:  uint32(num),
		})
	}

	return instructions
}

func main() {
	var instructionsFromBinaryFile []RawInstruction = DecodeFromBinaryFile()
	var instructionsFromHexFile []RawInstruction = DecodeFromHexFile()

	DecodeInstructionFromUInt32(instructionsFromBinaryFile)
	DecodeInstructionFromUInt32(instructionsFromHexFile)
}
