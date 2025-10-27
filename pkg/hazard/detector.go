package hazard

import "riscv-instruction-encoder/pkg/isa"

func HasDataHazard(currentInstruction isa.PipelineInstruction, executing []isa.PipelineInstruction, forwarding bool) bool {
	currMeta := currentInstruction.Instruction.GetMeta()

	for _, prev := range executing {
		if !prev.HasStarted || prev.HasCompleted {
			continue
		}

		prevMeta := prev.Instruction.GetMeta()
		if !prevMeta.WritesRegister || prevMeta.Rd == nil {
			continue
		}

		for _, rs := range currMeta.Rs {
			if rs == *prevMeta.Rd {
				if !forwarding {
					cyclesToProduce := int(isa.WB) - prev.CurrentStage
					cyclesToConsume := int(currMeta.ConsumeStage) - currentInstruction.CurrentStage

					if cyclesToProduce >= 0 && cyclesToConsume >= 0 && cyclesToProduce <= cyclesToConsume {
						return true
					}
					return true
				}

				cyclesToProduce := int(prevMeta.ProduceStage) - prev.CurrentStage
				cyclesToConsume := int(currMeta.ConsumeStage) - currentInstruction.CurrentStage

				if cyclesToProduce >= 0 && cyclesToConsume >= 0 && cyclesToProduce < cyclesToConsume {
					return true
				}

			}
		}
	}
	return false
}
