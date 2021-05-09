package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/nordicsense/lstproxy/lib/ts"
)

type cmttWriter struct {
	root   string
	header []string
}

func (mttw *cmttWriter) Write(dsname string, cmtts map[int]ts.MultiColumnTimeseries) error {
	var w *csv.Writer
	if f, err := os.Create(path.Join(mttw.root, fmt.Sprintf("%s.csv", dsname))); err == nil {
		defer func() { _ = f.Close() }()
		w = csv.NewWriter(f)
	} else {
		return err
	}
	record := append([]string{"id", "datetime"}, mttw.header...)

	if err := w.Write(record); err != nil {
		return err
	}
	for id, cmtt := range cmtts {
		if len(cmtt.Times) == 0 {
			continue
		}
		record[0] = strconv.Itoa(id)
		for i, tm := range cmtt.Times {
			record[1] = tm.Format("2006-01-02T15:04:05")
			for j, colname := range mttw.header {
				vector, ok := cmtt.Values[colname]
				if !ok {
					return fmt.Errorf("missing column %s", colname)
				}
				if len(vector) != len(cmtt.Times) {
					return fmt.Errorf("length mismatch for column %s", colname)
				}
				record[j+2] = strconv.FormatFloat(vector[i], 'f', 5, 64)
			}
			if err := w.Write(record); err != nil {
				return err
			}
		}
	}
	w.Flush()
	return nil
}
