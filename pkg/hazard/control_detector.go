package hazard

import "riscv-instruction-encoder/pkg/isa"

func HasControlHazard(currentInstruction isa.PipelineInstruction, executing []*isa.PipelineInstruction, forwarding bool) bool {
	for _, prev := range executing {
		if hasUnresolvedBranchHazard(currentInstruction, *prev) {
			return true
		}
	}

	return false
}

func hasUnresolvedBranchHazard(currentInstruction isa.PipelineInstruction, previousInstruction isa.PipelineInstruction) bool {
	prevMeta := previousInstruction.Instruction.GetMeta()
	if (prevMeta.IsBranch || prevMeta.IsJump) && !previousInstruction.HasCompleted && previousInstruction.CurrentStage < int(isa.WB) {
		return true
	}

	return false
}
