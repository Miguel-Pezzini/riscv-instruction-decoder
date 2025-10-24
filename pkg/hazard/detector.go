package hazard

import (
	"riscv-instruction-encoder/pkg/isa"
)

type HazardType string

const (
	HazardRAW     HazardType = "RAW"
	HazardWAW     HazardType = "WAW"
	HazardWAR     HazardType = "WAR"
	HazardControl HazardType = "CONTROL"
)

type Hazard struct {
	Type         HazardType
	From         int    // index of producer / earlier instruction
	To           int    // index of consumer / later instruction
	Reg          uint8  // register involved (0 if none)
	Description  string // human readable
	StallsNeeded int    // recommended number of NOPs to insert before To
}

// Nop representa instrução de bolha (no-op)
type Nop struct{}

func (n *Nop) String() string                     { return "NOP" }
func (n *Nop) Decode(inst uint32) isa.Instruction { return n }
func (n *Nop) GetMeta() isa.InstructionMeta       { return isa.InstructionMeta{} }
func (n *Nop) ExecuteFetchInstruction()           {}
func (n *Nop) ExecuteDecodeInstruction()          {}
func (n *Nop) ExecuteOperation()                  {}
func (n *Nop) ExecuteAccessOperand()              {}
func (n *Nop) ExecuteWriteBack()                  {}
func (n *Nop) GetRegisterUsage() isa.RegisterUsage {
	return isa.RegisterUsage{ReadRegs: []uint8{}, WriteRegs: []uint8{}}
}
