package ioman

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/guptarohit/asciigraph"
)

var _cont *controller

func TestMain(m *testing.M) {
	_cont = newController()

	result := m.Run()
	os.Exit(result)
}

func TestStates(t *testing.T) {

	flows, err := loadDebug()
	if err != nil {
		t.Fatalf("Failed to load debug files: %v", err)
	}
	t.Logf("Loaded %v flows", len(flows))

	vals := []float64{}
	for _, f := range flows {

		sensors := Sensors{
			Flow: f,
		}

		vals = append(vals, f.Val)
		_ = _cont.states(sensors)
	}

	graph := asciigraph.Plot(vals, asciigraph.Width(100), asciigraph.Height(10))

	fmt.Println(graph)

}

func loadDebug() ([]Flow, error) {
	datafile := "../../iodrivers/i2c/sfm3000/capture_datalog.csv"

	f, err := os.Open(datafile)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to open test data at %v: %v", datafile, err)
	}

	cs := csv.NewReader(f)
	records, err := cs.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read csv: %v", err)
	}

	flows := []Flow{}

	timenow := time.Now()

	for _, v := range records[:(len(records) - 1)] {

		offset, err := time.ParseDuration(v[0])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse duration: %v", v[0])
		}
		val, err := strconv.ParseFloat(v[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse value: %v", v[1])
		}

		flows = append(flows, Flow{
			Timestamp: timenow.Add(offset),
			Val:       val,
		})
	}
	return flows, nil
}
