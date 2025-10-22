package ITypeInstructions

import (
	"fmt"
	isa "riscv-instruction-encoder/instructions"
)

type Type struct {
	isa.BaseInstruction
	OpCode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

func (i *Type) Decode(inst uint32) isa.Instruction {
	i.OpCode = uint8(inst & 0x7F)
	i.Rd = uint8((inst >> 7) & 0x1F)
	i.Funct3 = uint8((inst >> 12) & 0x7)
	i.Rs1 = uint8((inst >> 15) & 0x1F)
	i.Imm = uint16((inst >> 20) & 0xFFF)
	return i
}

func (i *Type) String() string {
	return fmt.Sprintf("formato = I {opcode=%02X, rd=%d, funct3=%d, rs1=%d, imm=%d}",
		i.OpCode, i.Rd, i.Funct3, i.Rs1, i.Imm)
}
