package STypeInstructions

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

func (s *Type) Decode(inst uint32) isa.Instruction {
	s.Opcode = uint8(inst & 0x7F)
	imm4_0 := (inst >> 7) & 0x1F
	s.Funct3 = uint8((inst >> 12) & 0x7)
	s.Rs1 = uint8((inst >> 15) & 0x1F)
	s.Rs2 = uint8((inst >> 20) & 0x1F)
	imm11_5 := (inst >> 25) & 0x7F
	s.Imm = uint16((imm11_5 << 5) | imm4_0)
	return s
}

func (s *Type) String() string {
	return fmt.Sprintf("formato = S {opcode=%02X, funct3=%d, rs1=%d, rs2=%d, imm=%d}",
		s.Opcode, s.Funct3, s.Rs1, s.Rs2, s.Imm)
}
