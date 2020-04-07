package ioman

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kaelanfouwels/iodrivers/i2c"
	"github.com/kaelanfouwels/iodrivers/i2c/sfm3000"
)

const _sfm3000Address = 0x40
const _sfm3000IsAir = false
const _i2cBus = 1
const _ioLoopTime = 10 * time.Millisecond

const _breathInFlowThreshold = 0.5
const _breathOutFlowThreshold = -0.5

//IOMan ..
type IOMan struct {
	flow1    *sfm3000.SFM3000
	d        DataPacket
	i        internalPacket
	outputs  DataPacket
	moutputs sync.Mutex
}

//NewIOMan ..
func NewIOMan() (*IOMan, error) {
	iom := IOMan{}

	iom.initialize()

	return &iom, nil
}

func (i *IOMan) initialize() error {

	//Initialize I2C
	logf("ioman", "Initializing I2C on bus %v", _i2cBus)
	i21, err := i2c.NewI2C(_i2cBus)
	if err != nil {
		return fmt.Errorf("Failed to create I2C device: %w", err)
	}

	//Initialize SFM3000
	logf("ioman", "Initializing SFM3000 on I2C bus %v 0x%x as %v", _i2cBus, _sfm3000Address, "FLOW1")
	flow1, err := sfm3000.NewSFM3000(i21, _sfm3000Address, _sfm3000IsAir, "FLOW1")
	if err != nil {
		return fmt.Errorf("Failed to create SFM3000: %w", err)
	}

	logf("ioman", "Testing %v", flow1.Label())
	err = flow1.SoftReset()
	if err != nil {
		return fmt.Errorf("Failed to soft reset %v: %w", flow1.Label(), err)
	}
	logf("ioman", "%v soft-reset OK", flow1.Label())
	time.Sleep(50 * time.Millisecond)

	serial, err := flow1.GetSerial()
	if err != nil {
		return fmt.Errorf("Failed to get %v serial number: %w", flow1.Label(), err)
	}
	logf("ioman", "%v serial number OK: %v", flow1.Label(), serial)
	time.Sleep(50 * time.Millisecond)

	_, _, _ = flow1.GetValue() //First value is expected to be garbage
	time.Sleep(50 * time.Millisecond)

	value, crc, err := flow1.GetValue()
	if err != nil {
		return fmt.Errorf("Failed to get %v value: %w", flow1.Label(), err)
	}
	logf("ioman", "%v flow value OK: %v crc %v", flow1.Label(), value, crc)
	i.flow1 = flow1
	logf("ioman", "All sensors initialized")

	return nil
}

//Start ..
func (i *IOMan) Start(cherr chan<- error) {
	logf("ioman", "Starting at %v hz", 1/_ioLoopTime.Seconds())
	lt := time.NewTicker(_ioLoopTime)
	defer lt.Stop()

	for range lt.C {

		// Read inputs
		i.d.Sensors.Flow.Val, i.d.Sensors.Flow.CRC, i.d.Sensors.Flow.Err = i.flow1.GetValue()
		i.d.Timestamp = time.Now()

		// If inputs read ok
		if i.d.Sensors.Flow.Err == nil {
			i.d.Stats.OkReads++

			i.states()
			i.calculate()
			i.d.Valid = true

		} else {
			i.d.Stats.FailedReads++

			i.d.Valid = false
		}

		// Copy to output registers
		i.moutputs.Lock()
		i.outputs = i.d
		i.moutputs.Unlock()
	}

	cherr <- fmt.Errorf("io loop ended unexpectedly")
}

func (i *IOMan) states() {
	breathin := i.d.Sensors.Flow.Val > _breathInFlowThreshold
	breathout := i.d.Sensors.Flow.Val < _breathOutFlowThreshold

	var state enumState = StateRest

	if breathin {
		state = StateBreathingIn
	}
	if breathout {
		state = StateBreathingOut
	}
	if breathin && breathout {
		state = StateError
	}

	if state != i.i.lastState {
		i.i.lastStateChange = i.d.Timestamp
	}

	i.d.State = state
}

func (i *IOMan) calculate() {

	// If moving into breathing state, reset counters
	if i.d.State == StateBreathingIn && i.i.lastState != StateBreathingIn {
		i.i.flowAverageTotal = 0
		i.i.flowAverageN = 0
	}

	// If in breathing state, increment counters
	if i.d.State == StateBreathingIn {
		i.i.flowAverageTotal += float64(i.d.Sensors.Flow.Val)
		i.i.flowAverageN++
	}

	// If leaving breathing state, calculate integrated flow.
	if i.d.State != StateBreathingIn && i.i.lastState == StateBreathingIn {
		duration := i.d.Timestamp.Sub(i.i.lastStateChange)
		flow := (i.i.flowAverageTotal / float64(i.i.flowAverageN)) * duration.Seconds()

		i.d.Calculated.FlowIntegrated = flow
		i.d.Calculated.FlowIntegratedTimestamp = i.i.lastStateChange
	}

}

//GetDataPacket ..
func (i *IOMan) GetDataPacket() DataPacket {
	i.moutputs.Lock()
	defer i.moutputs.Unlock()
	return i.outputs
}

//Destroy ..
func (*IOMan) Destroy() {

}

func logf(owner string, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[%v] %v", owner, message)
}
