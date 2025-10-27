package stype

import isa "riscv-instruction-encoder/pkg/isa"

type SW struct {
	Type
}

func newSW(t Type) *SW {
	inst := &SW{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "SW",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        true,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: false,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1), int(t.Rs2)},
		Rd:             nil,
		ProduceStage:   0,
		ConsumeStage:   isa.ID,
	}
	return inst
}
