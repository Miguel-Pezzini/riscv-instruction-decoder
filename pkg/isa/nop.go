package isa

import (
	"fmt"
)

type NOP struct {
	BaseInstruction
}

func NewNOP() Instruction {
	inst := &NOP{BaseInstruction: BaseInstruction{}}
	inst.InstructionMeta = InstructionMeta{
		Name:           "NOP",
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		WritesRegister: false,
		ReadsRegister:  false,
	}
	return inst
}

func (i *NOP) Decode(inst uint32) Instruction {
	return NewNOP()
}

func (i *NOP) String() string {
	return fmt.Sprintf("name: %s",
		i.InstructionMeta.Name)
}

// Stages
func (t *NOP) ExecuteFetchInstruction() {
	fmt.Printf("[IF ] Fetching instruction: %s\n", t.InstructionMeta.Name)
}

func (t *NOP) ExecuteDecodeInstruction() {
	fmt.Printf("[ID ] Decoding instruction: %s\n", t.InstructionMeta.Name)
}

func (t *NOP) ExecuteOperation() {
	fmt.Printf("[EX ] Executing operation for instruction: %s\n", t.InstructionMeta.Name)
}

func (t *NOP) ExecuteAccessOperand() {
	fmt.Printf("[MEM] Accessing operands/memory for instruction: %s\n", t.InstructionMeta.Name)
}

func (t *NOP) ExecuteWriteBack() {
	fmt.Printf("[WB ] Writing back result of instruction: %s\n", t.InstructionMeta.Name)
}
