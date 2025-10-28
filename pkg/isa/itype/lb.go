package itype

import isa "riscv-instruction-encoder/pkg/isa"

type LB struct {
	Type
}

func newLB(t Type) *LB {
	inst := &LB{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "LB",
		OpCode:         uint32(t.OpCode),
		IsLoad:         true,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1)},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.MEM,
		ConsumeStage:   isa.ID,
	}
	return inst
}
