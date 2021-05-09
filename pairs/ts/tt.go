package ts

import (
	"math"
	"sort"
	"time"
)

type TimePoint struct {
	Time  time.Time
	Value float64
}

type Timeseries []TimePoint

func (tt Timeseries) Sort() {
	sort.Sort(tt)
}

func (tt Timeseries) Len() int           { return len(tt) }
func (tt Timeseries) Swap(i, j int)      { tt[i], tt[j] = tt[j], tt[i] }
func (tt Timeseries) Less(i, j int) bool { return tt[i].Time.Before(tt[j].Time) }

func (tt Timeseries) Interpolate(t time.Time) float64 {
	i := sort.Search(len(tt), func(i int) bool { return !tt[i].Time.Before(t) })
	if i < 0 || i == len(tt) {
		return math.NaN()
	}
	if t.Equal(tt[i].Time) {
		return tt[i].Value
	}
	if i == 0 {
		return math.NaN()
	}
	n := t.Sub(tt[i-1].Time)
	d := tt[i].Time.Sub(tt[i-1].Time)
	return tt[i-1].Value + (tt[i].Value-tt[i-1].Value)*n.Seconds()/d.Seconds()
}

// assume sorted
func (tt Timeseries) Slice(from time.Time, to time.Time) Timeseries {
	res := tt
	index := sort.Search(len(res), func(i int) bool { return from.Before(res[i].Time) })
	if index >= 0 {
		res = res[index:]
	}
	index = sort.Search(len(res), func(i int) bool { return to.Before(res[i].Time) })
	if index > 0 {
		res = res[:(index - 1)]
	}
	return res
}

// assume sorted, assume no NaNs
func (tt Timeseries) Dedupe() Timeseries {
	if len(tt) == 0 {
		return nil
	}
	res := Timeseries{tt[0]}
	for i := 1; i < len(tt); i++ {
		v := tt[i]
		last := res[i-1]
		if last.Time == v.Time {
			res[i-1].Value = 0.5 * (last.Value + v.Value)
		} else {
			res = append(res, v)
		}
	}
	return res
}

func (tt Timeseries) Transform(fun func(float64) float64) Timeseries {
	var res Timeseries
	for _, tp := range tt {
		res = append(res, TimePoint{Time: tp.Time, Value: fun(tp.Value)})
	}
	return res
}

func (tt Timeseries) ToMultiColumnTimeseries(colname string) MultiColumnTimeseries {
	res := MultiColumnTimeseries{
		Times:  make([]time.Time, len(tt)),
		Values: map[string][]float64{colname: make([]float64, len(tt))},
	}
	vector := res.Values[colname]
	for i, tp := range tt {
		res.Times[i] = tp.Time
		vector[i] = tp.Value
	}
	return res
}
