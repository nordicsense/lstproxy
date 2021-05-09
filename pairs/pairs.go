package main

import (
	"log"
	"path"
	"time"

	"github.com/nordicsense/modis"

	"github.com/nordicsense/lstproxy/pairs/fail"
	"github.com/nordicsense/lstproxy/pairs/ts"

	modists "github.com/nordicsense/modis/ts"
)

type datasetPair struct {
	name string
	root string
	pair modists.LayerPair
}

type station struct {
	name string
	ll   modis.LatLon
}

var (
	lstDay = modists.LayerPair{
		Time:  `^.+_LST:Day_view_time$`,
		Value: `^.+_LST:LST_Day_.+$`,
	}
	lstNight = modists.LayerPair{
		Time:  `^.+_LST:Night_view_time$`,
		Value: `^.+_LST:LST_Night_.+$`,
	}
)

func init() {
	log.SetFlags(0)
}

func main() {
	researchRoot := "/data/Research/Data"
	outputRoot := path.Join(researchRoot, "Analysis", "Fin-Meteo")
	meteoRoot := path.Join(researchRoot, "Data", "Meteo", "Finland")
	aquaRoot := path.Join(researchRoot, "Data", "MODIS", "MOD11A1", "Aqua-HDF")
	terraRoot := path.Join(researchRoot, "Data", "MODIS", "MOD11A1", "Terra-HDF")

	m := &runner{
		stationReader: &finnishStationReader{
			root:    meteoRoot,
			pattern: "Fin.+Stations.csv",
			skip:    1,
		},
		modisReader: &modisReader{},
		resultWriter: &cmttWriter{
			root:   outputRoot,
			header: []string{"LST", "t", "rh", "vis"},
		},
		matchers: map[string]interface {
			Match(map[int][]time.Time) (map[int][]float64, error)
		}{
			"t": &finnishMeteoReader{
				root:    meteoRoot,
				pattern: "obs-.+-t2m.csv",
			},
			"rh": &finnishMeteoReader{
				root:    meteoRoot,
				pattern: "obs-.+-rh.csv",
			},
			"vis": &finnishMeteoReader{
				root:    meteoRoot,
				pattern: "obs-.+-vis.csv",
			},
		},
		datasetsPairs: []datasetPair{
			{name: "Aqua-Day", root: aquaRoot, pair: lstDay},
			{name: "Aqua-Night", root: aquaRoot, pair: lstNight},
			{name: "Terra-Day", root: terraRoot, pair: lstDay},
			{name: "Terra-Night", root: terraRoot, pair: lstNight},
		},
	}
	fail.OnError(m.run())
}

type runner struct {
	stationReader interface {
		Read() (map[int]station, error)
	}
	modisReader interface {
		Read(ds datasetPair, stations map[int]station) (map[int]ts.MultiColumnTimeseries, error)
	}
	resultWriter interface {
		Write(dsname string, data map[int]ts.MultiColumnTimeseries) error
	}
	matchers map[string]interface {
		Match(map[int][]time.Time) (map[int][]float64, error)
	}
	datasetsPairs []datasetPair
}

func (r *runner) run() error {
	log.Println("Reading station data")
	stations, err := r.stationReader.Read()
	if err != nil {
		return err
	}

	for _, dsPair := range r.datasetsPairs {
		log.Println("Processing", dsPair.name)
		cmtts, err := r.modisReader.Read(dsPair, stations)
		if err != nil {
			return err
		}
		tms := make(map[int][]time.Time)
		for id, cmtt := range cmtts {
			tms[id] = cmtt.Times
		}
		for colname, matcher := range r.matchers {
			log.Println("Matching", colname)
			cols, err := matcher.Match(tms)
			if err != nil {
				return err
			}
			for id, col := range cols {
				cmtts[id].Values[colname] = col
			}
		}
		log.Println("Storing", dsPair.name)
		if err := r.resultWriter.Write(dsPair.name, cmtts); err != nil {
			return err
		}
	}
	return nil
}
