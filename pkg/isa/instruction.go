package isa

import (
	"fmt"
)

type Stage string

const (
	IF  Stage = "IF"
	ID  Stage = "ID"
	EX  Stage = "EX"
	MEM Stage = "MEM"
	WB  Stage = "WB"
)

var Stages = []Stage{IF, ID, EX, MEM, WB}

type RegisterUsage struct {
	ReadRegs  []uint8
	WriteRegs []uint8
}

type InstructionMeta struct {
	Name           string
	OpCode         uint32
	IsLoad         bool
	IsStore        bool
	IsBranch       bool
	WritesRegister bool

	Rs []int
	Rd *int

	ProduceStage Stage
	ConsumeStage Stage
}

type Instruction interface {
	String() string
	Decode(inst uint32) Instruction
	ExecuteFetchInstruction()
	ExecuteDecodeInstruction()
	ExecuteOperation()
	ExecuteAccessOperand()
	ExecuteWriteBack()
	GetMeta() InstructionMeta
}

type BaseInstruction struct {
	InstructionMeta InstructionMeta
}

func (b *BaseInstruction) GetMeta() InstructionMeta {
	return b.InstructionMeta
}

func (b *BaseInstruction) SetMeta(i InstructionMeta) {
	b.InstructionMeta = i
}

func (b *BaseInstruction) ExecuteFetchInstruction() {}

func (b *BaseInstruction) ExecuteDecodeInstruction() {}

func (b *BaseInstruction) ExecuteOperation() {}

func (b *BaseInstruction) ExecuteAccessOperand() {}

func (b *BaseInstruction) ExecuteWriteBack() {}

type RawInstruction struct {
	Origin string
	Value  uint32
}

func ExecuteStage(stage Stage, instruction Instruction) {
	switch stage {
	case IF:
		instruction.ExecuteFetchInstruction()
	case ID:
		instruction.ExecuteDecodeInstruction()
	case EX:
		instruction.ExecuteOperation()
	case MEM:
		instruction.ExecuteAccessOperand()
	case WB:
		instruction.ExecuteWriteBack()
	default:
		fmt.Printf("Stage not defined")
	}
}
