package JTypeInstructions

import (
	"fmt"
	isa "riscv-instruction-encoder/instructions"
)

const (
	JAL = 0x6F
)

type Type struct {
	isa.BaseInstruction
	OpCode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Imm    uint32 // 20 bits
}

func (j *Type) GetRegisterUsage() isa.RegisterUsage {
	return isa.RegisterUsage{
		ReadRegs:  nil,
		WriteRegs: []uint8{j.Rd},
	}
}

func (j *Type) Decode(inst uint32) isa.Instruction {
	j.OpCode = uint8(inst & 0x7F)
	j.Rd = uint8((inst >> 7) & 0x1F)
	j.Imm = uint32(inst>>12) & 0xFFFFF

	return j.findInstruction()
}

func (j *Type) String() string {
	return fmt.Sprintf("%s {opcode=%02X, rd=%d, imm=%d}",
		j.getInstructionName(), j.OpCode, j.Rd, j.Imm)
}

func (j *Type) findInstruction() isa.Instruction {
	switch j.OpCode {
	case JAL:
		return &jalInstruction{*j}
	}
	return j
}

func (j *Type) getInstructionName() string {
	switch j.OpCode {
	case JAL:
		return "JAL"
	}
	return "UNKNOWN_J"
}

func (j *Type) ExecuteFetchInstruction() {
	fmt.Printf("[IF ] Fetching instruction: %s\n", j.getInstructionName())
}

func (j *Type) ExecuteDecodeInstruction() {
	fmt.Printf("[ID ] Decoding instruction: %s\n", j.getInstructionName())
}

func (j *Type) ExecuteOperation() {
	fmt.Printf("[EX ] Executing operation for instruction: %s\n", j.getInstructionName())
}

func (j *Type) ExecuteAccessOperand() {
	fmt.Printf("[MEM] Accessing operands/memory for instruction: %s\n", j.getInstructionName())
}

func (j *Type) ExecuteWriteBack() {
	fmt.Printf("[WB ] Writing back result of instruction: %s\n", j.getInstructionName())
}
