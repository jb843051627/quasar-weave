package quality

import (
	"fmt"
	"sort"

	"github.com/jb843051627/quasar-weave/internal/model"
)

type RuleSet struct {
	gates map[string]model.QualityGate
}

func NewRuleSet(gates []model.QualityGate) *RuleSet {
	result := &RuleSet{gates: make(map[string]model.QualityGate, len(gates))}
	for _, gate := range gates {
		result.gates[gate.ID] = gate
	}
	return result
}

func (r *RuleSet) Get(id string) (model.QualityGate, error) {
	gate, ok := r.gates[id]
	if !ok {
		return model.QualityGate{}, fmt.Errorf("quality gate %q: %w", id, model.ErrNotFound)
	}
	if !gate.Active {
		return model.QualityGate{}, fmt.Errorf("quality gate %q is inactive", id)
	}
	return gate, nil
}

func (r *RuleSet) Replace(gates []model.QualityGate) {
	updated := make(map[string]model.QualityGate, len(gates))
	for _, gate := range gates {
		updated[gate.ID] = gate
	}
	r.gates = updated
}

func (r *RuleSet) IDs() []string {
	ids := make([]string, 0, len(r.gates))
	for id := range r.gates {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
