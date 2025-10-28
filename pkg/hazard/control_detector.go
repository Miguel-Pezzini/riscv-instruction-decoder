package hazard

import "riscv-instruction-encoder/pkg/isa"

func HasControlHazard(currentInstruction isa.PipelineInstruction, executing []*isa.PipelineInstruction, forwarding bool) bool {
	for _, prev := range executing {
		if hasBranchRAWHazard(currentInstruction, *prev, forwarding) || hasUnresolvedBranchHazard(currentInstruction, *prev) {
			return true
		}
	}

	return false
}

func hasBranchRAWHazard(currentInstruction isa.PipelineInstruction, previousInstruction isa.PipelineInstruction, forwarding bool) bool {
	currMeta := currentInstruction.Instruction.GetMeta()

	if currMeta.IsBranch {
		return false
	}

	prevMeta := previousInstruction.Instruction.GetMeta()

	if prevMeta.Rd == nil || previousInstruction.HasCompleted {
		return false
	}

	for _, rs := range currMeta.Rs {
		if rs == *prevMeta.Rd {
			cyclesToConsume := int(currMeta.ConsumeStage) - currentInstruction.CurrentStage
			if !forwarding {
				cyclesToProduce := int(isa.WB) - previousInstruction.CurrentStage
				if cyclesToProduce >= 0 && cyclesToConsume >= 0 && cyclesToProduce > cyclesToConsume {
					return true
				}
			} else {
				cyclesToProduce := int(prevMeta.ProduceStage) - previousInstruction.CurrentStage
				if cyclesToProduce >= 0 && cyclesToConsume >= 0 && cyclesToProduce > cyclesToConsume {
					return true
				}
			}
		}
	}
	return false
}

func hasUnresolvedBranchHazard(currentInstruction isa.PipelineInstruction, previousInstruction isa.PipelineInstruction) bool {
	prevMeta := previousInstruction.Instruction.GetMeta()
	if (prevMeta.IsBranch || prevMeta.IsJump) && !previousInstruction.HasCompleted && previousInstruction.CurrentStage < int(isa.EX) {
		return true
	}

	return false
}
