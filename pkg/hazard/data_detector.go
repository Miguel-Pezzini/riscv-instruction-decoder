package hazard

import "riscv-instruction-encoder/pkg/isa"

func HasDataHazard(currentInstruction isa.PipelineInstruction, executing []*isa.PipelineInstruction, forwarding bool) bool {
	for _, prev := range executing {
		if isRAWHazard(currentInstruction, *prev, forwarding) || isWARHazard(*prev, currentInstruction, forwarding) {
			return true
		}
	}
	return false
}

// Read after Write Hazard detection
func isRAWHazard(currentInstruction isa.PipelineInstruction, previousInstruction isa.PipelineInstruction, forwarding bool) bool {
	currMeta := currentInstruction.Instruction.GetMeta()
	prevMeta := previousInstruction.Instruction.GetMeta()

	if !prevMeta.WritesRegister || prevMeta.Rd == nil || !previousInstruction.HasStarted || previousInstruction.HasCompleted {
		return false
	}

	for _, rs := range currMeta.Rs {
		if rs == *prevMeta.Rd {
			cyclesToConsume := int(currMeta.ConsumeStage) - currentInstruction.CurrentStage
			var cyclesToProduce int
			if !forwarding {
				cyclesToProduce = int(isa.WB) - previousInstruction.CurrentStage
			} else {
				cyclesToProduce = int(prevMeta.ProduceStage) - previousInstruction.CurrentStage
			}
			if cyclesToProduce >= 0 && cyclesToConsume >= 0 && cyclesToProduce > cyclesToConsume {
				return true
			}
		}
	}

	return false
}

// Write after Read Hazard detection
func isWARHazard(prevInstruction, currInstruction isa.PipelineInstruction, forwarding bool) bool {
	prevMeta := prevInstruction.Instruction.GetMeta()
	currMeta := currInstruction.Instruction.GetMeta()

	if !prevMeta.ReadsRegister || !currMeta.WritesRegister || currMeta.Rd == nil || !prevInstruction.HasStarted || prevInstruction.HasCompleted {
		return false
	}

	for _, rs := range prevMeta.Rs {
		if rs == *currMeta.Rd {
			var cyclesToWrite int
			cyclesToRead := int(prevMeta.ConsumeStage) - prevInstruction.CurrentStage
			if !forwarding {
				cyclesToWrite = int(isa.WB) - currInstruction.CurrentStage
			} else {
				cyclesToWrite = int(currMeta.ProduceStage) - currInstruction.CurrentStage
			}
			if cyclesToWrite >= 0 && cyclesToRead >= 0 && cyclesToRead < cyclesToWrite {
				return true
			}
		}
	}

	return false
}
