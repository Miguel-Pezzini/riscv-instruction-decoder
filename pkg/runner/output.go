package runner

import (
	"fmt"
	"os"
)

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
		line := fmt.Sprintf("0x%08X\t%s\n", instr.OriginalPC, instr.Instruction.String())
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
