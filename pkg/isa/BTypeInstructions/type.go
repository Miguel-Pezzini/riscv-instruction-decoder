package BTypeInstructions

import (
	"fmt"
	isa "riscv-instruction-encoder/pkg/isa"
)

// Opcodes
const (
	BRANCH = 0x63
)

// Funct3
const (
	FUNCT3_BEQ  = 0x0
	FUNCT3_BNE  = 0x1
	FUNCT3_BLT  = 0x4
	FUNCT3_BGE  = 0x5
	FUNCT3_BLTU = 0x6
	FUNCT3_BGEU = 0x7
)

type Type struct {
	isa.BaseInstruction
	OpCode uint8  // 7 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Rs2    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

func (b *Type) Decode(inst uint32) isa.Instruction {
	b.OpCode = uint8(inst & 0x7F)
	imm11 := (inst >> 7) & 0x1
	imm4_1 := (inst >> 8) & 0xF
	b.Funct3 = uint8((inst >> 12) & 0x7)
	b.Rs1 = uint8((inst >> 15) & 0x1F)
	b.Rs2 = uint8((inst >> 20) & 0x1F)
	imm10_5 := (inst >> 25) & 0x3F
	imm12 := (inst >> 31) & 0x1
	b.Imm = uint16((imm12 << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1))
	return b.findInstruction()
}

func (b *Type) getInstructionName() string {
	switch b.OpCode {
	case BRANCH:
		switch b.Funct3 {
		case FUNCT3_BLT:
			return "BLT"
		case FUNCT3_BGE:
			return "BGE"
		case FUNCT3_BEQ:
			return "BEQ"
		case FUNCT3_BNE:
			return "BNE"
		}
	}
	return "UNKNOWN_B"
}

func (b *Type) findInstruction() isa.Instruction {
	switch b.OpCode {
	case BRANCH:
		return &bltInstruction{*b} // default para BLT, adapte se quiser outros
	}
	return b
}

func (b *Type) String() string {
	return fmt.Sprintf("%s {opcode=%02X, funct3=%d, rs1=%d, rs2=%d, imm=%d}",
		b.getInstructionName(), b.OpCode, b.Funct3, b.Rs1, b.Rs2, b.Imm)
}

func (b *Type) GetRegisterUsage() isa.RegisterUsage {
	// B-type lê Rs1 e Rs2, não escreve
	return isa.RegisterUsage{
		ReadRegs:  []uint8{b.Rs1, b.Rs2},
		WriteRegs: []uint8{},
	}
}

// Pipeline stages
func (b *Type) ExecuteFetchInstruction() {
	fmt.Printf("[IF ] Fetching instruction: %s\n", b.getInstructionName())
}

func (b *Type) ExecuteDecodeInstruction() {
	fmt.Printf("[ID ] Decoding instruction: %s\n", b.getInstructionName())
}

func (b *Type) ExecuteOperation() {
	fmt.Printf("[EX ] Executing operation for instruction: %s\n", b.getInstructionName())
}

func (b *Type) ExecuteAccessOperand() {
	fmt.Printf("[MEM] Accessing operands/memory for instruction: %s\n", b.getInstructionName())
}

func (b *Type) ExecuteWriteBack() {
	fmt.Printf("[WB ] Writing back result of instruction: %s\n", b.getInstructionName())
}
