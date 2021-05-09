package ts_test

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/nordicsense/lstproxy/lib/ts"
)

var tt = ts.Timeseries{
	ts.TimePoint{time.Date(2012, 1, 22, 13, 34, 56, 0, time.UTC), 208.4},
	ts.TimePoint{time.Date(2012, 2, 22, 13, 34, 56, 0, time.UTC), 234.2},
	ts.TimePoint{time.Date(2012, 2, 24, 12, 34, 56, 0, time.UTC), 256.2},
	ts.TimePoint{time.Date(2012, 2, 25, 11, 34, 56, 0, time.UTC), 212.4},
}

func TestTimeseries_Sort(t *testing.T) {
	actual := ts.Timeseries{tt[3], tt[0], tt[2], tt[1]}
	actual.Sort()
	if !reflect.DeepEqual(tt, actual) {
		t.Errorf("datasets not equal, expected %v, found %v", tt, actual)
	}
}

func TestTimeseries_Interpolate_timeBetweenNodes(t *testing.T) {
	tm := time.Date(2012, 2, 1, 13, 34, 56, 0, time.UTC)
	expected := tt[0].Value + (tt[1].Value-tt[0].Value)*10.0/31.0
	assertInterpolation(t, tm, expected)
}

func TestTimeseries_Interpolate_timeAtNode(t *testing.T) {
	assertInterpolation(t, tt[1].Time, tt[1].Value)
}

func TestTimeseries_Interpolate_timeBeforeMin(t *testing.T) {
	assertInterpolation(t, time.Date(2012, 1, 22, 13, 0, 0, 0, time.UTC), math.NaN())
}

func TestTimeseries_Interpolate_timeAfterMax(t *testing.T) {
	assertInterpolation(t, time.Date(2012, 2, 25, 13, 0, 0, 0, time.UTC), math.NaN())
}

func TestTimeseries_Interpolate_timeAtMin(t *testing.T) {
	assertInterpolation(t, tt[0].Time, tt[0].Value)
}

func TestTimeseries_Interpolate_timeAtMax(t *testing.T) {
	assertInterpolation(t, tt[3].Time, tt[3].Value)
}

func assertInterpolation(t *testing.T, at time.Time, expected float64) {
	actual := tt.Interpolate(at)
	if math.IsNaN(expected) && !math.IsNaN(actual) {
		t.Errorf("expected NaN, found %v", actual)
	}
	if math.Abs(actual-expected) > 1e-5 {
		t.Errorf("expected %v, found %v", expected, actual)
	}
}

func TestTimeseries_Interpolate_oneNodeNaN(t *testing.T) {
	tm := time.Date(2012, 2, 23, 13, 34, 56, 0, time.UTC)
	data := ts.Timeseries{tt[0], tt[1], ts.TimePoint{tt[2].Time, math.NaN()}, tt[3]}
	actual := data.Interpolate(tm)
	if !math.IsNaN(actual) {
		t.Errorf("expected NaN, found %v", actual)
	}
}
