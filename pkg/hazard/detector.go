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

func isRAWHazard(currentInstruction isa.PipelineInstruction, previousInstruction isa.PipelineInstruction, forwarding bool) bool {
	currMeta := currentInstruction.Instruction.GetMeta()
	prevMeta := previousInstruction.Instruction.GetMeta()

	if !prevMeta.WritesRegister || prevMeta.Rd == nil || !previousInstruction.HasStarted || previousInstruction.HasCompleted {
		return false
	}

	for _, rs := range currMeta.Rs {
		if rs == *prevMeta.Rd {
			if !forwarding {
				cyclesToProduce := int(isa.WB) - previousInstruction.CurrentStage
				cyclesToConsume := int(currMeta.ConsumeStage) - currentInstruction.CurrentStage
				if cyclesToProduce >= 0 && cyclesToConsume >= 0 && cyclesToProduce <= cyclesToConsume {
					return true
				}
			} else {
				cyclesToProduce := int(prevMeta.ProduceStage) - previousInstruction.CurrentStage
				cyclesToConsume := int(currMeta.ConsumeStage) - currentInstruction.CurrentStage
				if cyclesToProduce >= 0 && cyclesToConsume >= 0 && cyclesToProduce < cyclesToConsume {
					return true
				}
			}
		}
	}

	return false
}

func isWARHazard(prevInstruction, currInstruction isa.PipelineInstruction, forwarding bool) bool {
	prevMeta := prevInstruction.Instruction.GetMeta()
	currMeta := currInstruction.Instruction.GetMeta()

	if !prevMeta.ReadsRegister || !currMeta.WritesRegister || currMeta.Rd == nil || !prevInstruction.HasStarted || prevInstruction.HasCompleted {
		return false
	}

	for _, rs := range prevMeta.Rs {
		if rs == *currMeta.Rd {
			var cyclesToWrite, cyclesToRead int

			if !forwarding {
				cyclesToWrite = int(isa.WB) - currInstruction.CurrentStage
				cyclesToRead = int(prevMeta.ConsumeStage) - prevInstruction.CurrentStage
				if cyclesToWrite >= 0 && cyclesToRead >= 0 && cyclesToWrite <= cyclesToRead {
					return true
				}
			} else {
				cyclesToWrite = int(currMeta.ProduceStage) - currInstruction.CurrentStage
				cyclesToRead = int(prevMeta.ConsumeStage) - prevInstruction.CurrentStage
				if cyclesToWrite >= 0 && cyclesToRead >= 0 && cyclesToWrite < cyclesToRead {
					return true
				}
			}
		}
	}

	return false
}
