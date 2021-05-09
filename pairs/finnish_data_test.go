package main

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/nordicsense/lstproxy/pairs/ts"
	"github.com/nordicsense/modis"
)

func TestFinnishStationReader_Read_GoodPath(t *testing.T) {
	root, _ := os.Getwd()
	r := &finnishStationReader{
		root:    root,
		pattern: "Fin.+Stations.csv",
		skip:    1,
	}
	stations, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	expected := map[int]station{
		100919: {name: "Hammarland Notes", ll: modis.LatLon{60.3, 19.13}},
		100924: {name: "Pargas Fagerholm", ll: modis.LatLon{60.11, 21.7}},
		100928: {name: "Kumlinge Parish", ll: modis.LatLon{60.26, 20.75}},
	}
	if !reflect.DeepEqual(stations, expected) {
		t.Errorf("expected %#v, found %#v", expected, stations)
	}
}

func TestFinnishStationReader_Read_WithoutSkip(t *testing.T) {
	root, _ := os.Getwd()
	r := &finnishStationReader{
		root:    root,
		pattern: "Fin.+Stations.csv",
	}
	expectedError := "strconv.ParseFloat: parsing \"Lat\": invalid syntax"
	_, err := r.Read()
	if err == nil {
		t.Error("expected error, found none")
	} else if err.Error() != expectedError {
		t.Errorf("expected %s, found %s", expectedError, err.Error())
	}
}

func TestFinnishStationReader_Read_BadPattern(t *testing.T) {
	root, _ := os.Getwd()
	r := &finnishStationReader{
		root:    root,
		pattern: "XXX.+Stations.csv",
		skip:    1,
	}
	expectedError := "expected 1 file with station data, found 0"
	_, err := r.Read()
	if err == nil {
		t.Error("expected error, found none")
	} else if err.Error() != expectedError {
		t.Errorf("expected %s, found %s", expectedError, err.Error())
	}
}

func TestFinnishStationReader_Read_BadRoot(t *testing.T) {
	r := &finnishStationReader{
		root:    "/dev/null",
		pattern: "Fin.+Stations.csv",
		skip:    1,
	}
	expectedError := "readdirent: not a directory"
	_, err := r.Read()
	if err == nil {
		t.Error("expected error, found none")
	} else if err.Error() != expectedError {
		t.Errorf("expected %s, found %s", expectedError, err.Error())
	}
}

/*
102052,2003-01-04 16:00,-14.8000
102052,2003-01-04 17:00,-14.7000
102052,2003-01-04 23:00,-22.4000

100907,2003-01-04 17:00,-14.70
100907,2003-01-04 23:00,21

*/

func TestFinnishMeteoReader_Read_GoodPath(t *testing.T) {
	root, _ := os.Getwd()
	r := &finnishMeteoReader{
		root:    root,
		pattern: "obs-.+-t2m.csv",
	}
	data, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	parse := func(v string) time.Time {
		res, _ := time.Parse(dateTimeLayout, v)
		return res
	}

	expected := map[int]ts.Timeseries{
		102052: {
			ts.TimePoint{Time: parse("2003-01-04 16:00"), Value: -14.8},
			ts.TimePoint{Time: parse("2003-01-04 17:00"), Value: -14.7},
			ts.TimePoint{Time: parse("2003-01-04 23:00"), Value: -22.4},
		},
		100907: {
			ts.TimePoint{Time: parse("2003-01-04 17:00"), Value: -14.7},
			ts.TimePoint{Time: parse("2003-01-04 23:00"), Value: 21},
		},
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("expected %#v, found %#v", expected, data)
	}
}
