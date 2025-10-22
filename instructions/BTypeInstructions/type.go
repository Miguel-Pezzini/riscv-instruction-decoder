package BTypeInstructions

import (
	"fmt"
	isa "riscv-instruction-encoder/instructions"
)

type Type struct {
	isa.BaseInstruction
	Opcode uint8  // 7 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Rs2    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

func (b *Type) Decode(inst uint32) isa.Instruction {
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

func (b *Type) String() string {
	return fmt.Sprintf("formato = B {opcode=%02X, funct3=%d, rs1=%d, rs2=%d, imm=%d}",
		b.Opcode, b.Funct3, b.Rs1, b.Rs2, b.Imm)
}
