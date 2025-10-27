package ITypeInstructions

import isa "riscv-instruction-encoder/pkg/isa"

type ADDI struct {
	Type
}

func NewADDI(t Type) *ADDI {
	inst := &ADDI{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "ADDI",
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
