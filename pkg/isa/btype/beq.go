package btype

import isa "riscv-instruction-encoder/pkg/isa"

type BEQ struct {
	Type
}

func newBEQ(t Type) *BEQ {
	inst := &BEQ{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "BEQ",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       true,
		IsJump:         false,
		WritesRegister: false,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1), int(t.Rs2)},
		Rd:             nil,
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}
