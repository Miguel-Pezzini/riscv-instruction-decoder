package RTypeInstructions

type AddInstruction struct {
	Type
}

func (instruction *AddInstruction) ExecuteFetchInstruction() {
	print(instruction.String())
}
