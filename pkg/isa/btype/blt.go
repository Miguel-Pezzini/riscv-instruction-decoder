package btype

import isa "riscv-instruction-encoder/pkg/isa"

type BLT struct {
	Type
}

func newBLT(t Type) *BLT {
	inst := &BLT{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "BLT",
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
