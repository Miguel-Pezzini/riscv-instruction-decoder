package UTypeInstructions

import (
	"fmt"
	isa "riscv-instruction-encoder/pkg/isa"
)

type Type struct {
	isa.BaseInstruction
	Opcode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Imm    uint32 // 20 bits
}

func (u *Type) Decode(inst uint32) isa.Instruction {
	u.Opcode = uint8(inst & 0x7F)
	u.Rd = uint8((inst >> 7) & 0x1F)
	u.Imm = uint32(inst>>12) & 0xFFFFF
	u.InstructionMeta = isa.InstructionMeta{}
	return u
}

func (u *Type) String() string {
	return fmt.Sprintf("formato = U {opcode=%02X, rd=%d, imm=%d}",
		u.Opcode, u.Rd, u.Imm)
}
