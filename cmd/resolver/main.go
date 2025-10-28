package main

import (
	"fmt"
	"os"
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
	var formatChoice string
	fmt.Println("Select instruction format to decode (bin / hex):")
	_, err := fmt.Scanln(&formatChoice)
	if err != nil {
		fmt.Println("Invalid input. Defaulting to hex format.")
		os.Exit(1)
	}

	var format string
	var fileName string

	switch formatChoice {
	case "bin", "BIN":
		format = FORMAT_BIN
		fileName = BIN_INSTRUCTION_FILE_NAME
	case "hex", "HEX":
		format = FORMAT_HEX
		fileName = HEX_INSTRUCTION_FILE_NAME
	default:
		fmt.Println("Invalid format choice. Please select 'bin' or 'hex'.")
		os.Exit(1)
	}

	encodedInstructions := decoder.DecodeFromFile(fileName, format)

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

	decodedInstructions := decoder.DecodeInstructionFromUInt32(encodedInstructions)
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
