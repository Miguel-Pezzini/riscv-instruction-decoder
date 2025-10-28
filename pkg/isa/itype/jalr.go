package itype

import isa "riscv-instruction-encoder/pkg/isa"

type JALR struct {
	Type
}

func newJALR(t Type) *JALR {
	inst := &JALR{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "JALR",
		OpCode:         uint32(t.OpCode),
		IsLoad:         false,
		IsStore:        false,
		IsBranch:       true, // é um salto
		IsJump:         true,
		WritesRegister: true, // grava o PC+4 em Rd
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1)},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.EX, // calcula destino no EX
		ConsumeStage:   isa.ID,
	}
	return inst
}
