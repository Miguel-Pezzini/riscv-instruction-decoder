package main

import (
	"fmt"
	"riscv-instruction-encoder/pkg/decoder"
	"riscv-instruction-encoder/pkg/isa"
)

const (
	FORMAT_BIN = "bin"
	FORMAT_HEX = "hex"
)

const (
	BIN_INSTRUCTION_FILE_NAME = "fib_rec_binario.txt"
	HEX_INSTRUCTION_FILE_NAME = "fib_rec_hexadecimal.txt"
)

func execute(instructions []isa.Instruction) {
	numStages := len(isa.Stages)
	numInstr := len(instructions)
	totalCycles := numInstr + numStages - 1
	for cycle := 1; cycle <= totalCycles; cycle++ {
		fmt.Printf("Ciclo %d:\n", cycle)

		for i := 0; i < numInstr; i++ {
			stageIndex := cycle - i - 1
			if stageIndex >= 0 && stageIndex < numStages {
				isa.ExecuteStage(isa.Stages[stageIndex], instructions[i])
			}
		}
		fmt.Println()
	}
}

func main() {
	instructionsFromBinaryFile := decoder.DecodeFromFile(BIN_INSTRUCTION_FILE_NAME, FORMAT_BIN)
	// instructionsFromHexFile := DecodeFromFile(HEX_INSTRUCTION_FILE_NAME, FORMAT_HEX)
	// DecodeInstructionFromUInt32(instructionsFromHexFile)

	decodedInstructions := decoder.DecodeInstructionFromUInt32(instructionsFromBinaryFile)
	execute(decodedInstructions)
}
