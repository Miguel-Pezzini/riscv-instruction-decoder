package RTypeInstructions

import (
	"fmt"
	isa "riscv-instruction-encoder/pkg/isa"
)

type Type struct {
	isa.BaseInstruction
	Opcode uint8 // 7 bits
	Rd     uint8 // 5 bits
	Funct3 uint8 // 3 bits
	Rs1    uint8 // 5 bits
	Rs2    uint8 // 5 bits
	Funct7 uint8 // 7 bits
}

func (r *Type) Decode(inst uint32) isa.Instruction {
	r.Opcode = uint8(inst & 0x7F)
	r.Rd = uint8((inst >> 7) & 0x1F)
	r.Funct3 = uint8((inst >> 12) & 0x7)
	r.Rs1 = uint8((inst >> 15) & 0x1F)
	r.Rs2 = uint8((inst >> 20) & 0x1F)
	r.Funct7 = uint8((inst >> 25) & 0x7F)

	return r.findInstruction(r.Funct3, r.Funct7)
}

func (r *Type) String() string {
	return fmt.Sprintf("formato = R {opcode=%02X, rd=%d, funct3=%d, rs1=%d, rs2=%d, funct7=%d}",
		r.Opcode, r.Rd, r.Funct3, r.Rs1, r.Rs2, r.Funct7)
}

func (r *Type) findInstruction(funct3 uint8, funct7 uint8) isa.Instruction {
	switch {
	case funct7 == 0x00 && funct3 == 0x00:
		return &ADD{Type: *r}
	}

	return r
}

func (r *Type) GetRegisterUsage() isa.RegisterUsage {
	return isa.RegisterUsage{
		ReadRegs:  []uint8{r.Rs1, r.Rs2},
		WriteRegs: []uint8{r.Rd},
	}
}
