package runner

import (
	"riscv-instruction-encoder/pkg/isa"
)

type Pipeline struct {
	CurrentCycle          int
	Instructions          []isa.PipelineInstruction
	NumStages             int
	executingInstructions []int
}

func InstructionsToPipeline(instructions []isa.Instruction) []isa.PipelineInstruction {
	pipelineInstructions := make([]isa.PipelineInstruction, len(instructions))
	for i, instr := range instructions {
		pipelineInstructions[i] = isa.PipelineInstruction{
			Instruction:  instr,
			CurrentStage: 0,
			HasCompleted: false,
			HasStarted:   false,
			Id:           i + 1,
		}
	}
	return pipelineInstructions
}

func NewPipeline(instructions []isa.Instruction) *Pipeline {
	stages := len(isa.Stages)

	return &Pipeline{
		CurrentCycle: 0,
		Instructions: InstructionsToPipeline(instructions),
		NumStages:    stages,
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

func (p *Pipeline) getNextInstruction() *isa.PipelineInstruction {
	for i := range p.Instructions {
		instr := &p.Instructions[i]
		if !instr.HasStarted && !instr.HasCompleted {
			return instr
		}
	}
	return nil
}

func (p *Pipeline) Step() {
	for _, idx := range p.executingInstructions {
		instr := &p.Instructions[idx]
		instr.CurrentStage++

		if instr.CurrentStage >= p.NumStages {
			instr.HasCompleted = true
		}
	}

	active := make([]int, 0)
	for _, idx := range p.executingInstructions {
		if !p.Instructions[idx].HasCompleted {
			active = append(active, idx)
		}
	}
	p.executingInstructions = active

}

func Run(instructions []isa.Instruction) {
	p := NewPipeline(instructions)

	for !p.hasCompleted() {
		p.CurrentCycle++
		p.Step()
	}
}
