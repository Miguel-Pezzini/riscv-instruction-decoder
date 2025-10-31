package runner

import (
	"riscv-instruction-encoder/pkg/hazard"
	"riscv-instruction-encoder/pkg/isa"
)

type Pipeline struct {
	CurrentCycle          int
	Instructions          []*isa.PipelineInstruction
	NumStages             int
	executingInstructions []*isa.PipelineInstruction
	forwarding            bool
	data_hazard           bool
	control_hazard        bool
	file_path             string
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
			PC:           i * 4,
			OriginalPC:   i * 4,
		}
	}
	return pipelineInstructions
}

func NewPipeline(instructions []isa.Instruction, forwarding bool, data_hazard bool, control_hazard bool, file_path string) *Pipeline {
	stages := len(isa.Stages)

	return &Pipeline{
		CurrentCycle:   0,
		Instructions:   InstructionsToPipeline(instructions),
		NumStages:      stages,
		forwarding:     forwarding,
		data_hazard:    data_hazard,
		control_hazard: control_hazard,
		file_path:      file_path,
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

func (p *Pipeline) insertNOPAt(index int) {
	nop := createNOP()
	if index < len(p.Instructions) {
		nop.PC = p.Instructions[index].PC
	} else {
		nop.PC = len(p.Instructions) * 4
	}

	p.Instructions = append(
		p.Instructions[:index],
		append([]*isa.PipelineInstruction{nop}, p.Instructions[index:]...)...,
	)
	p.executingInstructions = append(p.executingInstructions, nop)
	for i := index + 1; i < len(p.Instructions); i++ {
		p.Instructions[i].PC += 4
	}
}

func (p *Pipeline) insertInstruction(instruction *isa.PipelineInstruction) {
	instruction.HasStarted = true
	instruction.CurrentStage = int(isa.IF)
	p.executingInstructions = append(p.executingInstructions, instruction)
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
		nextInstruction.CurrentStage = int(isa.IF)
		if (hazard.HasDataHazard(*nextInstruction, p.executingInstructions, p.forwarding) && p.data_hazard) || (hazard.HasControlHazard(*nextInstruction, p.executingInstructions, p.forwarding) && p.control_hazard) {
			p.insertNOPAt(index)
		} else {
			p.insertInstruction(nextInstruction)
		}
	}

	// for _, instruction := range p.executingInstructions {
	// 	fmt.Print(" - PC: ", instruction.PC, " | ")
	// 	isa.ExecuteStage(isa.Stage(instruction.CurrentStage), instruction.Instruction)
	// }
	// fmt.Print("\n")

	active := make([]*isa.PipelineInstruction, 0)
	for _, instruction := range p.executingInstructions {
		if !instruction.HasCompleted {
			active = append(active, instruction)
		}
	}
	p.executingInstructions = active
}

func Run(instructions []isa.Instruction, forwarding bool, data_hazard bool, control_hazard bool, file_path string) {
	p := NewPipeline(instructions, forwarding, data_hazard, control_hazard, file_path)

	for !p.hasCompleted() {
		p.CurrentCycle++
		p.Step()
	}
	p.printResult()
	p.writeFile()
}
