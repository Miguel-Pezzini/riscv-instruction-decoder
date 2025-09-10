package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Instruction interface {
	Encode() uint32
	String() string
}

type RType struct {
	Opcode uint8 // 7 bits
	Rd     uint8 // 5 bits
	Funct3 uint8 // 3 bits
	Rs1    uint8 // 5 bits
	Rs2    uint8 // 5 bits
	Funct7 uint8 // 7 bits
}

type IType struct {
	OpCode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

type SType struct {
	Opcode uint8  // 7 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Rs2    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

type BType struct {
	Opcode uint8  // 7 bits
	Funct3 uint8  // 3 bits
	Rs1    uint8  // 5 bits
	Rs2    uint8  // 5 bits
	Imm    uint16 // 12 bits
}

type UType struct {
	Opcode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Imm    uint32 // 20 bits
}

type JType struct {
	Opcode uint8  // 7 bits
	Rd     uint8  // 5 bits
	Imm    uint32 // 20 bits
}

func DecodeRType(inst uint32) RType {
	r := RType{
		Opcode: uint8(inst & 0x7F),         // bits 0-6
		Rd:     uint8((inst >> 7) & 0x1F),  // bits 7-11
		Funct3: uint8((inst >> 12) & 0x7),  // bits 12-14
		Rs1:    uint8((inst >> 15) & 0x1F), // bits 15-19
		Rs2:    uint8((inst >> 20) & 0x1F), // bits 20-24
		Funct7: uint8((inst >> 25) & 0x7F), // bits 25-31
	}
	return r
}

func DecodeIType(inst uint32) IType {
	i := IType{
		OpCode: uint8(inst & 0x7F),           // bits 0-6
		Rd:     uint8((inst >> 7) & 0x1F),    // bits 7-11
		Funct3: uint8((inst >> 12) & 0x7),    // bits 12-14
		Rs1:    uint8((inst >> 15) & 0x1F),   // bits 15-19
		Imm:    uint16((inst >> 20) & 0xFFF), // bits 20-31
	}
	return i
}

func DecodeInstructionFromUInt32(encodedInstructions []uint32) {
	for _, v := range encodedInstructions {
		checkOpCode := uint8(v & 0x7F)

		switch checkOpCode {
		case 0x33: // R-type
			r := DecodeRType(v)
			fmt.Printf("R-Type: %+v\n", r)

		case 0x13, 0x03, 0x67, 0x73: // I-type
			i := DecodeIType(v)
			fmt.Printf("I-Type: %+v\n", i)

		case 0x23: // S-type
			// s := DecodeSType(v)
			fmt.Printf("S-Type opcode %02X ainda não implementado\n", checkOpCode)

		case 0x63: // B-type
			// b := DecodeBType(v)
			fmt.Printf("B-Type opcode %02X ainda não implementado\n", checkOpCode)

		case 0x37, 0x17: // U-type
			// u := DecodeUType(v)
			fmt.Printf("U-Type opcode %02X ainda não implementado\n", checkOpCode)

		case 0x6F: // J-type
			// j := DecodeJType(v)
			fmt.Printf("J-Type opcode %02X ainda não implementado\n", checkOpCode)

		default:
			fmt.Printf("Opcode %02X não reconhecido\n", checkOpCode)
		}
	}
}

func DecodeFromBinaryFile() []uint32 {
	var instructions []uint32
	file, err := os.Open("fib_rec_binario.txt")
	if err != nil {
		log.Fatalf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		num, err := strconv.ParseUint(row, 2, 32)
		if err != nil {
			panic(err)
		}
		instructions = append(instructions, uint32(num))
	}

	return instructions
}

func DecodeFromHexFile() []uint32 {
	var instructions []uint32
	file, err := os.Open("fib_rec_hexadecimal.txt")
	if err != nil {
		log.Fatalf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		num, err := strconv.ParseUint(row, 16, 32)
		if err != nil {
			panic(err)
		}
		instructions = append(instructions, uint32(num))
	}

	return instructions
}

func main() {
	var instructionsFromBinaryFile []uint32 = DecodeFromBinaryFile()
	var instructionsFromHexFile []uint32 = DecodeFromHexFile()

	// for i, v := range instructionsFromBinaryFile {
	// 	fmt.Printf("nums[%d] = %d\n", i, v)
	// }

	// for i, v := range instructionsFromHexFile {
	// 	fmt.Printf("nums[%d] = %d\n", i, v)
	// }

	DecodeInstructionFromUInt32(instructionsFromBinaryFile)
	DecodeInstructionFromUInt32(instructionsFromHexFile)

}
