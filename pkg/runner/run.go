package runner

import (
	"fmt"
	"riscv-instruction-encoder/pkg/hazard"
	"riscv-instruction-encoder/pkg/isa"
)

type Pipeline struct {
	CurrentCycle          int
	Instructions          []*isa.PipelineInstruction
	NumStages             int
	executingInstructions []*isa.PipelineInstruction
	forwarding            bool
}

func InstructionsToPipeline(instructions []isa.Instruction) []*isa.PipelineInstruction {
	pipelineInstructions := make([]*isa.PipelineInstruction, len(instructions))
	for i, instr := range instructions {
		pipelineInstructions[i] = &isa.PipelineInstruction{
			Instruction:  instr,
			CurrentStage: 0,
			HasCompleted: false,
			HasStarted:   false,
			Id:           i + 1,
		}
	}
	return pipelineInstructions
}

func NewPipeline(instructions []isa.Instruction, forwarding bool) *Pipeline {
	stages := len(isa.Stages)

	return &Pipeline{
		CurrentCycle: 0,
		Instructions: InstructionsToPipeline(instructions),
		NumStages:    stages,
		forwarding:   forwarding,
	}
}

func (p *Pipeline) hasCompleted() bool {
	for _, instr := range p.Instructions {
		if !instr.HasCompleted {
			return false
		}
	}
	return true
}

func (p *Pipeline) getNextInstruction() (*isa.PipelineInstruction, int) {
	for i, instruction := range p.Instructions {
		if !instruction.HasStarted && !instruction.HasCompleted {
			return instruction, i
		}
	}
	return nil, -1
}

func createNOP() *isa.PipelineInstruction {
	return &isa.PipelineInstruction{
		Instruction:  isa.NewNOP(),
		CurrentStage: 1,
		HasStarted:   true,
		HasCompleted: false,
		Id:           -1,
	}
}

func (p *Pipeline) Step() {
	for _, instruction := range p.executingInstructions {
		instruction.CurrentStage++

		if instruction.CurrentStage >= p.NumStages {
			instruction.HasCompleted = true
		}
	}

	nextInstruction, index := p.getNextInstruction()
	if nextInstruction != nil {
		if hazard.HasDataHazard(*nextInstruction, p.executingInstructions, p.forwarding) || hazard.HasControlHazard(*nextInstruction, p.executingInstructions, p.forwarding) {
			nop := createNOP()
			p.executingInstructions = append(p.executingInstructions, nop)

			p.Instructions = append(
				p.Instructions[:index],
				append([]*isa.PipelineInstruction{nop}, p.Instructions[index:]...)...,
			)
		} else {
			nextInstruction.HasStarted = true
			nextInstruction.CurrentStage = int(isa.IF)
			p.executingInstructions = append(p.executingInstructions, nextInstruction)
		}
	}

	for _, instruction := range p.executingInstructions {
		isa.ExecuteStage(isa.Stage(instruction.CurrentStage), instruction.Instruction)
	}
	fmt.Print("\n")

	active := make([]*isa.PipelineInstruction, 0)
	for _, instruction := range p.executingInstructions {
		if !instruction.HasCompleted {
			active = append(active, instruction)
		}
	}
	p.executingInstructions = active
}

func Run(instructions []isa.Instruction, forwarding bool) {
	p := NewPipeline(instructions, forwarding)

	for !p.hasCompleted() {
		p.CurrentCycle++
		p.Step()
	}
}
