package jtype

import isa "riscv-instruction-encoder/pkg/isa"

type JAL struct {
	Type
}

func newJAL(t Type) *JAL {
	inst := &JAL{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "JAL",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         true,
		WritesRegister: true,
		ReadsRegister:  false,
		Rs:             nil,
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}
