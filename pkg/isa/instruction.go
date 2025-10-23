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

type Instruction interface {
	String() string
	Decode(inst uint32) Instruction
	ExecuteFetchInstruction()
	ExecuteDecodeInstruction()
	ExecuteOperation()
	ExecuteAccessOperand()
	ExecuteWriteBack()
	GetRegisterUsage() RegisterUsage
}

type BaseInstruction struct{}

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
