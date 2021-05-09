package main

import (
	"math"
	"runtime"
	"time"

	"github.com/nordicsense/lstproxy/lib/progress"
	"github.com/nordicsense/lstproxy/lib/ts"
	"github.com/nordicsense/modis/dataset"

	modists "github.com/nordicsense/modis/ts"
)

const LST = "LST"

type modisReader struct{}

func (r *modisReader) Read(dsPair datasetPair, stations map[int]station) (map[int]ts.MultiColumnTimeseries, error) {
	pairs, err := modists.ListAll(dsPair.root, dsPair.pair)
	if err != nil {
		return nil, err
	}

	bar := progress.Start("Read", len(pairs))
	defer bar.Finish()

	res := make(map[int]ts.MultiColumnTimeseries)
	for id := range stations {
		res[id] = ts.MultiColumnTimeseries{
			Times: make([]time.Time, len(pairs)),
			Values: map[string][]float64{
				LST: make([]float64, len(pairs)),
			},
		}
	}

	is := make(map[int]int)
	for _, dsPair := range pairs {
		bar.Add(1)
		ds, err := dataset.Open(dsPair.Value)
		if err != nil {
			return nil, err
		}
		tds, err := dataset.Open(dsPair.Time)
		if err != nil {
			ds.Close()
			return nil, err
		}

		err = nil
		for id, st := range stations {
			if !ds.ImageParams().Within(st.ll) {
				continue
			}
			var hdfval float64
			hdfval, err = ds.ReadAtLatLon(st.ll)
			if err != nil {
				break
			} else if math.IsNaN(hdfval) {
				continue
			}
			hdfval -= 273.15
			var tm time.Time
			tm, err = tds.ReadTimeAtLatLon(st.ll)
			if err != nil {
				break
			} else if tm == tds.ImageParams().Date() {
				continue
			}
			cmtt := res[id]
			cmtt.Times[is[id]] = tm
			cmtt.Values[LST][is[id]] = hdfval
			is[id] += 1
		}

		ds.Close()
		tds.Close()
		if err != nil {
			return nil, err
		}
	}
	for id, cmtt := range res {
		n := is[id]
		tms := make([]time.Time, n)
		copy(tms, cmtt.Times[:n])
		vals := make([]float64, n)
		copy(vals, cmtt.Values[LST][:n])
		cmtt = ts.MultiColumnTimeseries{Times: tms, Values: map[string][]float64{LST: vals}}
		cmtt.Sort()
		res[id] = cmtt
		runtime.GC()
	}
	return res, nil
}
