package itype

import isa "riscv-instruction-encoder/pkg/isa"

type ANDI struct {
	Type
}

func NewANDI(t Type) *ANDI {
	inst := &ANDI{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "ANDI",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1)},
		Rd:             intPtr(int(t.Rd)),
		ProduceStage:   isa.EX,
		ConsumeStage:   isa.ID,
	}
	return inst
}
