package ITypeInstructions

import isa "riscv-instruction-encoder/pkg/isa"

type ORI struct {
	Type
}

func NewORI(t Type) *ORI {
	inst := &ORI{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "ORI",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		WritesRegister: true,
		Rs:             []int{int(t.Rs1)},
		Rd:             intPtr(int(t.Rd)),
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}
	return inst
}
