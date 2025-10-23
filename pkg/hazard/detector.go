package hazard

import (
	"fmt"
	"strings"

	"riscv-instruction-encoder/pkg/isa"
)

type HazardType string

const (
	HazardRAW     HazardType = "RAW"
	HazardWAW     HazardType = "WAW"
	HazardWAR     HazardType = "WAR"
	HazardControl HazardType = "CONTROL"
)

type Hazard struct {
	Type         HazardType
	From         int    // index of producer / earlier instruction
	To           int    // index of consumer / later instruction
	Reg          uint8  // register involved (0 if none)
	Description  string // human readable
	StallsNeeded int    // recommended number of NOPs to insert before To
}

// Detector é stateless; cria instância para conveniência/expansão
type Detector struct{}

// Analyze varre a sequência de instruções e identifica hazards.
// - considera forwarding disponível para RAW (exceto load-use entre i e i+1).
// - detecta WAW e WAR como potenciais problemas (sem forwarding automático).
// - detecta instruções de controle (branch/jump) e recomenda stalls.
// Retorna slice de Hazard encontrados (ordenados por From,To).
func (d *Detector) Analyze(instrs []isa.Instruction) []Hazard {
	var hazards []Hazard
	n := len(instrs)

	isLoad := func(s string) bool {
		s = strings.ToUpper(s)
		return strings.Contains(s, "LW") || strings.Contains(s, "LB") || strings.Contains(s, "LD") || strings.Contains(s, "LH")
	}
	isControl := func(s string) bool {
		s = strings.ToUpper(s)
		// branches e saltos comuns
		branchNames := []string{"BEQ", "BNE", "BLT", "BGE", "BLTU", "BGEU", "JAL", "JALR"}
		for _, b := range branchNames {
			if strings.Contains(s, b) {
				return true
			}
		}
		return false
	}

	for i := 0; i < n; i++ {
		a := instrs[i]
		if a == nil {
			continue
		}
		usageA := a.GetRegisterUsage()
		strA := strings.ToUpper(a.String())

		// controle: se instrução i é branch/jump, provavelmente precisa de stalls para instruções subsequentes
		if isControl(strA) {
			// conservador: recomenda 1-2 ciclos dependendo do tipo (JAL/JALR resolvem PC diferente)
			stalls := 1
			if strings.Contains(strA, "JAL") || strings.Contains(strA, "JALR") {
				stalls = 1
			} else {
				stalls = 2
			}
			hazards = append(hazards, Hazard{
				Type:         HazardControl,
				From:         i,
				To:           i + 1,
				Reg:          0,
				Description:  fmt.Sprintf("Control hazard: %s at %d may change PC; recommend %d stall(s)", strA, i, stalls),
				StallsNeeded: stalls,
			})
		}

		// comparar com proximas instruções (janela típica de pipeline)
		for j := i + 1; j < n && j <= i+4; j++ {
			b := instrs[j]
			if b == nil {
				continue
			}
			usageB := b.GetRegisterUsage()

			// RAW: A escreve, B lê
			for _, wa := range usageA.WriteRegs {
				if wa == 0 {
					continue
				}
				for _, rb := range usageB.ReadRegs {
					if wa == rb {
						stalls := 0
						// load-use hazard: próxima instrução usa destino de load -> precisa de 1 stall se j == i+1
						if isLoad(strA) && j == i+1 {
							stalls = 1
						}
						hazards = append(hazards, Hazard{
							Type:         HazardRAW,
							From:         i,
							To:           j,
							Reg:          wa,
							Description:  fmt.Sprintf("RAW: instr %d writes x%d, instr %d reads x%d", i, wa, j, rb),
							StallsNeeded: stalls,
						})
					}
				}
			}

			// WAW: ambos escrevem mesmo registrador
			for _, wa := range usageA.WriteRegs {
				if wa == 0 {
					continue
				}
				for _, wb := range usageB.WriteRegs {
					if wa == wb {
						// conservador: 1 stall recomendado para preservar ordem de escrita em pipelines simples
						hazards = append(hazards, Hazard{
							Type:         HazardWAW,
							From:         i,
							To:           j,
							Reg:          wa,
							Description:  fmt.Sprintf("WAW: instr %d and instr %d both write x%d", i, j, wa),
							StallsNeeded: 1,
						})
					}
				}
			}

			// WAR: A lê, B escreve (pode ocorrer em pipelines com escrita antecipada)
			for _, ra := range usageA.ReadRegs {
				if ra == 0 {
					continue
				}
				for _, wb := range usageB.WriteRegs {
					if ra == wb {
						// conservador: 1 stall recomendado se j < i+3 (quando escrita poderia ocorrer antes da leitura)
						stalls := 0
						if j <= i+2 {
							stalls = 1
						}
						hazards = append(hazards, Hazard{
							Type:         HazardWAR,
							From:         i,
							To:           j,
							Reg:          ra,
							Description:  fmt.Sprintf("WAR: instr %d reads x%d, instr %d writes x%d", i, ra, j, wb),
							StallsNeeded: stalls,
						})
					}
				}
			}
		}
	}

	return hazards
}

// ResolveWithNops insere NOPs (istrução sem efeitos) para satisfazer stalls recomendados.
// Retorna nova slice de instruções e os hazards que foram aplicados.
// Estratégia simples: ao processar instruções em ordem, quando um hazard exige stalls entre i e i+1,
// insere NOPs imediatamente antes da instrução alvo (To).
func (d *Detector) ResolveWithNops(instrs []isa.Instruction) ([]isa.Instruction, []Hazard) {
	hazards := d.Analyze(instrs)
	// construir uma tabela de stalls por destino (index)
	stallsFor := map[int]int{}
	for _, h := range hazards {
		if h.StallsNeeded > 0 {
			// escolher maior necessidade se múltiplos hazards afetarem mesma posição
			cur := stallsFor[h.To]
			if h.StallsNeeded > cur {
				stallsFor[h.To] = h.StallsNeeded
			}
		}
	}
	// gerar nova lista inserindo NOPs
	var out []isa.Instruction
	for i := 0; i < len(instrs); i++ {
		// antes de adicionar instrs[i], verificar se há stalls para esta posição
		if s, ok := stallsFor[i]; ok && s > 0 {
			for k := 0; k < s; k++ {
				out = append(out, &Nop{})
			}
		}
		out = append(out, instrs[i])
	}
	// re-analisar no novo fluxo e retornar hazards detectados após inserção (útil para ver se restaram)
	newHazards := d.Analyze(out)
	return out, newHazards
}

// Nop representa instrução de bolha (no-op)
type Nop struct{}

func (n *Nop) String() string                     { return "NOP" }
func (n *Nop) Decode(inst uint32) isa.Instruction { return n }
func (n *Nop) ExecuteFetchInstruction()           {}
func (n *Nop) ExecuteDecodeInstruction()          {}
func (n *Nop) ExecuteOperation()                  {}
func (n *Nop) ExecuteAccessOperand()              {}
func (n *Nop) ExecuteWriteBack()                  {}
func (n *Nop) GetRegisterUsage() isa.RegisterUsage {
	return isa.RegisterUsage{ReadRegs: []uint8{}, WriteRegs: []uint8{}}
}

// Helpers utilitários (impressão)
func (d *Detector) PrintHazards(hazards []Hazard) {
	for _, h := range hazards {
		fmt.Printf("Hazard %s: from %d to %d reg x%d stalls=%d -- %s\n", h.Type, h.From, h.To, h.Reg, h.StallsNeeded, h.Description)
	}
}
