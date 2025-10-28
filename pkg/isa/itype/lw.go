package itype

import isa "riscv-instruction-encoder/pkg/isa"

type LW struct {
	Type
}

// LW – I-type, load word
func newLW(t Type) *LW {
	inst := &LW{Type: t}
	inst.InstructionMeta = isa.InstructionMeta{
		Name:           "LW",
		OpCode:         uint32(t.OpCode),
		IsLoad:         true,
		IsStore:        false,
		IsBranch:       false,
		IsJump:         false,
		WritesRegister: true,
		ReadsRegister:  true,
		Rs:             []int{int(t.Rs1)},
		Rd:             isa.IntPtr(int(t.Rd)),
		ProduceStage:   isa.MEM, // o dado é produzido no estágio MEM (após acessar memória)
		ConsumeStage:   isa.ID,
	}
	return inst
}
