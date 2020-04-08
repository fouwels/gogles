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
const _ioLoopTime = 1 * time.Millisecond

const _breathInFlowThreshold = 20
const _breathOutFlowThreshold = -20

//IOMan ..
type IOMan struct {
	flow1    *sfm3000.SFM3000
	o  DataPacket //Output data variables - external buffer. DataPacket is copied from working buffer to output buffer at end of io cycle.
	moutputs sync.Mutex
}

//NewIOMan ..
func NewIOMan() (*IOMan, error) {
	iom := IOMan{}

	iom.initialize()

	return &iom, nil
}

func (io *IOMan) initialize() error {

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
	io.flow1 = flow1
	logf("ioman", "All sensors initialized")

	return nil
}

//Start ..
func (io *IOMan) Start(cherr chan<- error) {
	logf("ioman", "Starting at %v hz", 1/_ioLoopTime.Seconds())
	lt := time.NewTicker(_ioLoopTime)
	defer lt.Stop()

	i  :=   internalPacket{} //Internal data variables

	for range lt.C {

		d := DataPacket{}

		// Read inputs
		fval, fcrc, ferr := io.flow1.GetValue()
		flow := Flow{
			Val: fval,
			CRC: fcrc,
			Err: ferr,
		} 

		d.Sensors = Sensors {
			Flow: flow,
		}
		d.Timestamp = time.Now()

		// If inputs read ok
		if d.Sensors.Flow.Err == nil {
			d.Stats.OkReads++

			io.states(&d, &i)
			io.calculate(&d, &i)
			d.Valid = true

		} else {
			d.Stats.FailedReads++

			d.Valid = false
		}

		// Copy to output registers
		io.moutputs.Lock()
		io.o = d
		io.moutputs.Unlock()
	}

	cherr <- fmt.Errorf("io loop ended unexpectedly")
}

func (io *IOMan) states(d *DataPacket, i *internalPacket) {
	breathin := d.Sensors.Flow.Val > _breathInFlowThreshold
	breathout := d.Sensors.Flow.Val < _breathOutFlowThreshold


	lastState := d.State
	state := StateRest

	if breathin {
		state = StateBreathingIn
		//logf("SPECIAL", "Breathing In %v", d.Sensors.Flow.Val)
	}
	if breathout {
		state = StateBreathingOut
		//logf("SPECIAL", "Breathing Out %v", d.Sensors.Flow.Val)
	}
	if breathin && breathout {
		state = StateError
		logf("ioman", "Breathing Error %v", d.Sensors.Flow.Val)
	}

	if state != lastState {
		i.lastStateChangeNmin1 = i.lastStateChangeN
		i.lastStateChangeN = d.Timestamp
	}
	i.lastState = lastState
	d.State = state
}

func (io *IOMan) calculate(d *DataPacket, i *internalPacket) {

	calc := Calculated {

	}

	// If moving into breathing state, reset counters
	if d.State == StateBreathingIn && i.lastState != StateBreathingIn {
		i.flowAverageTotal = 0
		i.flowAverageN = 0
	}

	// If in breathing state, increment counters
	if d.State == StateBreathingIn {
		i.flowAverageTotal += float64(d.Sensors.Flow.Val)
		i.flowAverageN++
	}

	// If leaving breathing state, calculate integrated flow.
	if d.State != StateBreathingIn && i.lastState == StateBreathingIn && d.State != StateError {
		duration := d.Timestamp.Sub(i.lastStateChangeNmin1)
		flow := (i.flowAverageTotal / float64(i.flowAverageN)) * duration.Minutes()

		calc.FlowIntegrated = flow
		calc.FlowIntegratedTimestamp = d.Timestamp
	}

	d.Calculated = calc
}

//GetDataPacket ..
func (io *IOMan) GetDataPacket() DataPacket {
	io.moutputs.Lock()
	defer io.moutputs.Unlock()
	return io.o
}

//Destroy ..
func (*IOMan) Destroy() {

}

func logf(owner string, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[%v] %v", owner, message)
}
