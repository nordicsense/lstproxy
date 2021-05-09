package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v2"
	"github.com/subchen/go-xmldom"
	"github.com/teris-io/cli"
)

const (
	serverTimeFormat = "2006-01-02T15:04:05Z"
	outTimeFormat    = "2006-01-02 15:04"
	urlBase          = "http://opendata.fmi.fi/wfs?request"
	command          = "getFeature"
	storedQueryID    = "fmi::observations::weather::timevaluepair"
	maxLocations     = 1
)

func main() {
	app := cli.New("Mass-downloader for Finnish meteo data").
		WithOption(cli.NewOption("from", "From date as 2006-12-25 (default: -1y from today)").WithChar('f')).
		WithOption(cli.NewOption("to", "Till date as 2006-12-25 (default: today)").WithChar('t')).
		WithOption(cli.NewOption("params", "Parameters, coma separated w/o spaces (default: all)").WithChar('p')).
		WithOption(cli.NewOption("col", "CSV column index (default: 2)").WithChar('c').WithType(cli.TypeInt)).
		WithOption(cli.NewOption("skip", "CSV lines to skip (default: 1)").WithType(cli.TypeInt)).
		WithOption(cli.NewOption("step", "Time step in minutes (default: 60)").WithChar('s').WithType(cli.TypeInt)).
		WithArg(cli.NewArg("stations", "Stations CSV file")).
		WithArg(cli.NewArg("path", "Output path (default: working directory)").AsOptional()).
		WithAction(func(args []string, options map[string]string) int {
			var params []string
			paramsStr, ok := options["params"]
			if ok {
				params = strings.Split(paramsStr, ",")
			}
			stationCol, ok := options["col"]
			if !ok {
				stationCol = "3"
			}
			colID, _ := strconv.Atoi(stationCol)
			colID-- // 1-based input
			skipStr, ok := options["skip"]
			if !ok {
				skipStr = "1"
			}
			skip, _ := strconv.Atoi(skipStr)
			stationIDs := readStationIDs(args[0], skip, colID)
			from := time.Now().AddDate(-1, 0, 0)
			fromStr, ok := options["from"]
			if ok {
				var err error
				from, err = time.Parse("2006-01-02", fromStr)
				fail(err)
			}
			to := time.Now()
			toStr, ok := options["to"]
			if ok {
				var err error
				to, err = time.Parse("2006-01-02", toStr)
				fail(err)
			}
			timeStepStr, ok := options["step"]
			if !ok {
				timeStepStr = "60"
			}
			timeStep, err := strconv.Atoi(timeStepStr)
			fail(err)
			root, err := os.Getwd()
			fail(err)
			if len(args) > 1 {
				root = args[1]
			}
			return run(stationIDs, from, to, timeStep, root, params...)
		})
	os.Exit(app.Run(os.Args, os.Stderr))
}

func readStationIDs(stationFile string, skip, colID int) []int {
	f, err := os.Open(stationFile)
	fail(err)
	defer func() { _ = f.Close() }()
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	fail(err)
	var res []int
	for i, line := range lines {
		if i < skip {
			continue
		}
		if len(line) <= colID {
			fail(errors.New("stations file line too short for column ID"))
		}
		id, err := strconv.Atoi(line[colID])
		fail(err)
		res = append(res, id)
	}
	return res
}

func run(stationIDs []int, from, to time.Time, timeStep int, root string, params ...string) int {
	var dates []time.Time
	for date := from; date.Before(to); date = date.AddDate(0, 0, 5) {
		dates = append(dates, date)
	}
	count := len(dates) * len(stationIDs)
	datachan := make(chan map[string]observations)
	throttling := make(chan bool, 10)
	for _, date := range dates {
		for _, stationID := range stationIDs {
			go func(stationID int, date time.Time) {
				throttling <- true
				start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
				end := start.AddDate(0, 0, 5).Add(-5 * time.Minute)
				u := url(stationID, start, end, timeStep, params...)
				data, err := extract(get(u), stationID)
				if err == nil {
					datachan <- data
				} else {
					fmt.Printf("Failed for %v: %v\n", u, err) // do not fail fully
				}
				<-throttling
			}(stationID, date)
		}

	}
	res := make(map[string]observations)
	bar := progressbar.New(count)
	for i := 0; i < count; i++ {
		for key, obss := range <-datachan {
			res[key] = append(res[key], obss...)
		}
		_ = bar.Add(1)
	}
	_ = bar.Finish()
	close(datachan)
	for name, obss := range res {
		sort.Sort(obss)
		write(root, name, obss)
	}
	return 0
}

func url(stationID int, start, end time.Time, timeStep int, params ...string) string {
	storedQueryArg := fmt.Sprintf("storedquery_id=%s", storedQueryID)
	fmisArg := fmt.Sprintf("fmisid=%d", stationID)
	maxLocationsArg := fmt.Sprintf("maxlocations=%d", maxLocations)
	startTimeArg := fmt.Sprintf("starttime=%s", start.Format(serverTimeFormat))
	endTimeArg := fmt.Sprintf("endtime=%s", end.Format(serverTimeFormat))
	timeStepArg := fmt.Sprintf("timestep=%d", timeStep)
	request := strings.Join([]string{command, storedQueryArg, fmisArg, maxLocationsArg, startTimeArg, endTimeArg, timeStepArg}, "&")
	if len(params) > 0 {
		request = fmt.Sprintf("%s&parameters=%s", request, strings.Join(params, ","))
	}
	return strings.Join([]string{urlBase, request}, "=")
}

func get(url string) io.ReadCloser {
	xmlReader, err := http.Get(url)
	fail(err)
	return xmlReader.Body
}

type observation struct {
	stationID int
	time      time.Time
	value     float64
}

type observations []observation

func (o observations) Len() int { return len(o) }
func (o observations) Less(i, j int) bool {
	if o[i].stationID != o[j].stationID {
		return o[i].stationID < o[j].stationID
	}
	return o[i].time.Before(o[j].time)
}
func (o observations) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

// XML parsing susceptible to nil dereference: no intention to fix due to CLI nature of the tool
func extract(xmlReader io.Reader, stationID int) (map[string]observations, error) {
	dom, err := xmldom.Parse(xmlReader)
	if err != nil {
		return nil, err
	}
	res := make(map[string]observations)
	for _, dsNode := range dom.Root.Children {
		tsNode := dsNode.Children[0]
		dataNode := tsNode.FindOneByName("MeasurementTimeseries")
		if dataNode == nil {
			return nil, errors.New("bad response")
		}
		nodeName := dataNode.GetAttribute("id").Value
		for _, pt := range dataNode.Children {
			if len(pt.Children) != 1 {
				return nil, errors.New("expected 1 MeasurementTVP child")
			}
			mtvp := pt.Children[0]
			timeNode := mtvp.FindOneByName("time")
			valueNode := mtvp.FindOneByName("value")
			if timeNode == nil || valueNode == nil {
				return nil, errors.New("expected one time and one value node")
			}
			tm, err := time.Parse(serverTimeFormat, timeNode.Text)
			if err != nil {
				return nil, err
			}
			val, err := strconv.ParseFloat(valueNode.Text, 64)
			if err != nil {
				return nil, err
			}
			res[nodeName] = append(res[nodeName], observation{stationID: stationID, time: tm, value: val})
		}
	}
	return res, nil
}

func write(root, name string, obss observations) {
	fname := path.Join(root, fmt.Sprintf("%s.csv", name))
	f, err := os.Create(fname)
	fail(err)
	defer func() { _ = f.Close() }()
	w := csv.NewWriter(f)
	fields := []string{"", "", ""}
	for _, obs := range obss {
		fields[0] = strconv.FormatInt(int64(obs.stationID), 10)
		fields[1] = obs.time.Format(outTimeFormat)
		fields[2] = strconv.FormatFloat(obs.value, 'f', 4, 64)
		fail(w.Write(fields))
	}
	w.Flush()
}

func fail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
