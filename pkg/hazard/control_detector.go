package hazard

import "riscv-instruction-encoder/pkg/isa"

func HasControlHazard(currentInstruction isa.PipelineInstruction, executing []*isa.PipelineInstruction, forwarding bool) bool {
	meta := currentInstruction.Instruction.GetMeta()

	if !meta.IsBranch && !meta.IsJump {
		return false
	}

	if meta.IsBranch && hasBranchRAWHazard(currentInstruction, executing, forwarding) {
		return true
	}

	if hasUnresolvedBranchHazard(currentInstruction, executing) {
		return true
	}

	return false
}

func hasBranchRAWHazard(currentInstruction isa.PipelineInstruction, executing []*isa.PipelineInstruction, forwarding bool) bool {
	currMeta := currentInstruction.Instruction.GetMeta()

	for _, prev := range executing {
		prevMeta := prev.Instruction.GetMeta()

		if !prevMeta.WritesRegister || prevMeta.Rd == nil || prev.HasCompleted {
			continue
		}

		for _, rs := range currMeta.Rs {
			if rs == *prevMeta.Rd {
				cyclesToConsume := int(currMeta.ConsumeStage) - currentInstruction.CurrentStage
				if !forwarding {
					cyclesToProduce := int(isa.WB) - prev.CurrentStage
					if cyclesToProduce >= 0 && cyclesToConsume >= 0 && cyclesToProduce > cyclesToConsume {
						return true
					}
				} else {
					cyclesToProduce := int(prevMeta.ProduceStage) - prev.CurrentStage
					if cyclesToProduce >= 0 && cyclesToConsume >= 0 && cyclesToProduce > cyclesToConsume {
						return true
					}
				}
			}
		}
	}
	return false
}

func hasUnresolvedBranchHazard(currentInstruction isa.PipelineInstruction, executing []*isa.PipelineInstruction) bool {
	for _, prev := range executing {
		prevMeta := prev.Instruction.GetMeta()
		if (prevMeta.IsBranch || prevMeta.IsJump) && !prev.HasCompleted && prev.CurrentStage < int(isa.EX) {
			return true
		}
	}
	return false
}
