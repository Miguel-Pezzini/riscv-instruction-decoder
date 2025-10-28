package runner

import (
	"fmt"
	"os"
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

func (p *Pipeline) writeFile() {
	file, err := os.Create(p.file_path)
	if err != nil {
		fmt.Printf("Error to create file %s: %v\n", p.file_path, err)
		return
	}
	defer file.Close()
	_, _ = file.WriteString("PC\tInstruction\n")
	_, _ = file.WriteString("===============================\n")
	for _, instr := range p.Instructions {
		line := fmt.Sprintf("0x%08X\t%s\n", instr.PC, instr.Instruction.String())
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Printf("Error to write in file %s: %v\n", p.file_path, err)
			return
		}
	}
}

func (p *Pipeline) printResult() {
	countNop := 0
	for _, instruction := range p.Instructions {
		if instruction.Instruction.GetMeta().Name == "NOP" {
			countNop++
		}
	}

	origCount := len(p.Instructions) - countNop
	totalCount := len(p.Instructions)
	overhead := float64(totalCount-origCount) / float64(origCount) * 100

	fmt.Printf("\nInput: fib_rec_binario.txt (%d instruções)\n", origCount)
	fmt.Println("Model pipeline: IF ID EX MEM WB")
	fmt.Println()

	var mode string
	if p.data_hazard && p.control_hazard {
		mode = "-- INTEGRATED"
	} else if p.data_hazard {
		mode = "-- DATA"
	} else if p.control_hazard {
		mode = "-- CONTROL"
	} else {
		mode = "-- NO CONTROL"
	}

	forwardingText := "sem forwarding"
	if p.forwarding {
		forwardingText = "com forwarding"
	}

	fmt.Printf("%s (%s)\n", mode, forwardingText)
	fmt.Printf("Output: %s\n", p.file_path)
	fmt.Printf("Instruções originais: %d\n", origCount)
	fmt.Printf("Instruções finais: %d\n", totalCount)
	fmt.Printf("NOPs inseridos: %d\n", countNop)
	fmt.Printf("Sobreacusto: +%.1f%%\n", overhead)
	fmt.Println("========================================")
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
