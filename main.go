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

func (s *SType) Encode() uint32 {
	var inst uint32
	inst |= uint32(s.Opcode & 0x7F)
	inst |= uint32(s.Imm&0x1F) << 7          // imm[4:0]
	inst |= uint32(s.Funct3&0x7) << 12
	inst |= uint32(s.Rs1&0x1F) << 15
	inst |= uint32(s.Rs2&0x1F) << 20
	inst |= uint32((s.Imm>>5)&0x7F) << 25    // imm[11:5]
	return inst
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
	return fmt.Sprintf("SType {opcode=%02X, funct3=%d, rs1=%d, rs2=%d, imm=%d}",
		s.Opcode, s.Funct3, s.Rs1, s.Rs2, s.Imm)
}

func (b *BType) Encode() uint32 {
	var inst uint32
	inst |= uint32(b.Opcode & 0x7F)
	inst |= (uint32(b.Imm>>11) & 0x1) << 7     // imm[11]
	inst |= (uint32(b.Imm>>1) & 0xF) << 8      // imm[4:1]
	inst |= uint32(b.Funct3&0x7) << 12
	inst |= uint32(b.Rs1&0x1F) << 15
	inst |= uint32(b.Rs2&0x1F) << 20
	inst |= (uint32(b.Imm>>5) & 0x3F) << 25    // imm[10:5]
	inst |= (uint32(b.Imm>>12) & 0x1) << 31    // imm[12]
	return inst
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
	return fmt.Sprintf("BType {opcode=%02X, funct3=%d, rs1=%d, rs2=%d, imm=%d}",
		b.Opcode, b.Funct3, b.Rs1, b.Rs2, b.Imm)
}

func (u *UType) Encode() uint32 {
	var inst uint32
	inst |= uint32(u.Opcode & 0x7F)
	inst |= uint32(u.Rd&0x1F) << 7
	inst |= (u.Imm & 0xFFFFF) << 12 // 20 bits
	return inst
}

func (u *UType) Decode(inst uint32) Instruction {
	u.Opcode = uint8(inst & 0x7F)
	u.Rd = uint8((inst >> 7) & 0x1F)
	u.Imm = uint32(inst >> 12) & 0xFFFFF
	return u
}

func (u *UType) String() string {
	return fmt.Sprintf("UType {opcode=%02X, rd=%d, imm=%d}",
		u.Opcode, u.Rd, u.Imm)
}

func (j *JType) Encode() uint32 {
	var inst uint32
	inst |= uint32(j.Opcode & 0x7F)
	inst |= uint32(j.Rd&0x1F) << 7
	inst |= (j.Imm & 0xFF000)          // imm[19:12]
	inst |= (j.Imm & 0x800) << 9       // imm[11]
	inst |= (j.Imm & 0x7FE) << 20      // imm[10:1]
	inst |= (j.Imm & 0x100000) << 11   // imm[20]
	return inst
}

func (j *JType) Decode(inst uint32) Instruction {
	j.Opcode = uint8(inst & 0x7F)
	j.Rd = uint8((inst >> 7) & 0x1F)

	imm19_12 := (inst >> 12) & 0xFF
	imm11 := (inst >> 20) & 0x1
	imm10_1 := (inst >> 21) & 0x3FF
	imm20 := (inst >> 31) & 0x1

	j.Imm = uint32((imm20 << 20) | (imm19_12 << 12) | (imm11 << 11) | (imm10_1 << 1))
	return j
}

func (j *JType) String() string {
	return fmt.Sprintf("JType {opcode=%02X, rd=%d, imm=%d}",
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

	DecodeInstructionFromUInt32(instructionsFromBinaryFile)
	DecodeInstructionFromUInt32(instructionsFromHexFile)
}
