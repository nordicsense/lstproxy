package ts

import (
	"sort"
	"time"
)

type MultiColumnTimeseries struct {
	Times  []time.Time
	Values map[string][]float64
}

func (cmtt MultiColumnTimeseries) Sort() {
	sort.Sort(cmtt)
}

func (cmtt MultiColumnTimeseries) Len() int { return len(cmtt.Times) }
func (cmtt MultiColumnTimeseries) Swap(i, j int) {
	cmtt.Times[i], cmtt.Times[j] = cmtt.Times[j], cmtt.Times[i]
	for _, v := range cmtt.Values {
		v[i], v[j] = v[j], v[i]
	}
}
func (cmtt MultiColumnTimeseries) Less(i, j int) bool { return cmtt.Times[i].Before(cmtt.Times[j]) }
