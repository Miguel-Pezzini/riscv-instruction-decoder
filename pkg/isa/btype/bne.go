package btype

import isa "riscv-instruction-encoder/pkg/isa"

type BNE struct {
	Type
}

func newBNE(t Type) *BNE {
	inst := &BNE{Type: t}

	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "BNE",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       true,
		IsJump:         false,
		WritesRegister: false,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1), int(t.Rs2)},
		Rd:             nil,
		ProduceStage:   isa.EX, // decisão do branch é feita no estágio EX
		ConsumeStage:   isa.ID, // registradores lidos no estágio ID
	}

	return inst
}
