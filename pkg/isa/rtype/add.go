package rtype

import isa "riscv-instruction-encoder/pkg/isa"

type ADD struct {
	Type
}

func newADD(t Type) *ADD {
	inst := &ADD{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "ADD",
		OpCode:         uint32(t.Opcode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{(int(t.Rs1)), (int(t.Rs2))},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}

	return inst
}
