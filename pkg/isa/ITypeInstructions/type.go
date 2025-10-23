package ITypeInstructions

import (
	"fmt"
	isa "riscv-instruction-encoder/pkg/isa"
)

// Definição de opcodes como constantes
const (
	OP_IMM = 0x13 // ADDI, ORI, ANDI, etc.
	LOAD   = 0x03 // LB, LW, etc.
	JALR   = 0x67
)

// Definição de funct3 para OP_IMM
const (
	FUNCT3_ADDI = 0x0
	FUNCT3_ORI  = 0x6
	FUNCT3_ANDI = 0x7
)

// Definição de funct3 para LOAD
const (
	FUNCT3_LB = 0x0
	FUNCT3_LW = 0x2
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
	return i.findInstruction()
}

func (i *Type) String() string {
	name := i.getInstructionName()
	return fmt.Sprintf("%s {opcode=%02X, rd=%d, funct3=%d, rs1=%d, imm=%d}",
		name, i.OpCode, i.Rd, i.Funct3, i.Rs1, i.Imm)
}

func (i *Type) getInstructionName() string {
	switch i.OpCode {
	case OP_IMM:
		switch i.Funct3 {
		case FUNCT3_ADDI:
			return "ADDI"
		case FUNCT3_ORI:
			return "ORI"
		case FUNCT3_ANDI:
			return "ANDI"
		}
	case LOAD:
		switch i.Funct3 {
		case FUNCT3_LB:
			return "LB"
		case FUNCT3_LW:
			return "LW"
		}
	case JALR:
		return "JALR"
	}
	return "UNKNOWN_I"
}

func (i *Type) findInstruction() isa.Instruction {
	switch i.OpCode {
	case OP_IMM:
		switch i.Funct3 {
		case FUNCT3_ADDI:
			return &addiInstruction{*i}
		case FUNCT3_ORI:
			return &oriInstruction{*i}
		case FUNCT3_ANDI:
			return &andiInstruction{*i}
		}
	case LOAD:
		switch i.Funct3 {
		case FUNCT3_LW:
			return &lwInstruction{*i}
		case FUNCT3_LB:
			return &lbInstruction{*i}
		}
	case JALR:
		return &jalrInstruction{*i}
	}
	return i
}

func (i *Type) GetRegisterUsage() isa.RegisterUsage {
	return isa.RegisterUsage{
		ReadRegs:  []uint8{i.Rs1},
		WriteRegs: []uint8{i.Rd},
	}
}

// Stages
func (t *Type) ExecuteFetchInstruction() {
	fmt.Printf("[IF ] Fetching instruction: %s\n", t.getInstructionName())
}

func (t *Type) ExecuteDecodeInstruction() {
	fmt.Printf("[ID ] Decoding instruction: %s\n", t.getInstructionName())
}

func (t *Type) ExecuteOperation() {
	fmt.Printf("[EX ] Executing operation for instruction: %s\n", t.getInstructionName())
}

func (t *Type) ExecuteAccessOperand() {
	fmt.Printf("[MEM] Accessing operands/memory for instruction: %s\n", t.getInstructionName())
}

func (t *Type) ExecuteWriteBack() {
	fmt.Printf("[WB ] Writing back result of instruction: %s\n", t.getInstructionName())
}
