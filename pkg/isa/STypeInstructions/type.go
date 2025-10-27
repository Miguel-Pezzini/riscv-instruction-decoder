package STypeInstructions

import (
	"fmt"
	isa "riscv-instruction-encoder/pkg/isa"
)

// Opcodes
const (
	STORE = 0x23
)

// Funct3
const (
	FUNCT3_SB = 0x0
	FUNCT3_SH = 0x1
	FUNCT3_SW = 0x2
)

type Type struct {
	isa.BaseInstruction
	OpCode uint8  // 7 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Rs2    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

func (s *Type) Decode(inst uint32) isa.Instruction {
	s.OpCode = uint8(inst & 0x7F)
	imm4_0 := (inst >> 7) & 0x1F
	s.Funct3 = uint8((inst >> 12) & 0x7)
	s.Rs1 = uint8((inst >> 15) & 0x1F)
	s.Rs2 = uint8((inst >> 20) & 0x1F)
	imm11_5 := (inst >> 25) & 0x7F
	s.Imm = uint16((imm11_5 << 5) | imm4_0)
	return s.findInstruction()
}

func (s *Type) String() string {
	return fmt.Sprintf("%s {opcode=%02X, funct3=%d, rs1=%d, rs2=%d, imm=%d}",
		s.getInstructionName(), s.OpCode, s.Funct3, s.Rs1, s.Rs2, s.Imm)
}

func (s *Type) getInstructionName() string {
	switch s.OpCode {
	case STORE:
		switch s.Funct3 {
		case FUNCT3_SB:
			return "SB"
		case FUNCT3_SH:
			return "SH"
		case FUNCT3_SW:
			return "SW"
		}
	}
	return "UNKNOWN_S"
}

func (s *Type) findInstruction() isa.Instruction {
	switch s.OpCode {
	case STORE:
		return &SW{*s} // default para SW, adapte se quiser SH/SB
	}
	return s
}

// Pipeline stages
func (s *Type) ExecuteFetchInstruction() {
	fmt.Printf("[IF ] Fetching instruction: %s\n", s.getInstructionName())
}

func (s *Type) ExecuteDecodeInstruction() {
	fmt.Printf("[ID ] Decoding instruction: %s\n", s.getInstructionName())
}

func (s *Type) ExecuteOperation() {
	fmt.Printf("[EX ] Executing operation for instruction: %s\n", s.getInstructionName())
}

func (s *Type) ExecuteAccessOperand() {
	fmt.Printf("[MEM] Accessing operands/memory for instruction: %s\n", s.getInstructionName())
}

func (s *Type) ExecuteWriteBack() {
	fmt.Printf("[WB ] Writing back result of instruction: %s\n", s.getInstructionName())
}
