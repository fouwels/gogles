package ioman

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

const _testSampleRate = (1 * time.Second) / 1000

var _cont *controller

func TestMain(m *testing.M) {
	_cont = newController(_testSampleRate)

	result := m.Run()
	os.Exit(result)
}

func TestStates(t *testing.T) {

	flows, err := loadDebug()
	if err != nil {
		t.Fatalf("Failed to load debug files: %v", err)
	}
	t.Logf("Loaded %v flows", len(flows))

	for _, f := range flows {

		sensors := Sensors{
			Flow: f,
		}

		time.Sleep(_testSampleRate)
		_cont.buffers(sensors)
		_cont.states(sensors)
		_cont.calculate(sensors)
	}

}

func loadDebug() ([]Flow, error) {
	datafile := "../../iodrivers/i2c/sfm3000/capture_datalog_2.csv"

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
