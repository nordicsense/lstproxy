package ts

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type MultiTimePoint struct {
	Time   time.Time
	Values []float64
}

type MultiTimeseries []MultiTimePoint

func (mtt MultiTimeseries) Sort() {
	sort.Sort(mtt)
}

func (mtt MultiTimeseries) Len() int           { return len(mtt) }
func (mtt MultiTimeseries) Swap(i, j int)      { mtt[i], mtt[j] = mtt[j], mtt[i] }
func (mtt MultiTimeseries) Less(i, j int) bool { return mtt[i].Time.Before(mtt[j].Time) }

func Join(tts ...Timeseries) MultiTimeseries {
	if len(tts) == 0 {
		return nil
	}
	data := make(map[time.Time]MultiTimePoint)
	for tti, tt := range tts {
		for _, tp := range tt {
			r, ok := data[tp.Time]
			if !ok {
				r = MultiTimePoint{Time: tp.Time, Values: make([]float64, len(tts))}
				for i := 0; i < len(tts); i++ {
					r.Values[i] = math.NaN()
				}
			}
			r.Values[tti] = tp.Value
			data[tp.Time] = r
		}
	}
	var res MultiTimeseries
	for _, r := range data {
		res = append(res, r)
	}
	res.Sort()
	return res
}

func Cbind(ts Timeseries, cols ...[]float64) (MultiTimeseries, error) {
	for j, col := range cols {
		if len(col) != len(ts) {
			return nil, fmt.Errorf("length of column %d differs from timeseries", j)
		}
	}
	res := make(MultiTimeseries, len(ts))
	for i, tp := range ts {
		mtp := MultiTimePoint{Time: tp.Time, Values: make([]float64, len(cols)+1)}
		for j := range mtp.Values {
			mtp.Values[j] = math.NaN()
		}
		mtp.Values[0] = tp.Value
		for j, col := range cols {
			mtp.Values[j+1] = col[i]
		}
		res[i] = mtp
	}
	return res, nil
}
