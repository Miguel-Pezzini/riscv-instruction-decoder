package main

import (
	"riscv-instruction-encoder/decoder"
	isa "riscv-instruction-encoder/instructions"
)

const (
	FORMAT_BIN = "bin"
	FORMAT_HEX = "hex"
)

const (
	BIN_INSTRUCTION_FILE_NAME = "fib_rec_binario.txt"
	HEX_INSTRUCTION_FILE_NAME = "fib_rec_hexadecimal.txt"
)

func main() {
	instructionsFromBinaryFile := decoder.DecodeFromFile(BIN_INSTRUCTION_FILE_NAME, FORMAT_BIN)
	// instructionsFromHexFile := DecodeFromFile(HEX_INSTRUCTION_FILE_NAME, FORMAT_HEX)
	// DecodeInstructionFromUInt32(instructionsFromHexFile)

	decodedInstructions := decoder.DecodeInstructionFromUInt32(instructionsFromBinaryFile)
	isa.Execute(decodedInstructions)
}
