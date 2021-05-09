package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/nordicsense/lstproxy/pairs/ts"

	io2 "github.com/nordicsense/lstproxy/pairs/io"
)

const dateTimeLayout = "2006-01-02 15:04"

type finnishMeteoReader struct {
	root    string
	pattern string
}

func (fm *finnishMeteoReader) Read() (map[int]ts.Timeseries, error) {
	files, err := io2.ScanTree(fm.root, fm.pattern)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("expected at least 1 file with ground data, found none")
	}

	res := make(map[int]ts.Timeseries)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}

		// wrapping into function to simplify file closing upon completion
		err = func() error {
			r := csv.NewReader(f)
			r.ReuseRecord = true
			r.FieldsPerRecord = 3
			for {
				record, err := r.Read()
				if err == io.EOF {
					// exit loop and function
					return nil
				} else if err == csv.ErrFieldCount {
					continue
				} else if err != nil {
					return err
				}
				v, err := strconv.ParseFloat(record[2], 64)
				if err != nil {
					return err
				}
				if math.IsNaN(v) {
					continue
				}
				id, err := strconv.Atoi(record[0])
				if err != nil {
					return err
				}
				dt, err := time.Parse(dateTimeLayout, record[1])
				if err != nil {
					return err
				}
				res[id] = append(res[id], ts.TimePoint{Time: dt, Value: v})
			}
		}()
		_ = f.Close()
		if err != nil {
			return nil, err
		}
	}
	for _, tt := range res {
		tt.Sort()
	}
	return res, nil
}

func (fm *finnishMeteoReader) Match(times map[int][]time.Time) (map[int][]float64, error) {
	cdata, err := fm.Read()
	if err != nil {
		return nil, err
	}
	res := make(map[int][]float64)
	for id, tms := range times {
		ctt, ok := cdata[id]
		var vals []float64
		for _, tm := range tms {
			val := math.NaN()
			if ok {
				val = ctt.Slice(tm.Add(-4*time.Hour), tm.Add(4*time.Hour)).Interpolate(tm)
			}
			vals = append(vals, val)
		}
		res[id] = vals
	}
	return res, nil
}

type finnishStationReader struct {
	root    string
	pattern string
	skip    int
}

func (fs *finnishStationReader) Read() (map[int]station, error) {
	files, err := io2.ScanTree(fs.root, fs.pattern)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 || len(files) > 1 {
		return nil, fmt.Errorf("expected 1 file with station data, found %d", len(files))
	}

	f, err := os.Open(files[0])
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	res := make(map[int]station)
	// --, Espoo Nuuksio, 852678, --, --, 60.29, 24.57, --
	counter := 0
	for _, record := range records {
		counter++
		if len(record) != 8 {
			return nil, fmt.Errorf("expected 8 items in a record, found %d", len(record))
		}
		if counter <= fs.skip {
			continue
		}
		var ll [2]float64
		ll[0], err = strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, err
		}
		ll[1], err = strconv.ParseFloat(record[6], 64)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}
		res[id] = station{
			name: record[1],
			ll:   ll,
		}
	}
	return res, nil
}
