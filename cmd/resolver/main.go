package main

import (
	"riscv-instruction-encoder/pkg/decoder"
	"riscv-instruction-encoder/pkg/runner"
)

const (
	FORMAT_BIN = "bin"
	FORMAT_HEX = "hex"
)

const (
	BIN_INSTRUCTION_FILE_NAME = "../../testdata/fib_rec_binario.txt"
	HEX_INSTRUCTION_FILE_NAME = "../../testdata/fib_rec_hexadecimal.txt"
)

func main() {
	instructionsFromBinaryFile := decoder.DecodeFromFile(BIN_INSTRUCTION_FILE_NAME, FORMAT_BIN)
	// instructionsFromHexFile := DecodeFromFile(HEX_INSTRUCTION_FILE_NAME, FORMAT_HEX)
	// DecodeInstructionFromUInt32(instructionsFromHexFile)

	executions := []struct {
		forwarding           bool
		dataHazardControl    bool
		controlHazardControl bool
		fileName             string
	}{
		{false, true, false, "../../pkg/files/output_data_no_forwarding.txt"},
		{true, true, false, "../../pkg/files/output_data_forwarding.txt"},
		{false, false, true, "../../pkg/files/output_control_no_forwarding.txt"},
		{true, false, true, "../../pkg/files/output_control_forwarding.txt"},
		{false, true, true, "../../pkg/files/output_integrated_no_forwarding.txt"},
		{true, true, true, "../../pkg/files/output_integrated_forwarding.txt"},
	}

	decodedInstructions := decoder.DecodeInstructionFromUInt32(instructionsFromBinaryFile)
	for _, exec := range executions {
		runner.Run(
			decodedInstructions,
			exec.forwarding,
			exec.dataHazardControl,
			exec.controlHazardControl,
			exec.fileName,
		)
	}
}
